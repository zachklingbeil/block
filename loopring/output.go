package loopring

import "fmt"

func (l *Loopring) ExtractPeerInfo() ([]string, []int64, []int64, error) {
	var blocks []Block
	if err := l.Factory.DiskToMem("loopring", &blocks); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load database: %w", err)
	}

	// Extract unique addresses, account IDs, and IDs without addresses
	address := make(map[string]struct{})
	id := make(map[int64]struct{})
	noAddress := make(map[int64]struct{})

	for _, block := range blocks {
		for _, tx := range block.Transactions {
			// Check if the address already exists in Peer.Map
			if tx.ToAddress != "" {
				if _, exists := l.Peers.Map[tx.ToAddress]; !exists {
					address[tx.ToAddress] = struct{}{}
				}
			} else if tx.To != 0 {
				// Check if the ID without an address already exists in Peer.Map
				if _, exists := l.Peers.Map[fmt.Sprintf("%d", tx.To)]; !exists {
					noAddress[tx.To] = struct{}{}
				}
			}

			// Check if the "From" account ID already exists in Peer.Map
			if tx.From != 0 {
				if _, exists := l.Peers.Map[fmt.Sprintf("%d", tx.From)]; !exists {
					id[tx.From] = struct{}{}
				}
			}

			// Check if the "To" account ID already exists in Peer.Map
			if tx.To != 0 {
				if _, exists := l.Peers.Map[fmt.Sprintf("%d", tx.To)]; !exists {
					id[tx.To] = struct{}{}
				}
			}
		}
	}
	return stringsToSlice(address), intsToSlice(id), intsToSlice(noAddress), nil
}

func (l *Loopring) FetchMissingAddresses(missingAccounts []int64) error {
	totalAccounts := len(missingAccounts)

	for i, accountID := range missingAccounts {
		l.Factory.Json.Print(fmt.Sprintf("%d/%d", i+1, totalAccounts))

		peer, err := l.Peers.FetchLoopringAddress(accountID)
		if err != nil {
			fmt.Printf("Failed to fetch address for account ID %d: %v\n", accountID, err)
			continue
		}

		if err := l.Peers.Update(peer); err != nil {
			fmt.Printf("Failed to insert or update address for account ID %d: %v\n", accountID, err)
		}
	}
	return nil
}

// Helper function to convert map keys to a slice of strings
func stringsToSlice(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

// Helper function to convert map keys to a slice of int64
func intsToSlice(m map[int64]struct{}) []int64 {
	keys := make([]int64, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
