package loopring

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/log"
	"github.com/zachklingbeil/block/universe"
)

func (l *Loopring) Loop() error {
	latestBlock := l.CurrentBlock()
	blocks, _ := l.Factory.State.Get("loop.block")
	startBlock := int64(1)

	if blocks != nil {
		if hb, ok := blocks.(float64); ok {
			startBlock = int64(hb) + 1
		} else {
			log.Error("Invalid type for blocks: %T", blocks)
			return fmt.Errorf("invalid type for blocks")
		}
	}
	for i := startBlock; i <= latestBlock; i++ {
		fmt.Println(i)
		if err := l.BlockByBlock(i); err != nil {
			log.Error("Failed to process block %d: %v", i, err)
			break
		}
		l.Factory.State.Count("loop.block", i)
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

func (l *Loopring) BlockByBlock(blockNumber int64) error {
	input := l.FetchBlock(blockNumber)
	if input == nil {
		return fmt.Errorf("failed to fetch block %d: got nil", blockNumber)
	}
	transactions, coordinate := l.Zero.Coordinates(input)
	txs := l.ProcessBlock(transactions)

	block := &universe.Block{
		Zero: *coordinate,
		Ones: txs,
	}

	if err := l.StoreBlock(block); err != nil {
		return fmt.Errorf("failed to store block: %w", err)
	}
	return nil
}

func (l *Loopring) FetchBlock(number int64) *universe.Raw {
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/block/getBlock?id=%d", number)
	response, err := l.Factory.Json.In(url, "")
	if err != nil {
		log.Error("Failed to fetch block data: %v", err)
		return nil
	}
	var input *universe.Raw
	if err := json.Unmarshal(response, &input); err != nil {
		log.Error("Failed to parse block data: %v", err)
		return nil
	}
	return input
}

func (l *Loopring) ProcessBlock(transactions []any) []universe.Tx {
	var txs []universe.Tx

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
			txs = append(txs, l.SwapToTx(txMap))
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

func (l *Loopring) StoreBlock(block *universe.Block) error {
	blockJSON, _ := json.Marshal(block)
	if err := l.Factory.Data.RB.SAdd(l.Factory.Ctx, "loopring", blockJSON).Err(); err != nil {
		return fmt.Errorf("failed to store block in Redis hash: %w", err)
	}
	return nil
}
