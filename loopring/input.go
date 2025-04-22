package loopring

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/zachklingbeil/factory/fx"
)

type Tx struct {
	Zero        any             `json:"zero,omitempty"`
	One         any             `json:"one,omitempty"`
	Value       string          `json:"value,omitempty"`
	Token       any             `json:"token,omitempty"`
	Fee         any             `json:"fee,omitempty"`
	FeeToken    int64           `json:"feeToken,omitempty"`
	OneValue    string          `json:"oneValue,omitempty"`
	OneToken    int64           `json:"oneToken,omitempty"`
	OneFee      any             `json:"oneFee,omitempty"`
	OneFeeToken int64           `json:"oneFeeToken,omitempty"`
	Type        string          `json:"type,omitempty"`
	Index       uint16          `json:"index"`
	Raw         json.RawMessage `json:"raw,omitempty"`
}

func (l *Loopring) FetchBlock(number int64) error {
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/block/getBlock?id=%d", number)
	response, err := l.Factory.Json.In(url, "")
	if err != nil {
		log.Error("Failed to fetch block data: %v", err)
		return nil
	}

	var block Raw
	if err := json.Unmarshal(response, &block); err != nil {
		log.Error("Failed to parse block data: %v", err)
		return nil
	}
	l.Raw = &block
	return nil
}

func (l *Loopring) Coordinates() fx.Zero {
	t := time.UnixMilli(l.Raw.Timestamp)
	return fx.Zero{
		Block:       l.Raw.Number,
		Year:        uint8(t.Year() - 2015),
		Month:       uint8(t.Month()),
		Day:         uint8(t.Day()),
		Hour:        uint8(t.Hour()),
		Minute:      uint8(t.Minute()),
		Second:      uint8(t.Second()),
		Millisecond: uint16(t.Nanosecond() / 1e6),
		Index:       0,
	}
}

func (l *Loopring) ProcessTransactions() error {
	coord := l.Coordinates() // Get the block's coordinates
	var processedTxs []Tx

	for i, tx := range l.Raw.Transactions {
		if txMap, ok := tx.(map[string]any); ok {
			if txType, ok := txMap["txType"].(string); ok {
				var processedTx Tx
				switch txType {
				case "Deposit":
					var deposit Deposit
					if err := mapToStruct(txMap, &deposit); err == nil {
						processedTx = l.DepositToTx(deposit)
					}
				case "Withdraw":
					var withdrawal Withdrawal
					if err := mapToStruct(txMap, &withdrawal); err == nil {
						processedTx = l.WithdrawToTx(withdrawal)
					}
				case "SpotTrade":
					var swap Swap
					if err := mapToStruct(txMap, &swap); err == nil {
						processedTx = l.SwapToTx(swap)
					}
				case "Transfer":
					var transfer Transfer
					if err := mapToStruct(txMap, &transfer); err == nil {
						processedTx = l.TransferToTx(transfer)
					}
				case "NftMint":
					var mint Mint
					if err := mapToStruct(txMap, &mint); err == nil {
						processedTx = l.MintToTx(mint)
					}
				case "AccountUpdate":
					var accountUpdate AccountUpdate
					if err := mapToStruct(txMap, &accountUpdate); err == nil {
						processedTx = l.AccountUpdateToTx(accountUpdate)
					}
				case "NftData":
					var nftData NftData
					if err := mapToStruct(txMap, &nftData); err == nil {
						processedTx = l.NftDataToTx(nftData)
					}
				case "AmmUpdate":
					var ammUpdate AmmUpdate
					if err := mapToStruct(txMap, &ammUpdate); err == nil {
						processedTx = l.AmmUpdateToTx(ammUpdate)
					}
				default:
					fmt.Printf("Unknown transaction type: %s\n", txType)
					continue
				}
				processedTx.Index = uint16(i + 1)
				processedTxs = append(processedTxs, processedTx)
			}
		}
	}
	l.Block = &Block{
		Coord:        coord,
		Transactions: processedTxs,
	}
	return nil
}

func (l *Loopring) ToMap() map[fx.Zero][]Tx {
	if l.Block == nil {
		return nil
	}
	return map[fx.Zero][]Tx{
		l.Block.Coord: l.Block.Transactions,
	}
}

func mapToStruct(data map[string]any, target any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}
