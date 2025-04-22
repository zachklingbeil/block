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
	Coordinates fx.Zero         `json:"coordinates"`
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

	var block Block
	if err := json.Unmarshal(response, &block); err != nil {
		log.Error("Failed to parse block data: %v", err)
		return nil
	}
	l.Block = &block
	return nil
}

func (l *Loopring) Index() error {
	for i, transaction := range l.Block.Transactions {
		if txMap, ok := transaction.(map[string]any); ok {
			txMap["coordinate"] = fx.Zero{
				Block:       l.Block.Coord.Block,
				Year:        l.Block.Coord.Year,
				Month:       l.Block.Coord.Month,
				Day:         l.Block.Coord.Day,
				Hour:        l.Block.Coord.Hour,
				Minute:      l.Block.Coord.Minute,
				Second:      l.Block.Coord.Second,
				Millisecond: l.Block.Coord.Millisecond,
				Index:       uint16(i + 1),
			}
			l.Block.Transactions[i] = txMap
		} else {
			return fmt.Errorf("transaction at index %d is not a map", i)
		}
	}
	return nil
}

func (l *Loopring) Coordinates() error {
	t := time.UnixMilli(l.Block.Timestamp)
	l.Block.Coord = fx.Zero{
		Block:       l.Block.Number,
		Year:        uint8(t.Year() - 2015),
		Month:       uint8(t.Month()),
		Day:         uint8(t.Day()),
		Hour:        uint8(t.Hour()),
		Minute:      uint8(t.Minute()),
		Second:      uint8(t.Second()),
		Millisecond: uint16(t.Nanosecond() / 1e6),
		Index:       0,
	}
	return nil
}

func (l *Loopring) ProcessTransactions() {
	var fx []Tx

	for _, tx := range l.Block.Transactions {
		if txMap, ok := tx.(map[string]any); ok {
			if txType, ok := txMap["txType"].(string); ok {
				switch txType {
				case "Deposit":
					var deposit DW
					if err := mapToStruct(txMap, &deposit); err == nil {
						fx = append(fx, l.DepositToTx(deposit))
					}
				case "Withdraw":
					var withdrawal DW
					if err := mapToStruct(txMap, &withdrawal); err == nil {
						fx = append(fx, l.WithdrawToTx(withdrawal))
					}
				case "SpotTrade":
					var swap Swap
					if err := mapToStruct(txMap, &swap); err == nil {
						fx = append(fx, l.SwapToTx(swap))
					}
				case "Transfer":
					var transfer Transfer
					if err := mapToStruct(txMap, &transfer); err == nil {
						fx = append(fx, l.TransferToTx(transfer))
					}
				case "NftMint":
					var mint Mint
					if err := mapToStruct(txMap, &mint); err == nil {
						fx = append(fx, l.MintToTx(mint))
					}
				case "AccountUpdate":
					var accountUpdate AccountUpdate
					if err := mapToStruct(txMap, &accountUpdate); err == nil {
						fx = append(fx, l.AccountUpdateToTx(accountUpdate))
					}
				case "AmmUpdate":
					var ammUpdate AmmUpdate
					if err := mapToStruct(txMap, &ammUpdate); err == nil {
						fx = append(fx, l.AmmUpdateToTx(ammUpdate))
					}
				case "NftData":
					var nftData NftData
					if err := mapToStruct(txMap, &nftData); err == nil {
						fx = append(fx, l.NftDataToTx(nftData))
					}
				default:
					fmt.Printf("Unknown transaction type: %s\n", txType)
				}
			}
		}
	}
	l.Block.Transactions = make([]any, len(fx))
	for i, tx := range fx {
		l.Block.Transactions[i] = tx
	}
}

func mapToStruct(data map[string]any, target any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}
