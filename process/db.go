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

	var rawTxs []any
	for rows.Next() {
		var txArray []byte
		if err := rows.Scan(new(int64), &txArray); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		var transactions []json.RawMessage
		if err := json.Unmarshal(txArray, &transactions); err != nil {
			return fmt.Errorf("failed to unmarshal transactions array: %w", err)
		}

		for _, tx := range transactions {
			rawTxs = append(rawTxs, tx)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over rows: %w", err)
	}
	p.RawTxs = rawTxs
	p.Counts["Input"] = len(p.RawTxs)
	return nil
}
