package loopring

import (
	"encoding/json"
	"fmt"
)

// CreateTable ensures the blocks table exists in the database.
func (l *Loopring) CreateTable() error {
	query := `
		  CREATE TABLE IF NOT EXISTS blocks (
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
        INSERT INTO blocks (block_id, block_size, created, tx_hash, transactions)
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

// DiskToMem loads all blocks from the database into the Loopring struct's Blocks field.
func (l *Loopring) DiskToMem() error {
	query := `
        SELECT block_id, block_size, created, tx_hash, transactions
        FROM blocks
    `

	rows, err := l.Db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query blocks from database: %w", err)
	}
	defer rows.Close()

	var blocks []Block
	for rows.Next() {
		var block Block
		var transactionsJSON []byte

		// Scan block data
		if err := rows.Scan(&block.Number, &block.Size, &block.Created, &block.TxHash, &transactionsJSON); err != nil {
			return fmt.Errorf("failed to scan block data: %w", err)
		}

		// Unmarshal transactions
		if err := json.Unmarshal(transactionsJSON, &block.Transactions); err != nil {
			return fmt.Errorf("failed to unmarshal transactions: %w", err)
		}

		blocks = append(blocks, block)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over rows: %w", err)
	}

	l.Blocks = blocks
	return nil
}

// EnsureTransactions checks if each block in the database has transactions.
// If a block does not have transactions, it fetches the block and updates the database.
func (l *Loopring) EnsureTransactions() error {
	query := `SELECT block_id, transactions FROM blocks`

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

		// Check if the transactions slice is empty
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
