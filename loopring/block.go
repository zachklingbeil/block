package loopring

import (
	"encoding/json"
	"fmt"
)

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

func (l *Loopring) ProcessBlock(number int64) error {
	response, err := l.fetchBlock(number)
	if err != nil {
		return err
	}

	var blockData map[string]any
	if err := json.Unmarshal(response, &blockData); err != nil {
		return fmt.Errorf("failed to parse block %d: %w", number, err)
	}

	blockNumber := int64(blockData["blockId"].(float64))
	blockTime := int64(blockData["createdAt"].(float64))
	transactions, ok := blockData["transactions"].([]any)
	if !ok {
		return fmt.Errorf("invalid transactions format in block %d", number)
	}

	l.Txs = l.flatten(blockNumber, blockTime, transactions)
	l.StoreTransactions(blockNumber, l.Txs)
	return nil
}

func (l *Loopring) fetchBlock(number int64) ([]byte, error) {
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/block/getBlock?id=%d", number)
	response, err := l.Factory.Json.In(url, "")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block %d: %w", number, err)
	}
	return response, nil
}

func (l *Loopring) flatten(blockNumber, blockTime int64, transactions []any) []any {
	var flattened []any
	for i, tx := range transactions {
		if txData, ok := tx.(map[string]any); ok {
			coordinates := l.coordinates(blockNumber, blockTime, int64(i+1))
			flatTx := l.Factory.Json.FlattenMap(txData, "")
			flatTx["coordinates"] = coordinates
			cleanedTx := l.Factory.Json.Cleanup(flatTx)
			flattened = append(flattened, cleanedTx)
		} else {
			fmt.Printf("Unexpected transaction format: %+v\n", tx)
		}
	}
	return flattened
}

// func flattenMap(input map[string]any, prefix string) map[string]any {
// 	flatMap := make(map[string]any)

// 	for key, value := range input {
// 		newKey := key
// 		if prefix != "" {
// 			newKey = prefix + "." + key
// 		}

// 		switch v := value.(type) {
// 		case map[string]any:
// 			maps.Copy(flatMap, flattenMap(v, newKey))
// 		case []any:
// 			for i, item := range v {
// 				arrayKey := fmt.Sprintf("%s[%d]", newKey, i)
// 				if nestedMap, ok := item.(map[string]any); ok {
// 					maps.Copy(flatMap, flattenMap(nestedMap, arrayKey))
// 				} else {
// 					flatMap[arrayKey] = item
// 				}
// 			}
// 		default:
// 			flatMap[newKey] = v
// 		}
// 	}
// 	return flatMap
// }

// func cleanup(data map[string]any) map[string]any {
// 	cleaned := make(map[string]any)
// 	for key, value := range data {
// 		switch v := value.(type) {
// 		case string:
// 			if v != "" {
// 				cleaned[key] = v
// 			}
// 		case []any:
// 			if len(v) > 0 {
// 				cleaned[key] = v
// 			}
// 		case map[string]any:
// 			nested := cleanup(v)
// 			if len(nested) > 0 {
// 				cleaned[key] = nested
// 			}
// 		default:
// 			if v != nil {
// 				cleaned[key] = v
// 			}
// 		}
// 	}
// 	return cleaned
// }
