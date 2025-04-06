package loopring

import (
	"encoding/json"
	"fmt"
)

func (l *Loopring) HelloPeers() error {
	// Fetch the highest block ID in the database
	highestBlock := l.blockHeight()
	if highestBlock == 0 {
		return fmt.Errorf("no blocks found in the database")
	}

	// Initialize a map to track unique IDs
	uniqueIDs := make(map[int64]struct{})

	// Process each block up to the highest block
	for blockID := int64(1); blockID <= highestBlock; blockID++ {
		// Fetch the block data from the database
		var blockJSON string
		err := l.Factory.Db.QueryRow(`SELECT transactions FROM loopring WHERE block_id = $1`, blockID).Scan(&blockJSON)
		if err != nil {
			return fmt.Errorf("failed to fetch transactions for block %d: %w", blockID, err)
		}

		// Unmarshal the transactions
		var transactions []Transaction
		if err := json.Unmarshal([]byte(blockJSON), &transactions); err != nil {
			return fmt.Errorf("failed to unmarshal transactions for block %d: %w", blockID, err)
		}

		// Extract unique IDs from the transactions
		for _, tx := range transactions {
			if tx.From != 0 {
				uniqueIDs[tx.From] = struct{}{}
			}
			if tx.To != 0 {
				uniqueIDs[tx.To] = struct{}{}
			}
		}
	}

	// Get the total number of unique IDs
	count := len(uniqueIDs)
	fmt.Printf("Processing %d unique IDs...\n", count)

	// Call HelloUniverse for each unique ID and log the countdown
	for id := range uniqueIDs {
		idStr := fmt.Sprintf("%d", id)
		l.Factory.Peer.HelloUniverse(idStr)
		count--
		fmt.Printf("%d\n", count)
	}

	fmt.Println("All HelloUniverse calls completed.")
	return nil
}
