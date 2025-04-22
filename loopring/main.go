package loopring

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/log"
	"github.com/zachklingbeil/factory"
	"github.com/zachklingbeil/factory/fx"
)

type Loopring struct {
	Factory *factory.Factory
	Block   *Block
}

type Block struct {
	Number       int64   `json:"blockId"`
	Timestamp    int64   `json:"createdAt"`
	Size         int64   `json:"blockSize"`
	Coord        fx.Zero `json:"coordinate"`
	Transactions []any   `json:"transactions"`
}

func Connect(factory *factory.Factory) *Loopring {
	loop := &Loopring{
		Factory: factory,
		Block:   &Block{},
	}
	return loop
}

func (l *Loopring) Loop() error {
	past, _ := l.Distance()
	current := l.currentBlock()

	for blockNumber := past + 1; blockNumber <= current; blockNumber++ {
		if err := l.FetchBlock(blockNumber); err != nil {
			log.Error("Failed to fetch block %d: %v", blockNumber, err)
			continue
		}
		l.Coordinates()
		if err := l.Index(); err != nil {
			log.Error("Failed to index block %d: %v", blockNumber, err)
			continue
		}
		l.ProcessTransactions()
		blockJSON, err := json.Marshal(l.Block)
		if err != nil {
			log.Error("Failed to serialize block %d: %v", blockNumber, err)
			continue
		}
		err = l.Factory.Db.Rdb.SAdd(l.Factory.Ctx, "blocks", blockJSON).Err()
		if err != nil {
			log.Error("Failed to store block %d in Redis: %v", blockNumber, err)
		} else {
			fmt.Printf("%d\n", blockNumber)
		}
	}

	return nil
}

// Simplified GetCurrentBlockNumber
func (l *Loopring) currentBlock() int64 {
	data, err := l.Factory.Json.In("https://api3.loopring.io/api/v3/block/getBlock", "")
	if err != nil {
		fmt.Printf("Failed to fetch block data: %v\n", err)
		return 0
	}
	var block struct {
		Number int64 `json:"blockId"`
	}
	err = json.Unmarshal(data, &block)
	if err != nil {
		fmt.Printf("Failed to parse block data: %v\n", err)
		return 0
	}
	return block.Number
}

// getHistory retrieves the highest block number from the Redis set
func (l *Loopring) getHistory() (int64, error) {
	blockJSONs, err := l.Factory.Db.Rdb.SMembers(l.Factory.Ctx, "blocks").Result()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve blocks from Redis: %w", err)
	}
	past := int64(0)
	for _, blockJSON := range blockJSONs {
		var block Block
		if err := json.Unmarshal([]byte(blockJSON), &block); err != nil {
			log.Error("Failed to deserialize block JSON: %v", err)
			continue
		}
		if block.Number > past {
			past = block.Number
		}
	}
	return past, nil
}
func (l *Loopring) Distance() (int64, error) {
	current := l.currentBlock()
	past, err := l.getHistory()
	if err != nil {
		return 0, fmt.Errorf("failed to get highest block from Redis: %w", err)
	}

	if current > past {
		return current, nil
	}
	return past, nil
}
