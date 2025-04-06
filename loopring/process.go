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
	TxType    string `json:"txType"`
	From      int64  `json:"accountId"`
	To        int64  `json:"toAccountId"`
	ToAddress string `json:"toAccountAddress"`
}

// Helper function to fetch the highest block ID
func (l *Loopring) blockHeight() int64 {
	var blockHeight int64
	err := l.Factory.Db.QueryRow(`SELECT COALESCE(MAX(block_id), 0) FROM loopring`).Scan(&blockHeight)
	if err != nil {
		fmt.Printf("Failed to fetch the highest block ID: %v\n", err)
		return 0
	}
	return blockHeight
}

// Simplified GetCurrentBlockNumber
func (l *Loopring) currentBlock() int64 {
	var block Block
	data, err := l.Factory.Json.In("https://api3.loopring.io/api/v3/block/getBlock", "")
	if err != nil {
		fmt.Printf("Failed to fetch block data: %v\n", err)
		return 0
	}
	err = json.Unmarshal(data, &block)
	if err != nil {
		fmt.Printf("Failed to parse block data: %v\n", err)
		return 0
	}
	return block.Number
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

// InsertBlock inserts a block into the database.
func (l *Loopring) InsertBlock(block *Block) error {
	query := `
        INSERT INTO loopring (block_id, block_size, created, tx_hash, transactions)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (created) DO NOTHING
    `

	transactions, err := json.Marshal(block.Transactions)
	if err != nil {
		return fmt.Errorf("failed to marshal transactions: %w", err)
	}

	if _, err := l.Factory.Db.Exec(query, block.Number, block.Size, block.Created, block.TxHash, transactions); err != nil {
		return fmt.Errorf("failed to insert block into database: %w", err)
	}
	l.Factory.Json.Print(block.Number)
	return nil
}

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
	missingCount := l.countMissingAddresses(uniqueAccounts)
	fmt.Printf("Need %d addresses.\n", missingCount)

	// Fetch missing addresses
	if err := l.fetchMissingAddresses(uniqueAccounts, &missingCount); err != nil {
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
func (l *Loopring) countMissingAddresses(uniqueAccounts map[int64]string) int {
	missingCount := 0
	for _, address := range uniqueAccounts {
		if address == "" {
			missingCount++
		}
	}
	return missingCount
}

// Helper function to fetch missing addresses
func (l *Loopring) fetchMissingAddresses(uniqueAccounts map[int64]string, missingCount *int) error {
	for accountId, address := range uniqueAccounts {
		if address == "" {
			idStr := fmt.Sprintf("%d", accountId)
			fmt.Printf("%d\n", *missingCount)
			fetchedAddress := l.Factory.Peer.NeedsAddress(idStr)
			if fetchedAddress != "" {
				uniqueAccounts[accountId] = fetchedAddress
				(*missingCount)--
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

	// Call HelloUniverse for each unique address
	fmt.Printf("Saying hellouniverse to %d peers...\n", len(addresses))
	for i, address := range addresses {
		l.Factory.Peer.HelloUniverse(address)
		fmt.Printf("%d\n", i+1)
	}
	fmt.Println("HelloUniverse")
}
