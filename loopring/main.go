package loopring

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/log"
	"github.com/zachklingbeil/block/circuit"
	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Circuit *circuit.Circuit
}

func Connect(factory *factory.Factory, circuit *circuit.Circuit) *Loopring {
	loop := &Loopring{
		Factory: factory,
		Circuit: circuit,
	}
	return loop
}

func (l *Loopring) BlockByBlock(blockNumber int64) error {
	input, err := l.FetchBlock(blockNumber)
	if err != nil {
		log.Error("Failed to fetch block %d: %v", blockNumber, err)
		return err
	}
	transactions, b, err := l.Circuit.Coordinates(input)
	if err != nil {
		log.Error("Failed to get coordinates: %v", err)
		return err
	}

	txs, err := l.ProcessBlock(transactions)
	if err != nil {
		log.Error("Failed to process transactions for block %d: %v", blockNumber, err)
		return err
	}

	block := &circuit.Block{
		Number:      b.Number,
		Year:        b.Year,
		Month:       b.Month,
		Day:         b.Day,
		Hour:        b.Hour,
		Minute:      b.Minute,
		Second:      b.Second,
		Millisecond: b.Millisecond,
		Index:       b.Index,
		Count:       uint16(len(transactions)),
		Txs:         txs,
	}

	blockJSON, err := json.Marshal(block)
	if err != nil {
		log.Error("Failed to serialize block %d: %v", blockNumber, err)
		return err
	}

	err = l.Factory.Redis.SAdd(l.Factory.Ctx, "blocks", blockJSON).Err()
	if err != nil {
		log.Error("Failed to store block %d in Redis: %v", blockNumber, err)
		return err
	}
	fmt.Printf("%d\n", blockNumber)
	return nil
}

func (l *Loopring) FetchBlock(number int64) (*circuit.Raw, error) {
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/block/getBlock?id=%d", number)
	response, err := l.Factory.Json.In(url, "")
	if err != nil {
		log.Error("Failed to fetch block data: %v", err)
		return nil, err
	}
	var input *circuit.Raw
	if err := json.Unmarshal(response, &input); err != nil {
		log.Error("Failed to parse block data: %v", err)
		return nil, err
	}
	return input, nil
}

func (l *Loopring) ProcessBlock(transactions []any) ([]circuit.Tx, error) {
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
			log.Warn("Unhandled transaction type: %s", txType)
			continue
		}
	}
	return txs, nil
}
