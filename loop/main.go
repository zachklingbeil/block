package loop

import (
	"encoding/json"
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Blocks  []NewBlock
}

func Connect(factory *factory.Factory) *Loopring {
	loop := &Loopring{
		Factory: factory,
	}
	go loop.Listen()
	// go loop.FetchBlocks()
	return loop
}

func (l *Loopring) FetchBlocks() {
	current := l.currentBlock()
	blockHeight := l.blockHeight()

	if blockHeight >= current {
		return
	}

	for i := blockHeight + 1; i <= current; i++ {
		fmt.Printf("%d\n", i)
		if err := l.ProcessBlock(i); err != nil {
			fmt.Printf("Failed to process block %d: %v\n", i, err)
			continue
		}
	}
}

// Helper function to fetch the highest block ID
func (l *Loopring) blockHeight() int64 {
	var blockHeight int64
	err := l.Factory.Db.QueryRow(`SELECT COALESCE(MAX(block), 0) FROM loopring`).Scan(&blockHeight)
	if err != nil {
		fmt.Printf("Failed to fetch the highest block ID: %v\n", err)
		return 0
	}
	return blockHeight
}

// Simplified GetCurrentBlockNumber
func (l *Loopring) currentBlock() int64 {
	data, err := l.Factory.Json.In("https://api3.loopring.io/api/v3/block/getBlock", "")
	if err != nil {
		fmt.Printf("Failed to fetch block data: %v\n", err)
		return 0
	}

	var blockData struct {
		BlockId int64 `json:"blockId"`
	}

	err = json.Unmarshal(data, &blockData)
	if err != nil {
		fmt.Printf("Failed to parse block data: %v\n", err)
		return 0
	}

	return blockData.BlockId
}
