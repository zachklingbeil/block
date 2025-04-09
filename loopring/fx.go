package loopring

import (
	"fmt"
	"time"
)

type Output struct {
	Number    int64
	Size      int64
	Timestamp int64
	Coords    string
}

// ProcessInputs converts a slice of Block into a slice of Output
func (l *Loopring) ProcessInputs(b []Block) []Output {
	blocks := make([]Output, len(b))

	for i, block := range b {
		blocks[i] = fx(block)
	}
	return blocks
}

// fx processes a single Block into a Output
func fx(block Block) Output {
	t := time.UnixMilli(block.Created)

	// Format the timestamp directly into a string representation of Coordinates
	formattedCoords := fmt.Sprintf("%d.%d.%d.%d.%d.%d.%d",
		t.Year()-2015,      // 0-based year
		int(t.Month()),     // Month
		t.Day(),            // Date of the month (1-31)
		t.Hour(),           // Hour
		t.Minute(),         // Minute
		t.Second(),         // Second
		t.Nanosecond()/1e6) // Millisecond as int64, uncapped

	// Return the Output with the formatted coordinates
	return Output{
		Coords:    formattedCoords,
		Number:    block.Number,
		Size:      block.Size,
		Timestamp: block.Created,
	}
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

	// Create a slice of Block
	var b []Block
	for rows.Next() {
		var block Block
		if err := rows.Scan(&block.Created, &block.Number, &block.Size); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}
		b = append(b, block)
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
