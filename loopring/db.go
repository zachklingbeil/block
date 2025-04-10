package loopring

import (
	"encoding/json"
	"fmt"
)

func (l *Loopring) CreateTable() error {
	createTableQuery := `
        CREATE TABLE IF NOT EXISTS loopring (
            block BIGINT PRIMARY KEY,
            tx JSONB NOT NULL
        );
    `
	_, err := l.Factory.Db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}

func (l *Loopring) StoreTransactions(blockNumber int64, transactions []any) error {
	txJSON, err := json.Marshal(transactions)
	if err != nil {
		return fmt.Errorf("failed to marshal transactions: %w", err)
	}

	query := `
        INSERT INTO loopring (block, tx)
        VALUES ($1, $2)
        ON CONFLICT (block) DO UPDATE
        SET tx = EXCLUDED.tx;
    `
	_, err = l.Factory.Db.Exec(query, blockNumber, txJSON)
	if err != nil {
		return fmt.Errorf("failed to store transactions: %w", err)
	}
	return nil
}
