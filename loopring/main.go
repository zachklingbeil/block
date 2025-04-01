package loopring

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Db      *sql.DB
}

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

// NewLoopring initializes a new Loopring instance and ensures the database table exists.
func NewLoopring(factory *factory.Factory) (*Loopring, error) {
	db, err := factory.Db.Connect("loopring")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Loopring database: %w", err)
	}

	loopring := &Loopring{
		Factory: factory,
		Db:      db,
	}

	if err := loopring.CreateTable(); err != nil {
		return nil, fmt.Errorf("failed to create blocks table: %w", err)
	}
	return loopring, nil
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

	// Insert the block into the database
	if err := l.InsertBlock(&block); err != nil {
		return fmt.Errorf("failed to insert block into database: %w", err)
	}

	return nil
}

// CurrentBlock fetches the latest block Number from the Loopring API.
func (l *Loopring) CurrentBlock() (int64, error) {
	response, err := l.Factory.Json.In("https://api3.loopring.io/api/v3/block/getBlock", "")
	if err != nil {
		return 0, fmt.Errorf("failed to fetch the latest block data: %w", err)
	}

	var block Block
	if err := json.Unmarshal(response, &block); err != nil {
		return 0, fmt.Errorf("failed to parse block data: %w", err)
	}
	return block.Number, nil
}

// FetchBlocks fetches blocks sequentially from the last fetched block to the current block and stores them in the database.
func (l *Loopring) FetchBlocks() error {
	// Get the current block number
	currentBlock, err := l.CurrentBlock()
	if err != nil {
		return fmt.Errorf("failed to fetch the current block number: %w", err)
	}

	// Get the highest block ID from the database
	blockHeight, err := l.ContinueFetch()
	if err != nil {
		return fmt.Errorf("failed to fetch the highest block ID: %w", err)
	}

	// Fetch and store each block sequentially
	for i := blockHeight + 1; i <= currentBlock; i++ {
		if err := l.GetBlock(int(i)); err != nil {
			fmt.Printf("Failed to fetch block %d: %v\n", i, err)
			continue
		}
	}

	fmt.Println("Finished fetching all blocks.")
	return nil
}

// GetHighestBlockID retrieves the highest block_id from the database.
func (l *Loopring) ContinueFetch() (int64, error) {
	query := `SELECT COALESCE(MAX(block_id), 0) FROM blocks`

	var blockHeight int64
	if err := l.Db.QueryRow(query).Scan(&blockHeight); err != nil {
		return 0, fmt.Errorf("failed to fetch the highest block_id from the database: %w", err)
	}
	return blockHeight, nil
}
