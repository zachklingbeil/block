package loopring

import (
	"encoding/json"
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Map     map[int64][]*Transaction
}

func NewLoopring(factory *factory.Factory) *Loopring {
	return &Loopring{
		Factory: factory,
		Map:     make(map[int64][]*Transaction),
	}
}

// CurrentBlock fetches the latest block Number from the Loopring API.
func (l *Loopring) CurrentBlock() (int64, error) {
	response, err := l.Factory.Json.In("https://api3.loopring.io/api/v3/block/getBlock", false, "")
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
	response, err := l.Factory.Json.In(url, false, "")
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
	l.Map[int64(number)] = transactions

	l.Factory.Json.Print(block.Size)
	return nil
}
