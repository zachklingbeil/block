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

// FetchBlocks fetches blocks sequentially from the last fetched block to the current block and stores them in the database.
func (l *Loopring) FetchBlocks() error {
	currentBlock := l.GetCurrentBlockNumber()
	blockHeight := l.GetHighestBlockID()

	if blockHeight == currentBlock {
		fmt.Println("blockHeight = currentBlock")
		return nil
	}

	if err := l.FetchAndStoreBlocks(blockHeight, currentBlock); err != nil {
		return err
	}

	if err := l.QualityControl(); err != nil {
		return err
	}

	if err := l.ExtractPeerInfo(); err != nil {
		return err
	}

	return nil
}

// GetCurrentBlockNumber fetches the current block number from the Loopring API.
func (l *Loopring) GetCurrentBlockNumber() int64 {
	response, err := l.Factory.Json.In("https://api3.loopring.io/api/v3/block/getBlock", "")
	if err != nil {
		fmt.Printf("Failed to fetch the latest block data: %v\n", err)
		return 0
	}

	var block Block
	if err := json.Unmarshal(response, &block); err != nil {
		fmt.Printf("Failed to parse block data: %v\n", err)
		return 0
	}

	return block.Number
}

// GetHighestBlockID retrieves the highest block ID stored in the database.
func (l *Loopring) GetHighestBlockID() int64 {
	query := `SELECT COALESCE(MAX(block_id), 0) FROM loopring`
	var blockHeight int64
	if err := l.Factory.Db.QueryRow(query).Scan(&blockHeight); err != nil {
		fmt.Printf("Failed to fetch the highest block ID: %v\n", err)
		return 0
	}
	return blockHeight
}

// FetchAndStoreBlocks fetches and stores blocks sequentially from the last fetched block to the current block.
func (l *Loopring) FetchAndStoreBlocks(startBlock, endBlock int64) error {
	for i := startBlock + 1; i <= endBlock; i++ {
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
