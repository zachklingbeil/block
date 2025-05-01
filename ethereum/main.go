package ethereum

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zachklingbeil/factory"
)

type Ethereum struct {
	Factory *factory.Factory
	// Map     map[Coordinate]*Tx
}

func NewEthereum(factory *factory.Factory) *Ethereum {
	return &Ethereum{
		Factory: factory,
	}
}

// Block holds information about a block.
type Block struct {
	Number       uint64
	Hash         string
	ParentHash   string
	Time         uint64
	GasUsed      uint64
	GasLimit     uint64
	BaseFee      *big.Int
	Transactions []*Transactions
}

// Transactions holds information about a transaction.
type Transactions struct {
	Hash              string
	From              string
	To                string
	Value             *big.Int
	Gas               uint64
	GasPrice          *big.Int
	Nonce             uint64
	DataLength        int
	Type              uint8
	Status            uint64
	CumulativeGasUsed uint64
	Logs              []*LogInfo
}

// LogInfo holds information about a transaction log.
type LogInfo struct {
	Address    string
	Topics     []string
	DataLength int
}

// ProcessBlocks processes the latest `count` blocks.
func (e *Ethereum) ProcessBlocks(count int) ([]*Block, error) {
	header, err := e.Factory.Eth.HeaderByNumber(e.Factory.Ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest header: %w", err)
	}
	latestBlock := header.Number.Uint64()

	var blocks []*Block
	for blockNum := latestBlock; blockNum > latestBlock-uint64(count); blockNum-- {
		block, err := e.Factory.Eth.BlockByNumber(e.Factory.Ctx, big.NewInt(int64(blockNum)))
		if err != nil {
			log.Printf("Error fetching block %d: %v", blockNum, err)
			continue
		}
		blockInfo := e.processBlock(e.Factory.Ctx, block)
		blocks = append(blocks, blockInfo)
		err = e.StoreBlock(int64(blockInfo.Number), blockInfo)
		if err != nil {
			log.Printf("Error storing block %d: %v", blockInfo.Number, err)
			continue
		}
		log.Printf("%d", blockInfo.Number)
	}
	return blocks, nil
}

func (e *Ethereum) StoreBlock(blockNumber int64, block any) error {
	// Serialize the block to JSON
	blockJSON, err := json.Marshal(block)
	if err != nil {
		return fmt.Errorf("failed to marshal block: %w", err)
	}

	// Define the Redis hash key for storing blocks
	hashKey := "ethereum"

	// Use the blockNumber as the field in the Redis hash
	err = e.Factory.Data.RB.HSet(e.Factory.Ctx, hashKey, fmt.Sprintf("%d", blockNumber), blockJSON).Err()
	if err != nil {
		return fmt.Errorf("failed to store block in Redis hash: %w", err)
	}
	return nil
}

// processBlock processes a single block and returns its information.
func (e *Ethereum) processBlock(ctx context.Context, block *types.Block) *Block {
	blockInfo := &Block{
		Number:     block.NumberU64(),
		Hash:       block.Hash().Hex(),
		ParentHash: block.ParentHash().Hex(),
		Time:       block.Time(),
		GasUsed:    block.GasUsed(),
		GasLimit:   block.GasLimit(),
		BaseFee:    block.BaseFee(),
	}

	for _, tx := range block.Transactions() {
		txInfo := e.processTransaction(ctx, tx)
		blockInfo.Transactions = append(blockInfo.Transactions, txInfo)
	}
	return blockInfo
}

