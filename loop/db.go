package loop

import (
	"encoding/json"
	"fmt"
	"strings"
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
	const batchSize = 1000 // Tune this for your DB
	l.Factory.Rw.RLock()
	defer l.Factory.Rw.RUnlock()

	coords := make([]Coordinate, 0, len(l.Map))
	txs := make([]*Tx, 0, len(l.Map))
	for coord, tx := range l.Map {
		coords = append(coords, coord)
		txs = append(txs, tx)
	}

	for i := 0; i < len(coords); i += batchSize {
		end := i + batchSize
		if end > len(coords) {
			end = len(coords)
		}

		valueStrings := make([]string, 0, end-i)
		valueArgs := make([]any, 0, (end-i)*9)
		argIdx := 1

		for j := i; j < end; j++ {
			txJSON, err := json.Marshal(txs[j])
			if err != nil {
				return fmt.Errorf("failed to marshal tx: %w", err)
			}
			valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
				argIdx, argIdx+1, argIdx+2, argIdx+3, argIdx+4, argIdx+5, argIdx+6, argIdx+7, argIdx+8))
			valueArgs = append(valueArgs,
				coords[j].Year, coords[j].Month, coords[j].Day, coords[j].Hour, coords[j].Minute,
				coords[j].Second, coords[j].Millisecond, coords[j].Index, txJSON,
			)
			argIdx += 9
		}

		query := `
            INSERT INTO tx (
                year, month, day, hour, minute, second, millisecond, index, tx
            ) VALUES ` + strings.Join(valueStrings, ",") + `
            ON CONFLICT (year, month, day, hour, minute, second, millisecond, index)
            DO UPDATE SET tx = EXCLUDED.tx;
        `

		tx, err := l.Factory.Db.BeginTx(l.Factory.Ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		_, err = tx.ExecContext(l.Factory.Ctx, query, valueArgs...)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to batch insert tx: %w", err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
	}
	return nil
}
