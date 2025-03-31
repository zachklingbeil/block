package loopring

import (
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
}

func NewLoopring(factory *factory.Factory) *Loopring {
	return &Loopring{
		Factory: factory,
	}
}

// GetBlock fetches block data from the Loopring API.
// If blockNumber is 0, it fetches the latest block.
// Otherwise, it fetches the block with the specified block number.
func (l *Loopring) GetBlock(blockNumber int) (map[string]interface{}, error) {
	var url string
	if blockNumber == 0 {
		url = "https://api3.loopring.io/api/v3/block/getBlock"
	} else {
		url = fmt.Sprintf("https://api3.loopring.io/api/v3/block/getBlock?id=%d", blockNumber)
	}

	// Use factory.Json.In to fetch and decode the JSON response
	var result map[string]interface{}
	err := l.Factory.Json.In(url, &result, false, "")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block: %w", err)
	}

	return result, nil
}
