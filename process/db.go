package process

import (
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
