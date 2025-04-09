package ethereum


import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BlockProcessor is responsible for processing blocks and transactions.
type BlockProcessor struct {
	Client *ethclient.Client
}

// BlockInfo holds information about a block.
type BlockInfo struct {
	Number       uint64
	Hash         string
	ParentHash   string
	Time         uint64
	GasUsed      uint64
	GasLimit     uint64
	BaseFee      *big.Int
	Transactions []*TransactionInfo
}

// TransactionInfo holds information about a transaction.
type TransactionInfo struct {
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
func (bp *BlockProcessor) ProcessBlocks(ctx context.Context, count int) ([]*BlockInfo, error) {
	header, err := bp.Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest header: %w", err)
	}
	latestBlock := header.Number.Uint64()

	var blocks []*BlockInfo
	for blockNum := latestBlock; blockNum > latestBlock-uint64(count); blockNum-- {
		block, err := bp.Client.BlockByNumber(ctx, big.NewInt(int64(blockNum)))
		if err != nil {
			log.Printf("Error fetching block %d: %v", blockNum, err)
			continue
		}
		blockInfo := bp.processBlock(ctx, block)
		blocks = append(blocks, blockInfo)
	}
	return blocks, nil
}

// processBlock processes a single block and returns its information.
func (bp *BlockProcessor) processBlock(ctx context.Context, block *types.Block) *BlockInfo {
	blockInfo := &BlockInfo{
		Number:     block.NumberU64(),
		Hash:       block.Hash().Hex(),
		ParentHash: block.ParentHash().Hex(),
		Time:       block.Time(),
		GasUsed:    block.GasUsed(),
		GasLimit:   block.GasLimit(),
		BaseFee:    block.BaseFee(),
	}

	for _, tx := range block.Transactions() {
		txInfo := bp.processTransaction(ctx, tx)
		blockInfo.Transactions = append(blockInfo.Transactions, txInfo)
	}
	return blockInfo
}

// processTransaction processes a single transaction and returns its information.
func (bp *BlockProcessor) processTransaction(ctx context.Context, tx *types.Transaction) *TransactionInfo {
	txInfo := &TransactionInfo{
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
	receipt, err := bp.Client.TransactionReceipt(ctx, tx.Hash())
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
