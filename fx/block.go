package fx

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/timefactoryio/block/zero/proto/sigprovider"
)

type Block struct {
	Number       *big.Int       `json:"number"`
	Hash         common.Hash    `json:"hash"`
	ParentHash   common.Hash    `json:"parentHash"`
	Timestamp    uint64         `json:"timestamp"`
	TxCount      uint           `json:"txCount"`
	GasLimit     uint64         `json:"gasLimit"`
	GasUsed      uint64         `json:"gasUsed"`
	BaseFee      *big.Int       `json:"baseFeePerGas,omitempty"`
	Transactions []*Transaction `json:"transactions"`
}

// Transaction pairs the intent with the outcome.
type Transaction struct {
	// Intent
	Hash      common.Hash     `json:"hash"`
	Nonce     uint64          `json:"nonce"`
	To        *common.Address `json:"to,omitempty"`
	Value     *big.Int        `json:"value,omitempty"`
	Gas       uint64          `json:"gas"`
	GasPrice  *big.Int        `json:"gasPrice,omitempty"`
	GasTipCap *big.Int        `json:"maxPriorityFeePerGas,omitempty"`
	GasFeeCap *big.Int        `json:"maxFeePerGas,omitempty"`
	Data      []byte          `json:"input,omitempty"`
	Type      uint8           `json:"type"`
	ChainID   *big.Int        `json:"chainId,omitempty"`

	// Outcome
	Status            uint64          `json:"status"`
	GasUsed           uint64          `json:"gasUsed"`
	CumulativeGasUsed uint64          `json:"cumulativeGasUsed"`
	EffectiveGasPrice *big.Int        `json:"effectiveGasPrice"`
	ContractAddress   *common.Address `json:"contractAddress,omitempty"`
	BlobGasUsed       uint64          `json:"blobGasUsed,omitempty"`
	BlobGasPrice      *big.Int        `json:"blobGasPrice,omitempty"`
	Logs              []*Log          `json:"logs,omitempty"`
}

// Log is a contract event — the economic activity.
type Log struct {
	Address common.Address `json:"address"`
	Topics  []common.Hash  `json:"topics"`
	Data    []byte         `json:"data,omitempty"`
	Index   uint           `json:"logIndex"`
	Removed bool           `json:"removed,omitempty"`

	// Decoded
	Event string `json:"event,omitempty"`
	Args  []*Arg `json:"args,omitempty"`
}

// Arg is a decoded event argument.
type Arg struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Indexed bool   `json:"indexed,omitempty"`
	Value   string `json:"value"`
}

func toLogs(logs []*types.Log) []*Log {
	out := make([]*Log, len(logs))
	for i, l := range logs {
		out[i] = &Log{
			Address: l.Address,
			Topics:  l.Topics,
			Data:    l.Data,
			Index:   l.Index,
			Removed: l.Removed,
		}
	}
	return out
}

func toArgs(inputs []*sigprovider.Argument) []*Arg {
	if len(inputs) == 0 {
		return nil
	}
	args := make([]*Arg, len(inputs))
	for i, in := range inputs {
		args[i] = &Arg{
			Name:    in.GetName(),
			Type:    in.GetType(),
			Indexed: in.GetIndexed(),
			Value:   in.GetValue(),
		}
	}
	return args
}

// decodeEvents batch-decodes all logs via sigprovider.
func (fx *Fx) decodeEvents(logs []*Log) {
	if len(logs) == 0 {
		return
	}

	reqs := make([]*sigprovider.GetEventAbiRequest, len(logs))
	for i, l := range logs {
		topics := make([]string, len(l.Topics))
		for j, t := range l.Topics {
			topics[j] = t.Hex()
		}
		reqs[i] = &sigprovider.GetEventAbiRequest{
			Data:   "0x" + hex.EncodeToString(l.Data),
			Topics: strings.Join(topics, ","),
		}
	}

	resp, err := fx.Sig.BatchGetEventAbis(fx.Context, &sigprovider.BatchGetEventAbisRequest{
		Requests: reqs,
	})
	if err != nil {
		return
	}

	for i, r := range resp.GetResponses() {
		if i >= len(logs) {
			break
		}
		abis := r.GetAbi()
		if len(abis) == 0 {
			continue
		}
		// Take the first match
		logs[i].Event = abis[0].GetName()
		logs[i].Args = toArgs(abis[0].GetInputs())
	}
}

func (fx *Fx) Block(number *big.Int) (*Block, error) {
	block, err := fx.Eth.BlockByNumber(fx.Context, number)
	if err != nil {
		return nil, fmt.Errorf("block: %w", err)
	}

	txs := make([]*Transaction, len(block.Transactions()))
	var allLogs []*Log

	for i, tx := range block.Transactions() {
		r, err := fx.Eth.TransactionReceipt(fx.Context, tx.Hash())
		if err != nil {
			return nil, fmt.Errorf("receipt[%d]: %w", i, err)
		}

		var contractAddr *common.Address
		if r.ContractAddress != (common.Address{}) {
			contractAddr = &r.ContractAddress
		}

		logs := toLogs(r.Logs)
		allLogs = append(allLogs, logs...)

		txs[i] = &Transaction{
			Hash:      tx.Hash(),
			Nonce:     tx.Nonce(),
			To:        tx.To(),
			Value:     tx.Value(),
			Gas:       tx.Gas(),
			GasPrice:  tx.GasPrice(),
			GasTipCap: tx.GasTipCap(),
			GasFeeCap: tx.GasFeeCap(),
			Data:      tx.Data(),
			Type:      tx.Type(),
			ChainID:   tx.ChainId(),

			Status:            r.Status,
			GasUsed:           r.GasUsed,
			CumulativeGasUsed: r.CumulativeGasUsed,
			EffectiveGasPrice: r.EffectiveGasPrice,
			ContractAddress:   contractAddr,
			BlobGasUsed:       r.BlobGasUsed,
			BlobGasPrice:      r.BlobGasPrice,
			Logs:              logs,
		}
	}

	fx.decodeEvents(allLogs)

	return &Block{
		Number:       block.Number(),
		Hash:         block.Hash(),
		ParentHash:   block.ParentHash(),
		Timestamp:    block.Time(),
		GasLimit:     block.GasLimit(),
		GasUsed:      block.GasUsed(),
		BaseFee:      block.BaseFee(),
		TxCount:      uint(len(txs)),
		Transactions: txs,
	}, nil
}
