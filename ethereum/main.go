package ethereum

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/zachklingbeil/block/value"
	"github.com/zachklingbeil/factory"
)

type Ethereum struct {
	Factory *factory.Factory
	Value   *value.Value
	Chain   *params.ChainConfig
}

func NewEthereum(factory *factory.Factory, value *value.Value) *Ethereum {
	return &Ethereum{
		Factory: factory,
		Value:   value,
		Chain:   params.MainnetChainConfig,
	}
}

// Signer returns a signer for Ethereum mainnet at the given block number and time.
func (e *Ethereum) Signer(blockNumber *big.Int, blockTime uint64) types.Signer {
	return types.MakeSigner(e.Chain, blockNumber, blockTime)
}

// ProcessBlocks processes the latest `count` blocks.
func (e *Ethereum) ProcessBlocks(count int) ([]*Block, error) {
	header, err := e.Factory.Eth.HeaderByNumber(e.Factory.Ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest header: %w", err)
	}
	latestBlock := header.Number.Uint64()

	var blocks []*Block
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
		log.Printf("%d", blockInfo.Number)
	}
	return blocks, nil
}

func (e *Ethereum) StoreBlock(blockNumber int64, block any) error {
	// Serialize the block to JSON
	blockJSON, err := json.Marshal(block)
	if err != nil {
		return fmt.Errorf("failed to marshal block: %w", err)
	}

	// Define the Redis hash key for storing blocks
	hashKey := "ethereum"

	// Use the blockNumber as the field in the Redis hash
	err = e.Factory.Data.RB.HSet(e.Factory.Ctx, hashKey, fmt.Sprintf("%d", blockNumber), blockJSON).Err()
	if err != nil {
		return fmt.Errorf("failed to store block in Redis hash: %w", err)
	}
	return nil
}
