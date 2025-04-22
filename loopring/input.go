package loopring

import (
	"github.com/ethereum/go-ethereum/log"
)

func (l *Loopring) ProcessTransactions(txs []any) ([]Tx, error) {
	var processedTxs []Tx

	for _, tx := range txs {
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

		var txObj Tx
		switch txType {
		case "Deposit":
			txObj = l.DepositToTx(txMap)
		case "Withdraw":
			txObj = l.WithdrawToTx(txMap)
		case "SpotTrade":
			txObj = l.SwapToTx(txMap)
		case "Transfer":
			txObj = l.TransferToTx(txMap)
		case "NftMint":
			txObj = l.MintToTx(txMap)
		case "AccountUpdate":
			txObj = l.AccountUpdateToTx(txMap)
		case "AmmUpdate":
			txObj = l.AmmUpdateToTx(txMap)
		case "NftData":
			txObj = l.NftDataToTx(txMap)
		default:
			log.Warn("Unhandled transaction type: %s", txType)
			continue
		}
		processedTxs = append(processedTxs, txObj)
	}
	return processedTxs, nil
}
