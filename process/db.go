package process

import (
	"encoding/json"
	"fmt"
)

func (p *Process) CreateTxTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS tx (
		year SMALLINT NOT NULL,
		month SMALLINT NOT NULL,
		day SMALLINT NOT NULL,
		hour SMALLINT NOT NULL,
		minute SMALLINT NOT NULL,
		second SMALLINT NOT NULL,
		millisecond SMALLINT NOT NULL,
		index SMALLINT NOT NULL,
		tx JSONB NOT NULL,
		PRIMARY KEY (year, month, day, hour, minute, second, millisecond, index)
		);
		`
	_, err := p.Factory.Db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}
	return nil
}

func (p *Process) LoadRecentBlocks(limit int) error {
	query := `
        SELECT block, tx
        FROM loopring
        ORDER BY block DESC
        LIMIT $1;
    `
	rows, err := p.Factory.Db.Query(query, limit)
	if err != nil {
		return fmt.Errorf("failed to query loopring table: %w", err)
	}
	defer rows.Close()

	var rawTxs []RawTx
	for rows.Next() {
		var txJSON []byte

		if err := rows.Scan(new(int64), &txJSON); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		var transactions []RawTx
		if err := json.Unmarshal(txJSON, &transactions); err != nil {
			return fmt.Errorf("failed to unmarshal transactions: %w", err)
		}
		rawTxs = append(rawTxs, transactions...)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over rows: %w", err)
	}

	p.RawTxs = rawTxs
	fmt.Printf("Loaded %d transactions from the database\n", len(p.RawTxs))
	return nil
}
