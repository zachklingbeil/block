package ethereum

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
)

type Raw struct {
	Number       uint64          `json:"number,omitempty"`
	Time         uint64          `json:"time,omitempty"`
	GasUsed      uint64          `json:"gasUsed,omitempty"`
	GasLimit     uint64          `json:"gasLimit,omitempty"`
	BaseFee      *big.Int        `json:"baseFee,omitempty"`
	Transactions []*Transactions `json:"transactions,omitempty"`
}

func (e *Ethereum) processBlock(ctx context.Context, block *types.Block) *Raw {
	signer := e.Signer(block.Number(), block.Time())
	txs := block.Transactions()
	transactions := make([]*Transactions, 0, len(txs))
	for _, tx := range txs {
		if txInfo := e.processTransaction(ctx, tx, signer); txInfo != nil {
			transactions = append(transactions, txInfo)
		}
	}
	return &Raw{
		Number:       block.NumberU64(),
		Time:         block.Time(),
		GasUsed:      block.GasUsed(),
		GasLimit:     block.GasLimit(),
		BaseFee:      block.BaseFee(),
		Transactions: transactions,
	}
}

func (e *Ethereum) processTransaction(ctx context.Context, tx *types.Transaction, signer types.Signer) *Transactions {
	txInfo := &Transactions{
		Value:    tx.Value(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Nonce:    tx.Nonce(),
	}

	// Set From address
	if addr, err := types.Sender(signer, tx); err == nil {
		txInfo.From = strings.ToLower(addr.Hex())
	}

	// Set To address or contract creation
	if to := tx.To(); to == nil {
		txInfo.To = "Contract Creation"
	} else {
		txInfo.To = strings.ToLower(to.Hex())
	}

	// Populate receipt info (logs, cumulative gas used, etc.)
	if receipt, err := e.Factory.Eth.TransactionReceipt(ctx, tx.Hash()); err == nil {
		e.populateReceiptInfo(txInfo, receipt)
	}
	return txInfo
}
