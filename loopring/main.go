package loopring

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Map     map[int64][]*Transaction
	Db      *sql.DB
}

func NewLoopring(factory *factory.Factory) (*Loopring, error) {
	db, err := factory.Db.Connect("loopring")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Loopring database: %w", err)
	}

	return &Loopring{
		Factory: factory,
		Map:     make(map[int64][]*Transaction),
		Db:      db,
	}, nil
}

// Helper function to update the map with transactions for a given block number.
func (l *Loopring) Write(blockNumber int64, transactions []*Transaction) {
	l.Factory.Mu.Lock()
	defer l.Factory.Mu.Unlock()
	l.Map[blockNumber] = transactions
}

// Helper function to read transactions from the map for a given block number.
func (l *Loopring) Read(blockNumber int64) ([]*Transaction, bool) {
	l.Factory.Mu.Lock()
	defer l.Factory.Mu.Unlock()
	transactions, exists := l.Map[blockNumber]
	return transactions, exists
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
	l.Factory.Json.Print(block.Number)
	return block.Number, nil
}

// GetBlock fetches block data from the Loopring API and updates the map with transactions.
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

	transactions := make([]*Transaction, len(block.Transactions))
	for i, tx := range block.Transactions {
		transactions[i] = &tx
	}

	// Use the helper function to update the map
	l.Write(int64(number), transactions)

	// Print the block size
	l.Factory.Json.Print(block.Size)
	return nil
}
