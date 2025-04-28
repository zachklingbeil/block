package loopring

import (
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/zachklingbeil/block/circuit"
	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory      *factory.Factory
	Circuit      *circuit.Circuit
	CurrentBlock int64
}

func Connect(factory *factory.Factory, circuit *circuit.Circuit) *Loopring {
	pb := factory.State.Get("processedBlocks")
	loop := &Loopring{
		Factory:      factory,
		Circuit:      circuit,
		CurrentBlock: pb.(int64),
	}
	go loop.Listen()
	return loop
}

func (l *Loopring) BlockByBlock(blockNumber int64) []byte {
	input := l.FetchBlock(blockNumber)
	transactions, block := l.Circuit.Coordinates(input)
	txs := l.ProcessBlock(transactions)
	block.Ones = txs
	blockJSON, _ := json.Marshal(block)
	return blockJSON
}

func (l *Loopring) StoreBlock(blockNumber int64, blockJSON []byte) error {
	score := float64(blockNumber)
	err := l.Factory.Data.RB.ZAdd(l.Factory.Ctx, "blocks", redis.Z{
		Score:  score,
		Member: blockJSON,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to store block in Redis: %w", err)
	}
	l.Factory.State.Add("blockHeight", l.CurrentBlock)
	return nil
}

// getHistory retrieves the lowest block number from the Redis set
func (l *Loopring) getHistory() (int64, error) {
	result, err := l.Factory.Data.RB.ZRevRangeWithScores(l.Factory.Ctx, "blocks", 0, 0).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve blocks from Redis: %w", err)
	}

	if len(result) == 0 {
		return 0, nil
	}

	highestBlock := int64(result[0].Score)
	l.Factory.State.Add("processedBlocks", highestBlock)
	return highestBlock, nil
}

func (l *Loopring) Loop() {
	past, distance := l.Distance()
	for blockNumber := past + distance; blockNumber > past; blockNumber-- {
		blockJSON := l.BlockByBlock(blockNumber)
		if err := l.StoreBlock(blockNumber, blockJSON); err != nil {
			fmt.Printf("Error storing block %d: %v\n", blockNumber, err)
		}
	}
}
