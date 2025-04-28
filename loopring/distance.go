package loopring

import (
	"encoding/json"
	"fmt"
)

func (l *Loopring) Distance() (int64, int64) {
	current := l.currentBlock()
	past, _ := l.getHistory()

	distance := current - past
	if distance > 0 {
		return past, distance
	}
	return past, 0
}

func (l *Loopring) currentBlock() int64 {
	data, err := l.Factory.Json.In("https://api3.loopring.io/api/v3/block/getBlock", "")
	if err != nil {
		fmt.Printf("Failed to fetch block data: %v\n", err)
		return 0
	}
	var block struct {
		Number int64 `json:"blockId"`
	}
	err = json.Unmarshal(data, &block)
	if err != nil {
		fmt.Printf("Failed to parse block data: %v\n", err)
		return 0
	}
	l.CurrentBlock = block.Number
	return block.Number
}
