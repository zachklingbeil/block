package process

import (
	"encoding/json"
	"fmt"
)

func (p *Process) CreateTxTable() error {
	createTableQuery := `
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
	_, err := p.Factory.Db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}
	return nil
}

func (p *Process) LoadBlocks() error {
	query := `
        SELECT block, tx
        FROM loopring
        ORDER BY block ASC;
    `

	rows, err := p.Factory.Db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query loopring table: %w", err)
	}
	defer rows.Close()

	var blocks []Block
	for rows.Next() {
		var block Block
		var txJSON []byte

		if err := rows.Scan(&block.Number, &txJSON); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		if err := json.Unmarshal(txJSON, &block.Transactions); err != nil {
			return fmt.Errorf("failed to unmarshal transactions: %w", err)
		}

		blocks = append(blocks, block)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over rows: %w", err)
	}

	p.Blocks = blocks
	fmt.Printf("Loaded %d blocks from the database\n", len(p.Blocks))
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

	var blocks []Block
	for rows.Next() {
		var block Block
		var txJSON []byte

		if err := rows.Scan(&block.Number, &txJSON); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		if err := json.Unmarshal(txJSON, &block.Transactions); err != nil {
			return fmt.Errorf("failed to unmarshal transactions: %w", err)
		}

		blocks = append(blocks, block)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over rows: %w", err)
	}

	// Reverse the order of blocks to maintain ascending order
	for i, j := 0, len(blocks)-1; i < j; i, j = i+1, j-1 {
		blocks[i], blocks[j] = blocks[j], blocks[i]
	}

	p.Blocks = blocks
	fmt.Printf("Loaded %d blocks from the database\n", len(p.Blocks))
	return nil
}
