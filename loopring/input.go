package loopring

import (
	"encoding/json"
	"fmt"
	"reflect"

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

type Raw struct {
	Number       int64   `json:"blockId"`
	Timestamp    int64   `json:"createdAt"`
	Size         int64   `json:"blockSize"`
	Coord        fx.Zero `json:"coordinate"`
	Transactions []any   `json:"transactions"`
}

type Block struct {
	Coord        fx.Zero `json:"coordinate"`
	Transactions []Tx    `json:"transactions"`
}

func (l *Loopring) FetchBlock(number int64) (fx.Zero, []any, error) {
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/block/getBlock?id=%d", number)
	response, err := l.Factory.Json.In(url, "")
	if err != nil {
		log.Error("Failed to fetch block data: %v", err)
		return fx.Zero{}, nil, err
	}

	var block Raw
	if err := json.Unmarshal(response, &block); err != nil {
		log.Error("Failed to parse block data: %v", err)
		return fx.Zero{}, nil, err
	}

	coord, txs, err := l.Factory.Circuit.Coordinates(block.Number, block.Timestamp, block.Transactions)
	if err != nil {
		log.Error("Failed to get coordinates: %v", err)
		return fx.Zero{}, nil, err
	}
	return coord, txs, nil
}

var txTypeToStruct = map[string]any{
	"Deposit":       &Deposit{},
	"Withdraw":      &Withdrawal{},
	"SpotTrade":     &Swap{},
	"Transfer":      &Transfer{},
	"NftMint":       &Mint{},
	"AccountUpdate": &AccountUpdate{},
	"AmmUpdate":     &AmmUpdate{},
	"NftData":       &NftData{},
}

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

		structPtr, exists := txTypeToStruct[txType]
		if !exists {
			log.Warn("Unknown transaction type: %s", txType)
			continue
		}
		structInstance := reflect.New(reflect.TypeOf(structPtr).Elem()).Interface()
		if err := mapToStruct(txMap, structInstance); err != nil {
			log.Error("Failed to unmarshal transaction: %v", err)
			continue
		}

		var txObj Tx
		switch txType {
		case "Deposit":
			txObj = l.DepositToTx(*structInstance.(*Deposit))
		case "Withdraw":
			txObj = l.WithdrawToTx(*structInstance.(*Withdrawal))
		case "SpotTrade":
			txObj = l.SwapToTx(*structInstance.(*Swap))
		case "Transfer":
			txObj = l.TransferToTx(*structInstance.(*Transfer))
		case "NftMint":
			txObj = l.MintToTx(*structInstance.(*Mint))
		case "AccountUpdate":
			txObj = l.AccountUpdateToTx(*structInstance.(*AccountUpdate))
		case "AmmUpdate":
			txObj = l.AmmUpdateToTx(*structInstance.(*AmmUpdate))
		case "NftData":
			txObj = l.NftDataToTx(*structInstance.(*NftData))
		default:
			log.Warn("Unhandled transaction type: %s", txType)
			continue
		}
		processedTxs = append(processedTxs, txObj)
	}

	return processedTxs, nil
}

func mapToStruct(data map[string]any, target any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	if err := json.Unmarshal(bytes, target); err != nil {
		return fmt.Errorf("failed to unmarshal data into target struct: %w", err)
	}
	return nil
}
