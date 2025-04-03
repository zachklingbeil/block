package loopring

import (
	"encoding/json"
	"fmt"
)

// FetchBlocks fetches blocks sequentially from the last fetched block to the current block and stores them in the database.
func (l *Loopring) FetchBlocks() error {
	// Current block number - highest blockId in database
	response, err := l.Factory.Json.In("https://api3.loopring.io/api/v3/block/getBlock", "")
	if err != nil {
		return fmt.Errorf("failed to fetch the latest block data: %w", err)
	}

	var block Block
	if err := json.Unmarshal(response, &block); err != nil {
		return fmt.Errorf("failed to parse block data: %w", err)
	}
	currentBlock := block.Number

	query := `SELECT COALESCE(MAX(block_id), 0) FROM blocks`
	var blockHeight int64
	if err := l.Db.QueryRow(query).Scan(&blockHeight); err != nil {
		return fmt.Errorf("failed to fetch the highest block ID: %w", err)
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

// CreateTable ensures the blocks table exists in the database.
func (l *Loopring) CreateTable() error {
	query := `
		  CREATE TABLE IF NOT EXISTS loopring (
				block_id BIGINT PRIMARY KEY,
				block_size BIGINT NOT NULL,
				created BIGINT UNIQUE NOT NULL, -- Add UNIQUE constraint
				tx_hash TEXT NOT NULL,
				transactions JSONB NOT NULL
		  );
	 `

	if _, err := l.Db.Exec(query); err != nil {
		return fmt.Errorf("failed to create blocks table: %w", err)
	}
	return nil
}

// InsertBlock inserts a block into the database.
func (l *Loopring) InsertBlock(block *Block) error {
	query := `
        INSERT INTO loopring (block_id, block_size, created, tx_hash, transactions)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (created) DO NOTHING
    `

	transactions, err := json.Marshal(block.Transactions)
	if err != nil {
		return fmt.Errorf("failed to marshal transactions: %w", err)
	}

	if _, err := l.Db.Exec(query, block.Number, block.Size, block.Created, block.TxHash, transactions); err != nil {
		return fmt.Errorf("failed to insert block into database: %w", err)
	}
	l.Factory.Json.Print(block.Number)
	return nil
}

// If a block does not have transactions, QualityControl it fetches the block and updates the database.
func (l *Loopring) QualityControl() error {
	query := `SELECT block_id, transactions FROM loopring`
	rows, err := l.Db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query blocks from the database: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var blockID int64
		var transactionsJSON []byte

		if err := rows.Scan(&blockID, &transactionsJSON); err != nil {
			return fmt.Errorf("failed to scan block row: %w", err)
		}

		if len(transactionsJSON) == 0 || string(transactionsJSON) == "[]" {
			fmt.Printf("Block %d has no transactions. Fetching block data...\n", blockID)
			if err := l.GetBlock(int(blockID)); err != nil {
				fmt.Printf("Failed to fetch block %d: %v\n", blockID, err)
				continue
			}
			fmt.Printf("Successfully updated block %d with transactions.\n", blockID)
		}
	}
	return nil
}
