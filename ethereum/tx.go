package ethereum

import (
	"context"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
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
		Hash:         block.Hash().Hex(),
		ParentHash:   block.ParentHash().Hex(),
		Time:         block.Time(),
		GasUsed:      block.GasUsed(),
		GasLimit:     block.GasLimit(),
		BaseFee:      block.BaseFee(),
		Transactions: transactions,
	}
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
		txInfo.Status = receipt.Status
		txInfo.CumulativeGasUsed = receipt.CumulativeGasUsed
		if logs := receipt.Logs; len(logs) > 0 {
			decodedLogs := make([]any, 0, len(logs))
			for _, log := range logs {
				if decoded := e.decodeLog(log); decoded != nil {
					decodedLogs = append(decodedLogs, decoded)
				}
			}
			txInfo.Logs = decodedLogs
		}
	}
	return txInfo
}

// decodeLog returns a map[string]any with decoded event info using EventSignature.
func (e *Ethereum) decodeLog(log *types.Log) any {
	eventType := ""
	var decodedTopics map[string]any

	if len(log.Topics) > 0 {
		sighash := log.Topics[0].Hex()
		if textSig, ok := e.EventSignature[sighash]; ok {
			eventType = textSig
			// Try to decode topics if ABI info is available
			if e.EventABI != nil {
				if event, ok := e.EventABI[sighash]; ok {
					decodedTopics = decodeIndexedTopics(event, log.Topics)
				}
			}
		}
	}

	return map[string]any{
		"address":       strings.ToLower(log.Address.Hex()),
		"topics":        topicsToStrings(log.Topics),
		"decodedTopics": decodedTopics,
		"data":          hex.EncodeToString(log.Data),
		"eventType":     eventType,
	}
}
func topicsToStrings(topics []common.Hash) []string {
	out := make([]string, len(topics))
	for i, t := range topics {
		out[i] = t.Hex()
	}
	return out
}

// Helper to decode indexed topics using ABI event definition
func decodeIndexedTopics(event abi.Event, topics []common.Hash) map[string]any {
	decoded := make(map[string]any)
	// topics[0] is the event signature hash, skip it
	indexedCount := 0
	for _, arg := range event.Inputs {
		if arg.Indexed {
			if len(topics) > indexedCount+1 {
				switch arg.Type.String() {
				case "address":
					decoded[arg.Name] = common.HexToAddress(topics[indexedCount+1].Hex()).Hex()
				case "uint256", "uint":
					decoded[arg.Name] = new(big.Int).SetBytes(topics[indexedCount+1].Bytes()).String()
				default:
					decoded[arg.Name] = topics[indexedCount+1].Hex()
				}
			}
			indexedCount++
		}
	}
	return decoded
}
