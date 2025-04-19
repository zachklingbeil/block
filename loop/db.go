package loop

import (
	"fmt"
)

func (l *Loopring) CreateTxTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS tx (
        coordinate JSONB NOT NULL,
        tx JSONB NOT NULL,
        PRIMARY KEY (coordinate)
    );
    `
	_, err := l.Factory.Db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}
	return nil
}
