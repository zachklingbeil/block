package ethereum

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/core/types"
)

// ListenForNewBlocksIPC subscribes to new block headers using IPC and updates e.Header.
func (e *Ethereum) Listen(ctx context.Context) error {
	headers := make(chan *types.Header)
	sub, err := e.Factory.Eth.SubscribeNewHead(ctx, headers)
	if err != nil {
		return fmt.Errorf("failed to subscribe to new heads: %w", err)
	}

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-sub.Err():
				log.Printf("Subscription error: %v", err)
				return
			case header := <-headers:
				e.Header = int64(header.Number.Uint64())
				log.Printf("New block: %d", e.Header)
			}
		}
	}()
	return nil
}

// ProcessBlocks processes the latest `count` blocks.
func (e *Ethereum) ProcessBlocks(count int) ([]*Raw, error) {
	header, err := e.Factory.Eth.HeaderByNumber(e.Factory.Ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest header: %w", err)
	}
	latestBlock := header.Number.Uint64()

	var blocks []*Raw
	for blockNum := latestBlock; blockNum > latestBlock-uint64(count); blockNum-- {
		block, err := e.Factory.Eth.BlockByNumber(e.Factory.Ctx, big.NewInt(int64(blockNum)))
		if err != nil {
			log.Printf("Error fetching block %d: %v", blockNum, err)
			continue
		}
		blockInfo := e.processBlock(e.Factory.Ctx, block)
		blocks = append(blocks, blockInfo)
		err = e.StoreBlock(int64(blockInfo.Number), blockInfo)
		if err != nil {
			log.Printf("Error storing block %d: %v", blockInfo.Number, err)
			continue
		}
		fmt.Printf("%d", blockInfo.Number)
	}
	return blocks, nil
}

func (e *Ethereum) StoreBlock(blockNumber int64, block *Raw) error {
	blockJSON, err := json.Marshal(block)
	if err != nil {
		return fmt.Errorf("failed to marshal block: %w", err)
	}
	err = e.Factory.Data.RB.HSet(e.Factory.Ctx, "ethereum", strconv.Itoa(int(blockNumber)), blockJSON).Err()
	if err != nil {
		return fmt.Errorf("failed to store block in Redis hash: %w", err)
	}
	return nil
}
