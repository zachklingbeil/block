package ethereum

import (
	"context"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Block struct {
	Number uint64
	// Hash         string
	// ParentHash   string
	Time         uint64
	GasUsed      uint64
	GasLimit     uint64
	BaseFee      *big.Int
	Transactions []*Transactions
}

type Transactions struct {
	// Hash     string
	From     string
	To       string
	Value    *big.Int
	Gas      uint64
	GasPrice *big.Int
	Nonce    uint64
	// DataLength        int
	// Type              uint8
	// Status            uint64
	CumulativeGasUsed uint64
	FunctionSignature string `json:"functionSignature,omitempty"`
	Logs              []any  `json:"logs,omitempty"`
}

func (e *Ethereum) processBlock(ctx context.Context, block *types.Block) *Block {
	signer := e.Signer(block.Number(), block.Time())
	txs := block.Transactions()
	transactions := make([]*Transactions, 0, len(txs))
	for _, tx := range txs {
		if txInfo := e.processTransaction(ctx, tx, signer); txInfo != nil {
			transactions = append(transactions, txInfo)
		}
	}
	return &Block{
		Number: block.NumberU64(),
		// Hash:         block.Hash().Hex(),
		// ParentHash:   block.ParentHash().Hex(),
		Time:         block.Time(),
		GasUsed:      block.GasUsed(),
		GasLimit:     block.GasLimit(),
		BaseFee:      block.BaseFee(),
		Transactions: transactions,
	}
}

func (e *Ethereum) processTransaction(ctx context.Context, tx *types.Transaction, signer types.Signer) *Transactions {
	txInfo := &Transactions{
		// Hash:     strings.ToLower(tx.Hash().Hex()),
		Value:    tx.Value(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Nonce:    tx.Nonce(),
		// DataLength: len(tx.Data()),
		// Type:       tx.Type(),
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

	// Set function signature if available using Signature map
	if data := tx.Data(); len(data) >= 4 {
		selector := "0x" + hex.EncodeToString(data[:4])
		if textSig, ok := e.Signature[selector]; ok {
			txInfo.FunctionSignature = textSig
		}
	}

	// Attach receipt info if available
	if receipt, err := e.Factory.Eth.TransactionReceipt(ctx, tx.Hash()); err == nil {
		// txInfo.Status = receipt.Status
		txInfo.CumulativeGasUsed = receipt.CumulativeGasUsed

		// Process logs here
		var decodedLogs []any
		for _, log := range receipt.Logs {
			if decoded := e.decodeLog(log); decoded != nil {
				decodedLogs = append(decodedLogs, decoded)
			}
		}
		txInfo.Logs = decodedLogs
	}
	return txInfo
}

func (e *Ethereum) decodeLog(log *types.Log) any {
	if len(log.Topics) == 0 {
		return nil
	}
	sighash := log.Topics[0].Hex()
	event, ok := e.EventABI[sighash]
	if !ok {
		return nil
	}

	result := make(map[string]any, len(event.Inputs)+1)
	idx := 1

	// Initialize all fields to nil
	for _, arg := range event.Inputs {
		result[arg.Name] = nil
	}

	// Decode indexed fields
	for _, arg := range event.Inputs {
		if arg.Indexed && len(log.Topics) > idx {
			switch arg.Type.String() {
			case "address":
				result[arg.Name] = common.HexToAddress(log.Topics[idx].Hex()).Hex()
			case "uint256", "uint":
				result[arg.Name] = new(big.Int).SetBytes(log.Topics[idx].Bytes()).String()
			case "bool":
				result[arg.Name] = log.Topics[idx].Big().Cmp(big.NewInt(0)) != 0
			default:
				result[arg.Name] = log.Topics[idx].Hex()
			}
			idx++
		}
	}

	// Decode non-indexed fields
	nonIndexed := event.Inputs.NonIndexed()
	if len(nonIndexed) > 0 && len(log.Data) > 0 {
		if values, err := nonIndexed.Unpack(log.Data); err == nil {
			for i, arg := range nonIndexed {
				switch arg.Type.String() {
				case "address":
					result[arg.Name] = values[i].(common.Address).Hex()
				case "uint256", "uint":
					result[arg.Name] = values[i].(*big.Int).String()
				case "bool":
					result[arg.Name] = values[i].(bool)
				default:
					result[arg.Name] = values[i]
				}
			}
		}
	}

	// Use thread-safe accessor for event signature
	if eventType, ok := e.GetEventSignature(sighash); ok {
		result["eventType"] = eventType
	} else {
		result["eventType"] = event.Name
	}

	return result
}
