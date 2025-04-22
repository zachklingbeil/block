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
			var Deposit Deposit
			txObj = l.DepositToTx(Deposit)
		case "Withdraw":
			var Withdrawal Withdrawal
			txObj = l.WithdrawToTx(Withdrawal)
		case "SpotTrade":
			var SpotTrade SpotTrade
			txObj = l.SwapToTx(SpotTrade)
		case "Transfer":
			var Transfer Transfer
			txObj = l.TransferToTx(Transfer)
		case "NftMint":
			var Mint Mint
			txObj = l.MintToTx(Mint)
		case "AccountUpdate":
			var AccountUpdate AccountUpdate
			txObj = l.AccountUpdateToTx(AccountUpdate)
		case "AmmUpdate":
			var AmmUpdate AmmUpdate
			txObj = l.AmmUpdateToTx(AmmUpdate)
		case "NftData":
			var NftData NftData
			txObj = l.NftDataToTx(NftData)
		default:
			log.Warn("Unhandled transaction type: %s", txType)
			continue
		}
		// Append the processed transaction
		processedTxs = append(processedTxs, txObj)
	}

	l.Factory.Json.Print(processedTxs)
	return processedTxs, nil
}
