package loopring

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/redis/go-redis/v9"
)

func (l *Loopring) BlockByBlock(blockNumber int64) error {
	input := l.FetchBlock(blockNumber)
	transactions, block := l.Coordinates(input)
	txs := l.ProcessBlock(transactions)
	block.Ones = txs

	if err := l.StoreBlock(blockNumber, block); err != nil {
		return fmt.Errorf("failed to store block: %w", err)
	}
	return nil
}

func (l *Loopring) CurrentBlock() int64 {
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

func (l *Loopring) FetchBlock(number int64) *Raw {
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/block/getBlock?id=%d", number)
	response, err := l.Factory.Json.In(url, "")
	if err != nil {
		log.Error("Failed to fetch block data: %v", err)
		return nil
	}
	var input *Raw
	if err := json.Unmarshal(response, &input); err != nil {
		log.Error("Failed to parse block data: %v", err)
		return nil
	}
	return input
}

func (l *Loopring) Coordinates(loop *Raw) ([]any, *Block) {
	for i := range loop.Transactions {
		if tx, ok := loop.Transactions[i].(map[string]any); ok {
			tx["index"] = i + 1
		}
	}
	transactions := l.Factory.Json.Simplify(loop.Transactions, "")
	depth := uint16(len(transactions))

	t := time.UnixMilli(loop.Timestamp)
	coordinate := Coordinate{
		Year:        uint8(t.Year() - 2015),
		Month:       uint8(t.Month()),
		Day:         uint8(t.Day()),
		Hour:        uint8(t.Hour()),
		Minute:      uint8(t.Minute()),
		Second:      uint8(t.Second()),
		Millisecond: uint16(t.Nanosecond() / 1e6),
		Index:       0,
		Depth:       depth,
	}

	block := &Block{
		Number: loop.Number,
		Zero:   coordinate,
		Ones:   make([]Tx, depth),
	}
	return transactions, block
}

func (l *Loopring) ProcessBlock(transactions []any) []Tx {
	var txs []Tx

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

func (l *Loopring) StoreBlock(blockNumber int64, block any) error {
	blockJSON, err := json.Marshal(block)
	if err != nil {
		return fmt.Errorf("failed to marshal block: %w", err)
	}
	score := float64(blockNumber)
	err = l.Factory.Data.RB.ZAdd(l.Factory.Ctx, "blocks", redis.Z{
		Score:  score,
		Member: blockJSON,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to store block in Redis: %w", err)
	}
	return nil
}
