package loopring

import (
	"encoding/json"
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
}

// NewLoopring initializes a new Loopring instance and ensures the database table exists.
func NewLoopring(factory *factory.Factory) (*Loopring, error) {
	loopring := &Loopring{
		Factory: factory,
	}

	if err := loopring.CreateTable(); err != nil {
		return nil, fmt.Errorf("failed to create blocks table: %w", err)
	}
	return loopring, nil
}

// FetchBlocks fetches blocks sequentially from the last fetched block to the current block and stores them in the database.
func (l *Loopring) FetchBlocks() error {
	// Fetch the current block number directly
	response, err := l.Factory.Json.In("https://api3.loopring.io/api/v3/block/getBlock", "")
	if err != nil {
		return fmt.Errorf("failed to fetch the latest block data: %w", err)
	}

	var block Block
	if err := json.Unmarshal(response, &block); err != nil {
		return fmt.Errorf("failed to parse block data: %w", err)
	}
	currentBlock := block.Number

	// Get the highest block ID from the database
	query := `SELECT COALESCE(MAX(block_id), 0) FROM loopring`
	var blockHeight int64
	if err := l.Factory.Db.QueryRow(query).Scan(&blockHeight); err != nil {
		return fmt.Errorf("failed to fetch the highest block ID: %w", err)
	}

	if blockHeight == currentBlock {
		fmt.Println("blockHeight = currentBlock")
		return nil
	}

	// Fetch and store each block sequentially
	for i := blockHeight + 1; i <= currentBlock; i++ {
		if err := l.GetBlock(int(i)); err != nil {
			fmt.Printf("Failed to fetch block %d: %v\n", i, err)
			continue
		}
	}
	l.QualityControl()
	return nil
}

// GetBlock fetches a block from the Loopring API and inserts it into the database.
func (l *Loopring) GetBlock(number int) error {
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/block/getBlock?id=%d", number)
	response, err := l.Factory.Json.In(url, "")
	if err != nil {
		return fmt.Errorf("failed to fetch block data for block number %d: %w", number, err)
	}

	var block Block
	if err := json.Unmarshal(response, &block); err != nil {
		return fmt.Errorf("failed to parse block data for block number %d: %w", number, err)
	}

	if err := l.InsertBlock(&block); err != nil {
		return fmt.Errorf("failed to insert block into database: %w", err)
	}
	return nil
}
