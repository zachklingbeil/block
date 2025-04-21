package loopring

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

type Type struct {
	Deposits       []DW            `json:"Deposit,omitempty"`
	Withdrawals    []DW            `json:"Withdraw,omitempty"`
	Swaps          []Swap          `json:"SpotTrade,omitempty"`
	Transfers      []Transfer      `json:"Transfer,omitempty"`
	Mints          []Mint          `json:"NftMint,omitempty"`
	AccountUpdates []AccountUpdate `json:"AccountUpdate,omitempty"`
	AmmUpdates     []AmmUpdate     `json:"AmmUpdate,omitempty"`
	NftData        []NftData       `json:"NftData,omitempty"`
	TBD            []any           `json:"tbd,omitempty"`
	*json.RawMessage
}

type Input struct {
	Number       int64 `json:"blockId"`
	Timestamp    int64 `json:"createdAt"`
	Size         int64 `json:"blockSize"`
	Transactions []any `json:"transactions"`
}

func (l *Loopring) FetchBlock(number int64) *Input {
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/block/getBlock?id=%d", number)
	response, err := l.Factory.Json.In(url, "")
	if err != nil {
		log.Error("Failed to fetch block data: %v", err)
		return nil
	}

	var block Input
	if err := json.Unmarshal(response, &block); err != nil {
		log.Error("Failed to parse block data: %v", err)
		return nil
	}

	l.Block.In = &block
	return &block
}

func (l *Loopring) Coordinates(input *Input) Coordinate {
	for i, tx := range input.Transactions {
		if txMap, ok := tx.(map[string]any); ok {
			txMap["index"] = int64(i + 1)
			input.Transactions[i] = txMap
		}
	}
	t := time.UnixMilli(input.Timestamp)
	coordinate := Coordinate{
		Block:       input.Number,
		Year:        int64(t.Year() - 2015),
		Month:       int64(t.Month()),
		Day:         int64(t.Day()),
		Hour:        int64(t.Hour()),
		Minute:      int64(t.Minute()),
		Second:      int64(t.Second()),
		Millisecond: int64(t.Nanosecond() / 1e6),
	}
	return coordinate
}

func (l *Loopring) ProcessTransactions(input *Input) []Tx {
	var transactions []Tx

	for _, tx := range input.Transactions {
		if txMap, ok := tx.(map[string]any); ok {
			if txType, ok := txMap["txType"].(string); ok {
				switch txType {
				case "Deposit":
					var deposit DW
					if err := mapToStruct(txMap, &deposit); err == nil {
						transactions = append(transactions, l.DepositToTx(deposit))
					}
				case "Withdraw":
					var withdrawal DW
					if err := mapToStruct(txMap, &withdrawal); err == nil {
						transactions = append(transactions, l.WithdrawToTx(withdrawal))
					}
				case "SpotTrade":
					var swap Swap
					if err := mapToStruct(txMap, &swap); err == nil {
						transactions = append(transactions, l.SwapToTx(swap))
					}
				case "Transfer":
					var transfer Transfer
					if err := mapToStruct(txMap, &transfer); err == nil {
						transactions = append(transactions, l.TransferToTx(transfer))
					}
				case "NftMint":
					var mint Mint
					if err := mapToStruct(txMap, &mint); err == nil {
						transactions = append(transactions, l.MintToTx(mint))
					}
				case "AccountUpdate":
					var accountUpdate AccountUpdate
					if err := mapToStruct(txMap, &accountUpdate); err == nil {
						transactions = append(transactions, l.AccountUpdateToTx(accountUpdate))
					}
				case "AmmUpdate":
					var ammUpdate AmmUpdate
					if err := mapToStruct(txMap, &ammUpdate); err == nil {
						transactions = append(transactions, l.AmmUpdateToTx(ammUpdate))
					}
				case "NftData":
					var nftData NftData
					if err := mapToStruct(txMap, &nftData); err == nil {
						transactions = append(transactions, l.NftDataToTx(nftData))
					}
				default:
					fmt.Printf("Unknown transaction type: %s\n", txType)
				}
			}
		}
	}
	return transactions
}

func mapToStruct(data map[string]any, target any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}
