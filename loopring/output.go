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

	// Initialize a map to track unique accountId -> address mappings
	uniqueAccounts := make(map[int64]string)

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

		// Extract unique accountId -> address mappings from the transactions
		for _, tx := range transactions {
			// Only map the To accountId to its ToAddress
			if tx.To != 0 && tx.ToAddress != "" {
				uniqueAccounts[tx.To] = tx.ToAddress
			}
		}
	}

	// Get the total number of unique accounts
	count := len(uniqueAccounts)
	fmt.Printf("Processing %d unique accounts...\n", count)

	// Call HelloUniverse for each unique accountId and log the countdown
	for accountId, address := range uniqueAccounts {
		idStr := fmt.Sprintf("%d", accountId)
		fmt.Printf("Calling HelloUniverse for accountId: %s, address: %s\n", idStr, address)
		l.Factory.Peer.HelloUniverse(idStr)
		count--
		fmt.Printf("%d\n", count)
	}

	fmt.Println("All HelloUniverse calls completed.")
	return nil
}
