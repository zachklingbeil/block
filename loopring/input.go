package loopring

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/log"
	"github.com/zachklingbeil/block/circuit"
)

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
