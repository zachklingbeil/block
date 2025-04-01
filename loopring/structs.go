package loopring

import (
	"encoding/json"
	"fmt"
)

// CurrentBlock fetches the latest block Number from the Loopring API.
func (l *Loopring) CurrentBlock() (int64, error) {
	response, err := l.Factory.Json.In("https://api3.loopring.io/api/v3/block/getBlock", "")
	if err != nil {
		return 0, fmt.Errorf("failed to fetch the latest block data: %w", err)
	}

	var block Block
	if err := json.Unmarshal(response, &block); err != nil {
		return 0, fmt.Errorf("failed to parse block data: %w", err)
	}
	l.Factory.Json.Print(block.Number)
	return block.Number, nil
}

// Helper function to read transactions from the map for a given block number.
func (l *Loopring) Read(blockNumber int64) (*Block, bool) {
	l.Factory.Mu.Lock()
	defer l.Factory.Mu.Unlock()
	block, exists := l.Map[blockNumber]
	if !exists {
		return nil, false
	}
	return block, true
}

// Helper function to update the map with transactions for a given block number.
func (l *Loopring) Write(block *Block) {
	l.Factory.Mu.Lock()
	defer l.Factory.Mu.Unlock()
	l.Map[block.Number] = block
}

func (l *Loopring) LoadBlocks() error {
	query := `SELECT created, block_id, block_size, tx_hash, transactions FROM blocks`

	rows, err := l.Db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query blocks from database: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var block Block
		var transactionsJSON []byte

		if err := rows.Scan(&block.Created, &block.Number, &block.Size, &block.TxHash, &transactionsJSON); err != nil {
			return fmt.Errorf("failed to scan block row: %w", err)
		}

		if err := json.Unmarshal(transactionsJSON, &block.Transactions); err != nil {
			return fmt.Errorf("failed to unmarshal transactions: %w", err)
		}

		// Write to the in-memory map
		l.Write(&block)
	}

	return nil
}
