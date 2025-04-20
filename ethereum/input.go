package ethereum

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

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
func (e *Ethereum) ProcessBlocks(ctx context.Context, count int) ([]*Block, error) {
	header, err := e.Factory.Eth.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest header: %w", err)
	}
	latestBlock := header.Number.Uint64()

	var blocks []*Block
	for blockNum := latestBlock; blockNum > latestBlock-uint64(count); blockNum-- {
		block, err := e.Factory.Eth.BlockByNumber(ctx, big.NewInt(int64(blockNum)))
		if err != nil {
			log.Printf("Error fetching block %d: %v", blockNum, err)
			continue
		}
		blockInfo := e.processBlock(ctx, block)
		blocks = append(blocks, blockInfo)
	}
	return blocks, nil
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
	signer := types.NewEIP155Signer(tx.ChainId())
	from, err := types.Sender(signer, tx)
	if err == nil {
		txInfo.From = from.Hex()
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
			txInfo.Logs = append(txInfo.Logs, logInfo)
		}
	}

	return txInfo
}
