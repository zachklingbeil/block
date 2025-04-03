package loopring

import (
	"encoding/json"
	"fmt"
)

type Block struct {
	Created      int64         `json:"createdAt"`
	Number       int64         `json:"blockId"`
	Size         int64         `json:"blockSize"`
	TxHash       string        `json:"txHash"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	TxType    TxType `json:"txType"`
	From      int64  `json:"accountId"`
	To        int64  `json:"toAccountId"`
	ToAddress string `json:"toAccountAddress"`
}

type TxType string

const (
	Transfer TxType = "Transfer, Deposit, Withdraw"
)

// FetchBlocks fetches blocks sequentially from the last fetched block to the current block and stores them in the database.
func (l *Loopring) FetchBlocks() error {
	// Fetch the current block number directly
	response, err := l.Factory.Json.In("https://api3.loopring.io/api/v3/block/getBlock", "")
	if err != nil {
		return fmt.Errorf("failed to fetch the latest block data: %w", err)
	}

	var block Block
	if err := json.Unmarshal(response, &block); err != nil {
		return fmt.Errorf("failed to parse block data: %w", err)
	}
	currentBlock := block.Number

	// Get the highest block ID from the database
	query := `SELECT COALESCE(MAX(block_id), 0) FROM blocks`
	var blockHeight int64
	if err := l.Db.QueryRow(query).Scan(&blockHeight); err != nil {
		return fmt.Errorf("failed to fetch the highest block ID: %w", err)
	}

	// Fetch and store each block sequentially
	for i := blockHeight + 1; i <= currentBlock; i++ {
		if err := l.GetBlock(int(i)); err != nil {
			fmt.Printf("Failed to fetch block %d: %v\n", i, err)
			continue
		}
	}
	return nil
}

// GetBlock fetches a block from the Loopring API and inserts it into the database.
func (l *Loopring) GetBlock(number int) error {
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/block/getBlock?id=%d", number)
	response, err := l.Factory.Json.In(url, "")
	if err != nil {
		return fmt.Errorf("failed to fetch block data for block number %d: %w", number, err)
	}

	var block Block
	if err := json.Unmarshal(response, &block); err != nil {
		return fmt.Errorf("failed to parse block data for block number %d: %w", number, err)
	}

	if err := l.InsertBlock(&block); err != nil {
		return fmt.Errorf("failed to insert block into database: %w", err)
	}
	return nil
}

// PeerMap extracts a map of account IDs to their corresponding addresses and identifies IDs without addresses.
func (l *Loopring) PeerMap() (map[int64]string, []int64) {
	peerMap := make(map[int64]string)
	missingAddresses := []int64{}

	// Iterate over all blocks and transactions
	for _, block := range l.Blocks {
		for _, tx := range block.Transactions {
			// Map "From" account ID to its address
			if _, exists := peerMap[tx.From]; !exists {
				peerMap[tx.From] = tx.ToAddress
			}
			// Map "To" account ID to its address
			if _, exists := peerMap[tx.To]; !exists {
				peerMap[tx.To] = tx.ToAddress
			}
		}
	}

	// Identify account IDs without addresses
	for accountID, address := range peerMap {
		if address == "" {
			missingAddresses = append(missingAddresses, accountID)
		}
	}

	// Calculate the number of unique addresses
	uniqueAddresses := make(map[string]struct{})
	for _, address := range peerMap {
		uniqueAddresses[address] = struct{}{}
	}

	// Print the size of the address map and unique counts
	fmt.Printf("Number of unique accounts: %d\n", len(peerMap))
	fmt.Printf("Number of unique addresses: %d\n", len(uniqueAddresses))
	return peerMap, missingAddresses
}

// PeerTable creates the addresses table if it doesn't already exist.
func (l *Loopring) PeerTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS peers (
        address TEXT PRIMARY KEY,       -- Ethereum address
        id BIGINT,                      -- Loopring account ID
        ens TEXT,                       -- [peer].eth
        loopringEns TEXT                -- [peer].loopring.eth
    );`
	if _, err := l.Db.Exec(query); err != nil {
		return fmt.Errorf("failed to create addresses table: %w", err)
	}
	return nil
}

// MapToTable stores the addressMap in the addresses table and resolves conflicts by updating existing rows.
func (l *Loopring) MapToTable(addressMap map[int64]string) error {
	if err := l.PeerTable(); err != nil {
		return fmt.Errorf("failed to create addresses table: %w", err)
	}

	query := `
    INSERT INTO peers (id, address)
    VALUES ($1, $2)
    ON CONFLICT (address) DO UPDATE
    SET id = EXCLUDED.id;
    `

	// Use a transaction for batch inserts
	tx, err := l.Db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Insert each entry in the addressMap
	for id, address := range addressMap {
		if _, err := stmt.Exec(id, address); err != nil {
			return fmt.Errorf("failed to insert or update address in addresses table: %w", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (l *Loopring) GetPeers() error {
	// Step 1: Load all blocks from the database into memory
	if err := l.DiskToMem(); err != nil {
		return fmt.Errorf("failed to load database: %w", err)
	}

	// Step 2: Extract the account-to-address map and missing addresses
	peerMap, missingAddresses := l.PeerMap()

	// Step 3: Store the initial address map in the database
	if err := l.MapToTable(peerMap); err != nil {
		return fmt.Errorf("failed to store address map: %w", err)
	}

	// Step 4: Fetch and update missing addresses
	if len(missingAddresses) > 0 {
		if err := l.FetchMissingAddresses(peerMap, missingAddresses); err != nil {
			return fmt.Errorf("failed to fetch missing addresses: %w", err)
		}

		// Step 5: Store the updated address map in the database
		if err := l.MapToTable(peerMap); err != nil {
			return fmt.Errorf("failed to update address map: %w", err)
		}
	}

	return nil
}

// FetchMissingAddresses fetches addresses for account IDs that do not have addresses in the PeerMap.
func (l *Loopring) FetchMissingAddresses(peerMap map[int64]string, missingAccounts []int64) error {
	totalAccounts := len(missingAccounts)

	for i, accountID := range missingAccounts {
		// Print progress using l.Factory.Json.Print
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
		// Update the PeerMap with the fetched address
		peerMap[accountID] = accountData.Owner
	}
	return nil
}
