package loopring

import (
	"encoding/json"
	"fmt"
)

// InsertBlock inserts a block into the database.
func (l *Loopring) InsertBlock(block *Block) error {
	query := `
        INSERT INTO blocks ( block_id, block_size, created, tx_hash, transactions)
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