// processTransaction processes a single transaction and returns its information.
func (e *Ethereum) processTransaction(ctx context.Context, tx *types.Transaction) *Transactions {
	txInfo := &Transactions{
		Hash:       tx.Hash().Hex(),
		Value:      tx.Value(),
		Gas:        tx.Gas(),
		GasPrice:   tx.GasPrice(),
		Nonce:      tx.Nonce(),
		DataLength: len(tx.Data()),
		Type:       tx.Type(),
	}

	// Get sender address
	var from string
	var err error
	if tx.ChainId() == nil || tx.ChainId().Sign() == 0 {
		addr, err := types.Sender(types.HomesteadSigner{}, tx)
		if err == nil {
			from = addr.Hex()
		}
	} else {
		signer := types.LatestSignerForChainID(tx.ChainId())
		addr, err2 := types.Sender(signer, tx)
		if err2 == nil {
			from = addr.Hex()
		}
		err = err2
	}
	if err == nil {
		txInfo.From = from
	}

	// To address (contract creation if nil)
	if tx.To() == nil {
		txInfo.To = "Contract Creation"
	} else {
		txInfo.To = tx.To().Hex()
	}

	// Get receipt for transaction status, logs, etc.
	receipt, err := e.Factory.Eth.TransactionReceipt(ctx, tx.Hash())
	if err == nil {
		txInfo.Status = receipt.Status
		txInfo.CumulativeGasUsed = receipt.CumulativeGasUsed
		for _, log := range receipt.Logs {
			logInfo := &LogInfo{
				Address:    log.Address.Hex(),
				DataLength: len(log.Data),
			}
			for _, topic := range log.Topics {
				logInfo.Topics = append(logInfo.Topics, topic.Hex())
			}

			// ERC20/ERC721 Transfer event signature
			const transferEvent = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
			if len(log.Topics) == 3 && log.Topics[0].Hex() == transferEvent && len(log.Data) == 32 {
				fromAddr := "0x" + log.Topics[1].Hex()[26:]
				toAddr := "0x" + log.Topics[2].Hex()[26:]
				amount := new(big.Int).SetBytes(log.Data)
				logInfo.Topics = append(logInfo.Topics,
					fmt.Sprintf("Transfer: from %s to %s value %s", fromAddr, toAddr, amount.String()))
			}

			// ERC1155 TransferSingle event
			const transfer1155SingleEvent = "0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62"
			if len(log.Topics) == 4 && log.Topics[0].Hex() == transfer1155SingleEvent && len(log.Data) == 64 {
				operator := "0x" + log.Topics[1].Hex()[26:]
				fromAddr := "0x" + log.Topics[2].Hex()[26:]
				toAddr := "0x" + log.Topics[3].Hex()[26:]
				id := new(big.Int).SetBytes(log.Data[:32])
				value := new(big.Int).SetBytes(log.Data[32:])
				logInfo.Topics = append(logInfo.Topics,
					fmt.Sprintf("ERC1155 TransferSingle: operator %s from %s to %s id %s value %s", operator, fromAddr, toAddr, id.String(), value.String()))
			}

			// ERC1155 TransferBatch event (decode ids and values arrays)
			const transfer1155BatchEvent = "0x4a39dc06d4c0dbc64b70b1b5fdcf9a43c3b840ecb9c7aafb5c62c0124c6a16e3"
			if len(log.Topics) == 4 && log.Topics[0].Hex() == transfer1155BatchEvent && len(log.Data) >= 64 {
				operator := "0x" + log.Topics[1].Hex()[26:]
				fromAddr := "0x" + log.Topics[2].Hex()[26:]
				toAddr := "0x" + log.Topics[3].Hex()[26:]

				// Decode dynamic arrays for ids and values
				// Data layout: offset_ids (32 bytes) | offset_values (32 bytes) | ids[] | values[]
				if len(log.Data) >= 128 {
					idsOffset := new(big.Int).SetBytes(log.Data[:32]).Int64()
					valuesOffset := new(big.Int).SetBytes(log.Data[32:64]).Int64()
					idsStart := int(idsOffset)
					valuesStart := int(valuesOffset)

					// ids array
					idsLen := new(big.Int).SetBytes(log.Data[idsStart : idsStart+32]).Int64()
					var ids []string
					for i := int64(0); i < idsLen; i++ {
						id := new(big.Int).SetBytes(log.Data[idsStart+32+int(i)*32 : idsStart+32+int(i+1)*32])
						ids = append(ids, id.String())
					}

					// values array
					valuesLen := new(big.Int).SetBytes(log.Data[valuesStart : valuesStart+32]).Int64()
					var values []string
					for i := int64(0); i < valuesLen; i++ {
						val := new(big.Int).SetBytes(log.Data[valuesStart+32+int(i)*32 : valuesStart+32+int(i+1)*32])
						values = append(values, val.String())
					}

					logInfo.Topics = append(logInfo.Topics,
						fmt.Sprintf("ERC1155 TransferBatch: operator %s from %s to %s ids %v values %v", operator, fromAddr, toAddr, ids, values))
				} else {
					logInfo.Topics = append(logInfo.Topics,
						fmt.Sprintf("ERC1155 TransferBatch: operator %s from %s to %s (unable to decode ids/values)", operator, fromAddr, toAddr))
				}
			}

			txInfo.Logs = append(txInfo.Logs, logInfo)
		}
	}

	return txInfo
}
