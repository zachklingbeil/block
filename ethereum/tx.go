package ethereum

import (
	"context"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Raw struct {
	Number       uint64
	Hash         string
	ParentHash   string
	Time         uint64
	GasUsed      uint64
	GasLimit     uint64
	BaseFee      *big.Int
	Transactions []*Transactions
}

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
	FunctionSignature string `json:"functionSignature,omitempty"`
	Logs              []any  `json:"logs,omitempty"`
}

func (e *Ethereum) processBlock(ctx context.Context, block *types.Block) *Raw {
	blockInfo := &Raw{
		Number:     block.NumberU64(),
		Hash:       block.Hash().Hex(),
		ParentHash: block.ParentHash().Hex(),
		Time:       block.Time(),
		GasUsed:    block.GasUsed(),
		GasLimit:   block.GasLimit(),
		BaseFee:    block.BaseFee(),
	}
	signer := e.Signer(block.Number(), block.Time())
	for _, tx := range block.Transactions() {
		txInfo := e.processTransaction(ctx, tx, signer)
		blockInfo.Transactions = append(blockInfo.Transactions, txInfo)
	}
	return blockInfo
}

func (e *Ethereum) processTransaction(ctx context.Context, tx *types.Transaction, signer types.Signer) *Transactions {
	txInfo := &Transactions{
		Hash:       strings.ToLower(tx.Hash().Hex()),
		Value:      tx.Value(),
		Gas:        tx.Gas(),
		GasPrice:   tx.GasPrice(),
		Nonce:      tx.Nonce(),
		DataLength: len(tx.Data()),
		Type:       tx.Type(),
	}
	if addr, err := types.Sender(signer, tx); err == nil {
		txInfo.From = strings.ToLower(addr.Hex())
	}
	if tx.To() == nil {
		txInfo.To = "Contract Creation"
	} else {
		txInfo.To = strings.ToLower(tx.To().Hex())
	}
	if data := tx.Data(); len(data) >= 4 && e.HexToText != nil {
		selector := "0x" + hex.EncodeToString(data[:4])
		if textSig, ok := e.GetHexText(selector); ok {
			txInfo.FunctionSignature = textSig
		}
	}
	if receipt, err := e.Factory.Eth.TransactionReceipt(ctx, tx.Hash()); err == nil {
		txInfo.Status = receipt.Status
		txInfo.CumulativeGasUsed = receipt.CumulativeGasUsed
		for _, log := range receipt.Logs {
			if decoded := e.decodeLog(log); decoded != nil {
				txInfo.Logs = append(txInfo.Logs, decoded)
			}
		}
	}
	return txInfo
}

// decodeLog returns a map[string]any with decoded event info using HexToText.
func (e *Ethereum) decodeLog(log *types.Log) any {
	eventType := ""
	if len(log.Topics) > 0 && e.HexToText != nil {
		if textSig, ok := e.GetHexText(log.Topics[0].Hex()); ok {
			eventType = textSig
		}
	}
	result := map[string]any{
		"address":   strings.ToLower(log.Address.Hex()),
		"topics":    topicsToStrings(log.Topics),
		"data":      hex.EncodeToString(log.Data),
		"eventType": eventType,
	}
	return result
}

func topicsToStrings(topics []common.Hash) []string {
	out := make([]string, len(topics))
	for i, t := range topics {
		out[i] = t.Hex()
	}
	return out
}
