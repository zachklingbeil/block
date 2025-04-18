package loop

import (
	"encoding/json"
	"fmt"
	"maps"
	"strings"
)

func (l *Loopring) CreateTxTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS tx (
        block BIGINT NOT NULL,
        year SMALLINT NOT NULL,
        month SMALLINT NOT NULL,
        day SMALLINT NOT NULL,
        hour SMALLINT NOT NULL,
        minute SMALLINT NOT NULL,
        second SMALLINT NOT NULL,
        millisecond SMALLINT NOT NULL,
        index SMALLINT NOT NULL,
        tx JSONB NOT NULL,
        PRIMARY KEY (block, year, month, day, hour, minute, second, millisecond, index)
    );
    `
	_, err := l.Factory.Db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}
	return nil
}

func (l *Loopring) SaveMap() error {
	const batchSize = 1000

	l.Factory.Rw.RLock()
	mapCopy := make(map[Coordinate]*Tx, len(l.Map))
	maps.Copy(mapCopy, l.Map)
	l.Factory.Rw.RUnlock()

	queryTemplate := `
        INSERT INTO tx (
            block, year, month, day, hour, minute, second, millisecond, index, tx
        ) VALUES %s
        ON CONFLICT (block, year, month, day, hour, minute, second, millisecond, index)
        DO UPDATE SET tx = EXCLUDED.tx;
    `

	coords := make([]Coordinate, 0, len(mapCopy))
	txs := make([]*Tx, 0, len(mapCopy))
	for coord, tx := range mapCopy {
		coords = append(coords, coord)
		txs = append(txs, tx)
	}

	for i := 0; i < len(coords); i += batchSize {
		end := min(i+batchSize, len(coords))

		valueStrings := make([]string, 0, end-i)
		valueArgs := make([]any, 0, (end-i)*10)
		argIdx := 1

		for j := i; j < end; j++ {
			txJSON, err := json.Marshal(txs[j])
			if err != nil {
				return fmt.Errorf("failed to marshal tx: %w", err)
			}
			valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
				argIdx, argIdx+1, argIdx+2, argIdx+3, argIdx+4, argIdx+5, argIdx+6, argIdx+7, argIdx+8, argIdx+9))
			valueArgs = append(valueArgs,
				coords[j].Block, coords[j].Year, coords[j].Month, coords[j].Day, coords[j].Hour, coords[j].Minute,
				coords[j].Second, coords[j].Millisecond, coords[j].Index, txJSON,
			)
			argIdx += 10
		}

		query := fmt.Sprintf(queryTemplate, strings.Join(valueStrings, ","))

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
