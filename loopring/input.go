package loopring

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/log"
	"github.com/zachklingbeil/block/circuit"
)

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
	return block.Number
}

// getHistory retrieves the highest block number from the Redis set
func (l *Loopring) getHistory() (int64, error) {
	blockJSONs, err := l.Factory.Redis.SMembers(l.Factory.Ctx, "blocks").Result()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve blocks from Redis: %w", err)
	}
	past := int64(0)
	for _, blockJSON := range blockJSONs {
		var block circuit.Block
		if err := json.Unmarshal([]byte(blockJSON), &block); err != nil {
			log.Error("Failed to deserialize block JSON: %v", err)
			continue
		}
		if block.Number > past {
			past = block.Number
		}
	}
	return past, nil
}

func (l *Loopring) FetchBlock(number int64) *circuit.Raw {
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/block/getBlock?id=%d", number)
	response, err := l.Factory.Json.In(url, "")
	if err != nil {
		log.Error("Failed to fetch block data: %v", err)
		return nil
	}
	var input *circuit.Raw
	if err := json.Unmarshal(response, &input); err != nil {
		log.Error("Failed to parse block data: %v", err)
		return nil
	}
	return input
}

func (l *Loopring) ProcessBlock(transactions []any) []circuit.Tx {
	var txs []circuit.Tx

	for _, tx := range transactions {
		txMap, ok := tx.(map[string]any)
		if !ok {
			log.Error("Invalid transaction format: %v", tx)
			continue
		}

		txType, ok := txMap["txType"].(string)
		if !ok {
			log.Error("Transaction missing txType field: %v", tx)
			continue
		}

		switch txType {
		case "Deposit":
			txs = append(txs, l.DepositToTx(txMap))
		case "Withdraw":
			txs = append(txs, l.WithdrawToTx(txMap))
		case "SpotTrade":
			spotTxs := l.SwapToTx(txMap)
			txs = append(txs, spotTxs...)
		case "Transfer":
			txs = append(txs, l.TransferToTx(txMap))
		case "NftMint":
			txs = append(txs, l.MintToTx(txMap))
		case "AccountUpdate":
			txs = append(txs, l.AccountUpdateToTx(txMap))
		case "AmmUpdate":
			txs = append(txs, l.AmmUpdateToTx(txMap))
		case "NftData":
			txs = append(txs, l.NftDataToTx(txMap))
		default:
			log.Warn("Unhandled type: %s", txType)
			continue
		}
	}
	return txs
}
