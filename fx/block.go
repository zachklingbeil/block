package fx

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

type Block struct {
	Number       *big.Int       `json:"number"`
	Hash         common.Hash    `json:"hash"`
	ParentHash   common.Hash    `json:"parentHash"`
	Timestamp    uint64         `json:"timestamp"`
	TxCount      int            `json:"txCount"`
	GasLimit     uint64         `json:"gasLimit"`
	GasUsed      uint64         `json:"gasUsed"`
	BaseFee      *big.Int       `json:"baseFeePerGas,omitempty"`
	Transactions []*Transaction `json:"transactions"`
}

type Transaction struct {
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

	Status            uint64          `json:"status"`
	GasUsed           uint64          `json:"gasUsed"`
	CumulativeGasUsed uint64          `json:"cumulativeGasUsed"`
	EffectiveGasPrice *big.Int        `json:"effectiveGasPrice"`
	ContractAddress   *common.Address `json:"contractAddress,omitempty"`
	BlobGasUsed       uint64          `json:"blobGasUsed,omitempty"`
	BlobGasPrice      *big.Int        `json:"blobGasPrice,omitempty"`
	Logs              []*Log          `json:"logs,omitempty"`
}

func (fx *Fx) Block(number *big.Int) (*Block, error) {
	block, err := fx.Eth.BlockByNumber(fx.Context, number)
	if err != nil {
		return nil, fmt.Errorf("block: %w", err)
	}

	ethTxs := block.Transactions()
	n := len(ethTxs)

	receipts, err := fx.receipts(ethTxs)
	if err != nil {
		return nil, err
	}

	txs := make([]*Transaction, n)
	for i, tx := range ethTxs {
		r := receipts[i]

		var contractAddr *common.Address
		if r.ContractAddress != (common.Address{}) {
			contractAddr = &r.ContractAddress
		}

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
			Logs:              fx.Logs(r.Logs),
		}
	}

	return &Block{
		Number:       block.Number(),
		Hash:         block.Hash(),
		ParentHash:   block.ParentHash(),
		Timestamp:    block.Time(),
		GasLimit:     block.GasLimit(),
		GasUsed:      block.GasUsed(),
		BaseFee:      block.BaseFee(),
		TxCount:      n,
		Transactions: txs,
	}, nil
}

func (fx *Fx) receipts(txs types.Transactions) ([]*types.Receipt, error) {
	n := len(txs)
	if n == 0 {
		return nil, nil
	}
	receipts := make([]*types.Receipt, n)
	batch := make([]rpc.BatchElem, n)
	for i, tx := range txs {
		receipts[i] = &types.Receipt{}
		batch[i] = rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []any{tx.Hash()},
			Result: receipts[i],
		}
	}
	if err := fx.Rpc.BatchCallContext(fx.Context, batch); err != nil {
		return nil, fmt.Errorf("batch receipts: %w", err)
	}
	for i, elem := range batch {
		if elem.Error != nil {
			return nil, fmt.Errorf("receipt[%d]: %w", i, elem.Error)
		}
	}
	return receipts, nil
}
