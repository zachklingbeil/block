package loopring

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/log"
	"github.com/zachklingbeil/block/circuit"
	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Circuit *circuit.Circuit
}

func Connect(factory *factory.Factory, circuit *circuit.Circuit) *Loopring {
	loop := &Loopring{
		Factory: factory,
		Circuit: circuit,
	}
	go loop.Listen()
	return loop
}

func (l *Loopring) BlockByBlock(blockNumber int64) {
	input := l.FetchBlock(blockNumber)
	transactions, block := l.Circuit.Coordinates(input)
	txs := l.ProcessBlock(transactions)
	block.Ones = txs
	blockJSON, _ := json.Marshal(block)
	l.Factory.Redis.SAdd(l.Factory.Ctx, "blocks", blockJSON).Err()
	fmt.Printf("%d\n", blockNumber)
}

func (l *Loopring) Loop() {
	past, distance := l.Distance()
	for blockNumber := past + distance; blockNumber > past; blockNumber-- {
		l.BlockByBlock(blockNumber)
	}
}

func (l *Loopring) Distance() (int64, int64) {
	current := l.currentBlock()
	past, _ := l.getHistory()

	distance := current - past
	if distance > 0 {
		return past, distance
	}
	return past, 0
}

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

// getHistory retrieves the lowest block number from the Redis set
func (l *Loopring) getHistory() (int64, error) {
	blockJSONs, err := l.Factory.Redis.SMembers(l.Factory.Ctx, "blocks").Result()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve blocks from Redis: %w", err)
	}
	past := int64(^uint64(0) >> 1)
	for _, blockJSON := range blockJSONs {
		var block circuit.Block
		if err := json.Unmarshal([]byte(blockJSON), &block); err != nil {
			log.Error("Failed to deserialize block JSON: %v", err)
			continue
		}
		if block.Number < past {
			past = block.Number
		}
	}
	if past == int64(^uint64(0)>>1) {
		return 0, nil
	}
	return past, nil
}
