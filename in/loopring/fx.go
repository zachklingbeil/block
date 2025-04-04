package loopring

import (
	"encoding/json"
	"fmt"
)

// need to use the peers table, not loopring table.
func (l *Loopring) GetPeers() error {
	var blocks []Block

	if err := l.Factory.DiskToMem("loopring", &blocks); err != nil {
		return fmt.Errorf("failed to load database: %w", err)
	}

	peerMap, missingAddresses := l.PeerMap(blocks)
	if len(missingAddresses) > 0 {
		if err := l.FetchMissingAddresses(peerMap, missingAddresses); err != nil {
			return fmt.Errorf("failed to fetch missing addresses: %w", err)
		}
	}
	if err := l.MapToTable(peerMap); err != nil {
		return fmt.Errorf("failed to store address map: %w", err)
	}
	return nil
}

func (l *Loopring) PeerMap(blocks []Block) (map[int64]string, []int64) {
	peerMap := make(map[int64]string)
	missingAddresses := []int64{}

	for _, block := range blocks {
		for _, tx := range block.Transactions {
			if _, exists := peerMap[tx.From]; !exists {
				peerMap[tx.From] = tx.ToAddress
			}
			if _, exists := peerMap[tx.To]; !exists {
				peerMap[tx.To] = tx.ToAddress
			}
		}
	}

	for accountID, address := range peerMap {
		if address == "" {
			missingAddresses = append(missingAddresses, accountID)
		}
	}
	return peerMap, missingAddresses
}

func (l *Loopring) FetchMissingAddresses(peerMap map[int64]string, missingAccounts []int64) error {
	totalAccounts := len(missingAccounts)

	for i, accountID := range missingAccounts {
		l.Factory.Json.Print(fmt.Sprintf("%d/%d", i+1, totalAccounts))

		url := fmt.Sprintf("https://api3.loopring.io/api/v3/account?accountId=%d", accountID)
		response, err := l.Factory.Json.In(url, "")
		if err != nil {
			fmt.Printf("Failed to fetch address for account ID %d: %v\n", accountID, err)
			continue
		}

		var accountData struct {
			Owner string `json:"owner"`
		}
		if err := json.Unmarshal(response, &accountData); err != nil {
			fmt.Printf("Failed to parse address for account ID %d: %v\n", accountID, err)
			continue
		}
		peerMap[accountID] = accountData.Owner
	}
	return nil
}

func (l *Loopring) MapToTable(addressMap map[int64]string) error {
	query := `
    INSERT INTO peers (id, address)
    VALUES ($1, $2)
    ON CONFLICT (address) DO UPDATE
    SET id = EXCLUDED.id;
    `

	tx, err := l.Factory.Db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for id, address := range addressMap {
		if _, err := stmt.Exec(id, address); err != nil {
			fmt.Printf("Failed to insert or update address for ID %d: %v\n", id, err)
			continue
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
