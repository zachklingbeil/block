package loopring

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/log"
	"github.com/zachklingbeil/factory/fx"
)

type Raw struct {
	Number       int64   `json:"blockId"`
	Timestamp    int64   `json:"createdAt"`
	Size         int64   `json:"blockSize"`
	Coord        fx.Zero `json:"coordinate"`
	Transactions []any   `json:"transactions"`
}

type Block struct {
	Coord        fx.Zero `json:"coordinate"`
	Transactions []Tx    `json:"transactions"`
}

func (l *Loopring) FetchBlock(number int64) (fx.Zero, []any, error) {
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/block/getBlock?id=%d", number)
	response, err := l.Factory.Json.In(url, "")
	if err != nil {
		log.Error("Failed to fetch block data: %v", err)
		return fx.Zero{}, nil, err
	}

	var block Raw
	if err := json.Unmarshal(response, &block); err != nil {
		log.Error("Failed to parse block data: %v", err)
		return fx.Zero{}, nil, err
	}

	simpleTxs := l.Factory.Json.Simplify(block.Transactions, "")
	coord, txs, err := l.Factory.Circuit.Coordinates(block.Number, block.Timestamp, simpleTxs)
	if err != nil {
		log.Error("Failed to get coordinates: %v", err)
		return fx.Zero{}, nil, err
	}

	return coord, txs, nil
}
