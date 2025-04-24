package loopring

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/log"
	"github.com/zachklingbeil/block/circuit"
)

func (l *Loopring) Loop() error {
	past, distance, err := l.Distance()
	if err != nil {
		log.Error("Failed to calculate block distance: %v", err)
		return err
	}

	for blockNumber := past + distance; blockNumber > past; blockNumber-- {
		if err := l.BlockByBlock(blockNumber); err != nil {
			log.Error("Error processing block %d: %v", blockNumber, err)
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
	blockJSONs, err := l.Factory.Redis.SMembers(l.Factory.Ctx, "blocks").Result()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve blocks from Redis: %w", err)
	}
	past := int64(0)
	for _, blockJSON := range blockJSONs {
		var block circuit.Raw
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

func (l *Loopring) Distance() (int64, int64, error) {
	current := l.currentBlock()
	past, err := l.getHistory()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get highest block from Redis: %w", err)
	}

	distance := current - past
	if distance > 0 {
		return past, distance, nil
	}
	return past, 0, nil
}
