package loopring

import (
	"encoding/json"
	"fmt"
	"time"

	"maps"

	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Txs     []any
}

func NewLoopring(factory *factory.Factory) *Loopring {
	return &Loopring{Factory: factory}
}

func (l *Loopring) ProcessBlock(number int) error {
	response, err := l.fetchBlock(number)
	if err != nil {
		return err
	}

	var blockData map[string]any
	if err := json.Unmarshal(response, &blockData); err != nil {
		return fmt.Errorf("failed to parse block %d: %w", number, err)
	}

	// Extract transactions and metadata
	blockNumber := int64(blockData["blockId"].(float64))
	blockTime := int64(blockData["createdAt"].(float64))
	transactions, ok := blockData["transactions"].([]any)
	if !ok {
		return fmt.Errorf("invalid transactions format in block %d", number)
	}

	// Process and store transactions
	l.Txs = l.flattenTransactions(blockNumber, blockTime, transactions)
	l.StoreTransactions(blockNumber, l.Txs)

	return nil
}

func (l *Loopring) flattenTransactions(blockNumber, blockTime int64, transactions []any) []any {
	var flattened []any

	for i, tx := range transactions {
		if txData, ok := tx.(map[string]any); ok {
			coordinates := generateCoordinates(blockTime, int64(i+1))
			flatTx := flattenMap(txData, "")
			flatTx["block"] = blockNumber
			flatTx["index"] = int64(i + 1)
			flatTx["coordinates"] = coordinates
			flattened = append(flattened, flatTx)
		} else {
			fmt.Printf("Unexpected transaction format: %+v\n", tx)
		}
	}
	return flattened
}

func generateCoordinates(timestamp int64, index int64) string {
	t := time.UnixMilli(timestamp)

	year := int64(t.Year() - 2015)
	month := int64(t.Month())
	day := int64(t.Day())
	hour := int64(t.Hour())
	minute := int64(t.Minute())
	second := int64(t.Second())
	millisecond := int64(t.Nanosecond() / 1e6)

	// Format the coordinates string
	return fmt.Sprintf("%d.%d.%d.%d.%d.%d.%d.%d", year, month, day, hour, minute, second, millisecond, index)
}

func (l *Loopring) fetchBlock(number int) ([]byte, error) {
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/block/getBlock?id=%d", number)
	response, err := l.Factory.Json.In(url, "")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block %d: %w", number, err)
	}
	return response, nil
}

func flattenMap(input map[string]any, prefix string) map[string]any {
	flatMap := make(map[string]any)

	for key, value := range input {
		newKey := key
		if prefix != "" {
			newKey = prefix + "." + key
		}

		switch v := value.(type) {
		case map[string]any:
			maps.Copy(flatMap, flattenMap(v, newKey))
		case []any:
			for i, item := range v {
				arrayKey := fmt.Sprintf("%s[%d]", newKey, i)
				if nestedMap, ok := item.(map[string]any); ok {
					maps.Copy(flatMap, flattenMap(nestedMap, arrayKey))
				} else {
					flatMap[arrayKey] = item
				}
			}
		default:
			flatMap[newKey] = v
		}
	}
	return flatMap
}
func (l *Loopring) StoreTransactions(blockNumber int64, transactions []any) error {
	// Ensure the table exists
	createTableQuery := `
        CREATE TABLE IF NOT EXISTS loopring (
            block BIGINT PRIMARY KEY,
            transactions JSONB NOT NULL
        );
    `
	_, err := l.Factory.Db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	// Convert transactions to JSON
	txJSON, err := json.Marshal(transactions)
	if err != nil {
		return fmt.Errorf("failed to marshal transactions: %w", err)
	}

	// Insert into the database
	query := `
        INSERT INTO loopring (block, transactions)
        VALUES ($1, $2)
        ON CONFLICT (block) DO UPDATE
        SET transactions = EXCLUDED.transactions;
    `
	_, err = l.Factory.Db.Exec(query, blockNumber, txJSON)
	if err != nil {
		return fmt.Errorf("failed to store transactions: %w", err)
	}

	return nil
}

// func (l *Loopring) StoreTransactions(blockNumber int64, transactions []any) error {
// 	txJSON, err := json.Marshal(transactions)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal transactions: %w", err)
// 	}

// 	// Insert into the database
// 	query := `
//         INSERT INTO loopring (block, transactions)
//         VALUES ($1, $2)
//         ON CONFLICT (block) DO UPDATE
//         SET transactions = EXCLUDED.transactions;
//     `
// 	_, err = l.Factory.Db.Exec(query, blockNumber, txJSON)
// 	if err != nil {
// 		return fmt.Errorf("failed to store transactions: %w", err)
// 	}

// 	return nil
// }
