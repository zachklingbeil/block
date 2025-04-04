package loopring

import "fmt"

func (l *Loopring) ExtractPeerInfo() error {
	var blocks []Block
	if err := l.Factory.DiskToMem("loopring", &blocks); err != nil {
		return fmt.Errorf("failed to load database: %w", err)
	}

	// Extract IDs without addresses
	noAddress := make(map[int64]struct{})

	for _, block := range blocks {
		for _, tx := range block.Transactions {
			// Check if the address already exists in Peer.Map
			if tx.ToAddress != "" {
				formattedAddress := l.Peers.FormatAddress(tx.ToAddress)
				peer, exists := l.Peers.Map[formattedAddress]
				if !exists || peer == nil {
					// Address exists but no valid Peer, skip further processing
					continue
				}
			}

			// Check if the ID is in Peer.Map
			if tx.To != 0 {
				idKey := fmt.Sprintf("%d", tx.To)
				peer, exists := l.Peers.Map[idKey]
				if !exists || peer == nil {
					noAddress[tx.To] = struct{}{}
				}
			}
		}
	}

	// Convert noAddress map to a slice and fetch missing addresses
	missingAccounts := intsToSlice(noAddress)
	if err := l.FetchMissingAddresses(missingAccounts); err != nil {
		return fmt.Errorf("failed to fetch missing addresses: %w", err)
	}

	return nil
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

// Helper function to convert map keys to a slice of int64
func intsToSlice(m map[int64]struct{}) []int64 {
	keys := make([]int64, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
