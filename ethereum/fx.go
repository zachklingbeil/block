package ethereum

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
)

// ListenForNewBlocksIPC subscribes to new block headers using IPC and updates e.Header.
func (e *Ethereum) Listen() error {
	headers := make(chan *types.Header)
	sub, err := e.Factory.Eth.SubscribeNewHead(e.Factory.Ctx, headers)
	if err != nil {
		return fmt.Errorf("failed to subscribe to new heads: %w", err)
	}

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case <-e.Factory.Ctx.Done():
				return
			case err := <-sub.Err():
				log.Printf("Subscription error: %v", err)
				return
			case header := <-headers:
				e.Header = header.Number
				e.BlockByBlock()
			}
		}
	}()
	return nil
}

// BlockByBlock processes and stores the latest `count` blocks, one at a time.
func (e *Ethereum) BlockByBlock() error {
	block, err := e.Factory.Eth.BlockByNumber(e.Factory.Ctx, e.Header)
	if err != nil {
		log.Printf("Error fetching block %d: %v", e.Header.Uint64(), err)
	}
	blockInfo := e.processBlock(e.Factory.Ctx, block)
	err = e.StoreBlock(int64(blockInfo.Number), blockInfo)
	if err != nil {
		log.Printf("Error storing block %d: %v", blockInfo.Number, err)
	}
	fmt.Printf("%d\n", e.Header.Uint64())
	return nil
}

func (e *Ethereum) StoreBlock(blockNumber int64, block *Raw) error {
	blockJSON, err := json.Marshal(block)
	if err != nil {
		return fmt.Errorf("failed to marshal block: %w", err)
	}
	err = e.Factory.Data.RB.SAdd(e.Factory.Ctx, "ethereum", blockJSON).Err()
	if err != nil {
		return fmt.Errorf("failed to store block in Redis hash: %w", err)
	}
	return nil
}
