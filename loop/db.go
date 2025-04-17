package loop

import (
	"encoding/json"
	"fmt"
)

func (l *Loopring) CreateTxTable() error {
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
	_, err := l.Factory.Db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}
	return nil
}

func (l *Loopring) SaveMap() error {
	l.Factory.Rw.RLock()
	defer l.Factory.Rw.RUnlock()
	query := `
        INSERT INTO tx (
            year, month, day, hour, minute, second, millisecond, index, tx
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9
        )
        ON CONFLICT (year, month, day, hour, minute, second, millisecond, index)
        DO UPDATE SET tx = EXCLUDED.tx;
    `
	for coord, tx := range l.Map {
		txJSON, err := json.Marshal(tx)
		if err != nil {
			return fmt.Errorf("failed to marshal tx: %w", err)
		}
		_, err = l.Factory.Db.ExecContext(
			l.Factory.Ctx,
			query,
			coord.Year, coord.Month, coord.Day, coord.Hour, coord.Minute,
			coord.Second, coord.Millisecond, coord.Index, txJSON,
		)
		if err != nil {
			return fmt.Errorf("failed to insert tx: %w", err)
		}
	}
	return nil
}
