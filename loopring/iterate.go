// package for methods used to while iterating to absolute data
package loopring

import (
	"encoding/json"
	"fmt"
	"os"
)

// InsertBlock inserts a block into the database.
func (l *Loopring) InsertBlock(in *BlockIn) error {
	query := `
        INSERT INTO loopring (block_id, block_size, created, tx_hash, transactions)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (created) DO NOTHING
    `
	transactions, err := json.Marshal(in.Transactions)
	if err != nil {
		return fmt.Errorf("failed to marshal transactions: %w", err)
	}

	if _, err := l.Factory.Db.Exec(query, in.Number, in.Size, in.Created, in.TxHash, transactions); err != nil {
		return fmt.Errorf("failed to insert in into database: %w", err)
	}
	l.Factory.Json.Print(in.Number)
	return nil
}

// QualityControl checks if each block in the database has transactions.
// If a block does not have transactions, it fetches the block and updates the database.
func (l *Loopring) QualityControl() error {
	query := `SELECT block_id, transactions FROM loopring`

	rows, err := l.Factory.Db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query blocks from the database: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var blockID int64
		var transactionsJSON []byte

		if err := rows.Scan(&blockID, &transactionsJSON); err != nil {
			return fmt.Errorf("failed to scan block row: %w", err)
		}

		// Check if the transactions slice is empty
		if len(transactionsJSON) == 0 || string(transactionsJSON) == "[]" {
			fmt.Printf("BlockIn %d has no transactions. Fetching block data...\n", blockID)
			if err := l.GetBlock(int(blockID)); err != nil {
				fmt.Printf("Failed to fetch block %d: %v\n", blockID, err)
				continue
			}
			fmt.Printf("Successfully updated block %d with transactions.\n", blockID)
		}
	}
	return nil
}

// LoadBlocks queries the loopring table, processes the data, and inserts Blocks into the coords table
func (l *Loopring) LoadBlocks() error {
	// Query the loopring table to fetch data
	query := `
        SELECT created, block_id, block_size
        FROM loopring
    `
	rows, err := l.Factory.Db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query loopring table: %w", err)
	}
	defer rows.Close()

	// Create a slice of BlockIn
	var b []BlockIn
	for rows.Next() {
		var in BlockIn
		if err := rows.Scan(&in.Created, &in.Number, &in.Size); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}
		b = append(b, in)
	}

	// Process the b into Blocks
	blocks := l.ProcessInputs(b)

	// Insert the Blocks into the coords table
	for _, block := range blocks {
		if err := l.InsertBlockToCoords(&block); err != nil {
			return fmt.Errorf("failed to insert block into coords table: %w", err)
		}
	}

	return nil
}

// InsertBlockToCoords inserts a block into the coords table
func (l *Loopring) InsertBlockToCoords(o *Output) error {
	query := `
        INSERT INTO coords (block_id, block_size, created, coords)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (coords) DO NOTHING
    `
	if _, err := l.Factory.Db.Exec(query, o.Number, o.Size, o.Timestamp, o.Coords); err != nil {
		return fmt.Errorf("failed to insert block into coords table: %w", err)
	}
	return nil
}

func (l *Loopring) CreateCoordsTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS coords (
            block_id BIGINT NOT NULL,
            block_size BIGINT NOT NULL,
            created BIGINT NOT NULL,
            coords TEXT NOT NULL,
            PRIMARY KEY (coords) -- Use coords as the primary key
        )
    `
	if _, err := l.Factory.Db.Exec(query); err != nil {
		return fmt.Errorf("failed to create coords table: %w", err)
	}
	return nil
}

func (l *Loopring) OutputCoordsAsJSON() error {
	query := `
        SELECT block_id, block_size, created, coords
        FROM coords
    `
	rows, err := l.Factory.Db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query coords table: %w", err)
	}
	defer rows.Close()

	// Create a slice to hold the results
	var results []Output
	for rows.Next() {
		var output Output
		if err := rows.Scan(&output.Number, &output.Size, &output.Timestamp, &output.Coords); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, output)
	}

	// Convert the results to JSON
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results to JSON: %w", err)
	}

	// Write JSON to a file or print to console
	file, err := os.Create("coords.json")
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(jsonData); err != nil {
		return fmt.Errorf("failed to write JSON to file: %w", err)
	}

	fmt.Println("Coords table exported to coords.json")
	return nil
}

func (l *Loopring) OutputLoopringAsJSON() error {
	query := `
        SELECT block_id, block_size, created, tx_hash, transactions
        FROM loopring
    `
	rows, err := l.Factory.Db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query loopring table: %w", err)
	}
	defer rows.Close()

	// Create a slice to hold the results
	var results []map[string]interface{}
	for rows.Next() {
		var blockID int64
		var blockSize int64
		var created int64
		var txHash string
		var transactionsJSON []byte

		if err := rows.Scan(&blockID, &blockSize, &created, &txHash, &transactionsJSON); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		// Store the data in a map to preserve all fields
		blockData := map[string]interface{}{
			"block_id":     blockID,
			"block_size":   blockSize,
			"created":      created,
			"tx_hash":      txHash,
			"transactions": json.RawMessage(transactionsJSON), // Preserve raw JSON
		}

		results = append(results, blockData)
	}

	// Convert the results to JSON
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results to JSON: %w", err)
	}

	// Write JSON to a file or print to console
	file, err := os.Create("loopring.json")
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(jsonData); err != nil {
		return fmt.Errorf("failed to write JSON to file: %w", err)
	}

	fmt.Println("Loopring table exported to loopring.json")
	return nil
}
