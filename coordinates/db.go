package coordinates

import (
	"encoding/json"
	"fmt"
)

// Create the transactions table in the database
func (c *Coordinates) CreateTransactionsTable() error {
	createTableQuery := `
        CREATE TABLE IF NOT EXISTS transactions (
            year SMALLINT NOT NULL,
            month SMALLINT NOT NULL,
            day SMALLINT NOT NULL,
            hour SMALLINT NOT NULL,
            minute SMALLINT NOT NULL,
            second SMALLINT NOT NULL,
            millisecond SMALLINT NOT NULL,
            index SMALLINT NOT NULL,
            data JSONB NOT NULL,
            PRIMARY KEY (year, month, day, hour, minute, second, millisecond, index)
        );
    `
	_, err := c.Factory.Db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}
	return nil
}

// Fetch all transactions from the database and return a map
func (c *Coordinates) FetchTransactions() (map[string]any, error) {
	query := `
        SELECT year, month, day, hour, minute, second, millisecond, index, data
        FROM transactions;
    `
	rows, err := c.Factory.Db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	defer rows.Close()

	transactions := make(map[string]any)

	for rows.Next() {
		var coord Coord
		var dataJSON []byte

		err := rows.Scan(
			&coord.Year,
			&coord.Month,
			&coord.Day,
			&coord.Hour,
			&coord.Minute,
			&coord.Second,
			&coord.Millisecond,
			&coord.Index,
			&dataJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Create a composite key as a string
		key := fmt.Sprintf("%d-%d-%d-%d-%d-%d-%d-%d",
			coord.Year,
			coord.Month,
			coord.Day,
			coord.Hour,
			coord.Minute,
			coord.Second,
			coord.Millisecond,
			coord.Index,
		)

		// Unmarshal the JSONB data
		var data any
		if err := json.Unmarshal(dataJSON, &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal transaction data: %w", err)
		}

		// Add the key-value pair to the map
		transactions[key] = data
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return transactions, nil
}

// Insert a transaction into the database
func (c *Coordinates) InsertTransaction(coord Coord, data any) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction data: %w", err)
	}

	query := `
        INSERT INTO transactions (year, month, day, hour, minute, second, millisecond, index, data)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (year, month, day, hour, minute, second, millisecond, index) DO UPDATE
        SET data = EXCLUDED.data;
    `
	_, err = c.Factory.Db.Exec(query,
		coord.Year,
		coord.Month,
		coord.Day,
		coord.Hour,
		coord.Minute,
		coord.Second,
		coord.Millisecond,
		coord.Index,
		dataJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}
	return nil
}
