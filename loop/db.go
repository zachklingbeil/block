package loop

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

func (l *Loopring) fetchBlock(number int64) ([]byte, error) {
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/block/getBlock?id=%d", number)
	response, err := l.Factory.Json.In(url, "")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block %d: %w", number, err)
	}
	return response, nil
}
