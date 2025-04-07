package loopring

import (
	"encoding/json"
	"fmt"
)

func (l *Loopring) HelloPeers() error {
	highestBlock := l.blockHeight()
	if highestBlock == 0 {
		return fmt.Errorf("no blocks found in the database")
	}

	// Initialize a map to track unique accountId -> address mappings
	uniqueAccounts := make(map[int64]string)

	// Process each block up to the highest block
	if err := l.collectUniqueAccounts(highestBlock, uniqueAccounts); err != nil {
		return err
	}

	// Count and log the number of missing addresses
	count := l.missingAddresses(uniqueAccounts)
	fmt.Printf("Need %d addresses.\n", count)

	// Fetch missing addresses
	if err := l.fetchMissingAddresses(uniqueAccounts, &count); err != nil {
		return err
	}

	// Create a slice of addresses and call HelloUniverse
	l.callHelloUniverse(uniqueAccounts)

	return nil
}

// Helper function to collect unique accounts from all blocks
func (l *Loopring) collectUniqueAccounts(highestBlock int64, uniqueAccounts map[int64]string) error {
	for blockID := int64(1); blockID <= highestBlock; blockID++ {
		// Fetch and unmarshal the block data
		var blockJSON string
		if err := l.Factory.Db.QueryRow(`SELECT transactions FROM loopring WHERE block_id = $1`, blockID).Scan(&blockJSON); err != nil {
			return fmt.Errorf("failed to fetch transactions for block %d: %w", blockID, err)
		}

		var transactions []Transaction
		if err := json.Unmarshal([]byte(blockJSON), &transactions); err != nil {
			return fmt.Errorf("failed to unmarshal transactions for block %d: %w", blockID, err)
		}

		// Extract unique accountId -> address mappings
		for _, tx := range transactions {
			if tx.To != 0 {
				if _, exists := uniqueAccounts[tx.To]; !exists {
					uniqueAccounts[tx.To] = tx.ToAddress
				}
			}
			if tx.From != 0 {
				if _, exists := uniqueAccounts[tx.From]; !exists {
					uniqueAccounts[tx.From] = ""
				}
			}
		}
	}
	return nil
}

// Helper function to count missing addresses
func (l *Loopring) missingAddresses(uniqueAccounts map[int64]string) int {
	count := 0
	for _, address := range uniqueAccounts {
		if address == "" {
			count++
		}
	}
	return count
}

// Helper function to fetch missing addresses
func (l *Loopring) fetchMissingAddresses(uniqueAccounts map[int64]string, count *int) error {
	for accountId, address := range uniqueAccounts {
		if address == "" {
			idStr := fmt.Sprintf("%d", accountId)
			fmt.Printf("%d\n", *count)
			fetchedAddress := l.Factory.Peer.NeedsAddress(idStr)
			if fetchedAddress != "" {
				uniqueAccounts[accountId] = fetchedAddress
				(*count)--
			} else {
				// Log the failure and continue instead of exiting
				fmt.Printf("Warning: Failed to fetch address for accountId: %s\n", idStr)
			}
		}
	}
	return nil
}

// Helper function to call HelloUniverse for all addresses
func (l *Loopring) callHelloUniverse(uniqueAccounts map[int64]string) {
	// Create a slice of addresses
	addresses := make([]string, 0, len(uniqueAccounts))
	for _, address := range uniqueAccounts {
		addresses = append(addresses, address)
	}

	// Get the new peers created by BatchCreatePeers
	newPeers := l.Factory.Peer.BatchCreatePeers(addresses)

	// Call HelloUniverse only for the new peers
	fmt.Printf("Saying hellouniverse to %d new peers...\n", len(newPeers))
	for i, address := range newPeers {
		l.Factory.Peer.HelloUniverse(address)
		fmt.Printf("%d\n", i+1)
	}
	fmt.Println("HelloUniverse")
}
