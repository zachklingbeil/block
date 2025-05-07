package ethereum

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
)

// ProcessBlocks processes and stores the latest `count` blocks, one at a time.
func (e *Ethereum) ProcessBlocks(count int) error {
	header, err := e.Factory.Eth.HeaderByNumber(e.Factory.Ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to get latest header: %w", err)
	}
	latestBlock := header.Number.Uint64()

	for blockNum := latestBlock; blockNum > latestBlock-uint64(count); blockNum-- {
		block, err := e.Factory.Eth.BlockByNumber(e.Factory.Ctx, big.NewInt(int64(blockNum)))
		if err != nil {
			log.Printf("Error fetching block %d: %v", blockNum, err)
			continue
		}
		blockInfo := e.processBlock(e.Factory.Ctx, block)
		if err := e.StoreBlock(int64(blockInfo.Number), blockInfo); err != nil {
			log.Printf("Error storing block %d: %v", blockInfo.Number, err)
			continue
		}
		fmt.Printf("%d", blockInfo.Number)
	}
	return nil
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
