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
	Raw     *Raw
}

func Connect(factory *factory.Factory) *Loopring {
	loop := &Loopring{
		Factory: factory,
	}
	return loop
}

func (l *Loopring) BlockByBlock(blockNumber int64) error {
	coord, transactions, err := l.FetchBlock(blockNumber)
	if err != nil {
		log.Error("Failed to fetch block %d: %v", blockNumber, err)
		return err
	}

	txs, err := l.ProcessTransactions(transactions)
	if err != nil {
		log.Error("Failed to process transactions for block %d: %v", blockNumber, err)
		return err
	}

	block := &Block{
		Coord:        coord,
		Transactions: txs,
	}

	blockJSON, err := json.Marshal(block)
	if err != nil {
		log.Error("Failed to serialize block %d: %v", blockNumber, err)
		return err
	}

	err = l.Factory.Db.Rdb.SAdd(l.Factory.Ctx, "blocks", blockJSON).Err()
	if err != nil {
		log.Error("Failed to store block %d in Redis: %v", blockNumber, err)
		return err
	}
	fmt.Printf("%d\n", blockNumber)
	return nil
}

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
	blockJSONs, err := l.Factory.Db.Rdb.SMembers(l.Factory.Ctx, "blocks").Result()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve blocks from Redis: %w", err)
	}
	past := int64(0)
	for _, blockJSON := range blockJSONs {
		var block Raw
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

func (l *Loopring) ToMap() map[fx.Zero][]Tx {
	if l.Block == nil {
		return nil
	}
	return map[fx.Zero][]Tx{
		l.Block.Coord: l.Block.Transactions,
	}
}
