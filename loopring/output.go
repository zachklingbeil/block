package loopring

import (
	"encoding/json"
	"fmt"
)

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
