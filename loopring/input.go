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

func (l *Loopring) FetchBlocks() {
	current := l.currentBlock()
	blockHeight := l.blockHeight()
	if blockHeight == current {
		fmt.Println("blockHeight = currentBlock")
		return
	}
	for i := blockHeight + 1; i <= current; i++ {
		if err := l.GetBlock(int(i)); err != nil {
			fmt.Printf("Failed to fetch block %d: %v\n", i, err)
			continue
		}
	}
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

// GetBlock fetches a block and inserts it into the database.
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

	// Process the block's transactions to extract peers
	l.getPeers(&block)

	// Insert the block into the database
	if err := l.InsertBlock(&block); err != nil {
		return fmt.Errorf("failed to insert block into database: %w", err)
	}
	return nil
}

// getPeers extracts unique addresses from a block's transactions.
func (l *Loopring) getPeers(block *Block) {
	one := make(map[int64]string)

	for _, tx := range block.Transactions {
		if tx.To != 0 {
			one[tx.To] = tx.ToAddress
		}
		if tx.From != 0 && one[tx.From] == "" {
			one[tx.From] = ""
		}
	}

	// Convert unique accounts to a slice of addresses
	addresses := make([]string, 0, len(one))
	for _, address := range one {
		if address != "" {
			addresses = append(addresses, address)
		}
	}

	go func() {
		l.Factory.Peer.NewBlock(addresses)
	}()
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

// QualityControl checks if each block in the database has transactions.
// If a block does not have transactions, it fetches the block and updates the database.
func (l *Loopring) QualityControl() error {
	query := `SELECT block_id, transactions FROM loopring`

	rows, err := l.Factory.Db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query blocks from the database: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var blockID int64
		var transactionsJSON []byte

		if err := rows.Scan(&blockID, &transactionsJSON); err != nil {
			return fmt.Errorf("failed to scan block row: %w", err)
		}

		// Check if the transactions slice is empty
		if len(transactionsJSON) == 0 || string(transactionsJSON) == "[]" {
			fmt.Printf("Block %d has no transactions. Fetching block data...\n", blockID)
			if err := l.GetBlock(int(blockID)); err != nil {
				fmt.Printf("Failed to fetch block %d: %v\n", blockID, err)
				continue
			}
			fmt.Printf("Successfully updated block %d with transactions.\n", blockID)
		}
	}
	return nil
}
