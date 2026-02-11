package fx

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Block struct {
	Number       *big.Int    `json:"number"`
	Hash         common.Hash `json:"hash"`
	ParentHash   common.Hash `json:"parentHash"`
	Timestamp    uint64      `json:"timestamp"`
	TxCount      int         `json:"txCount"`
	GasLimit     uint64      `json:"gasLimit"`
	GasUsed      uint64      `json:"gasUsed"`
	BaseFee      *big.Int    `json:"baseFeePerGas,omitempty"`
	Transactions []*Receipt  `json:"transactions"`
}

type Receipt struct {
	TxHash            common.Hash     `json:"transactionHash"`
	From              common.Address  `json:"from"`
	To                *common.Address `json:"to,omitempty"`
	Value             *big.Int        `json:"value,omitempty"`
	Input             []byte          `json:"input,omitempty"`
	Status            uint64          `json:"status"`
	GasUsed           uint64          `json:"gasUsed"`
	EffectiveGasPrice *big.Int        `json:"effectiveGasPrice"`
	ContractAddress   *common.Address `json:"contractAddress,omitempty"`
	Logs              []*Log          `json:"logs,omitempty"`
}

type Log struct {
	Address common.Address `json:"address"`
	Topics  []common.Hash  `json:"topics"`
	Data    []byte         `json:"data,omitempty"`
	Index   uint           `json:"logIndex"`
	TxIndex uint           `json:"transactionIndex"`
}

func (fx *Fx) Block(number *big.Int) (*Block, error) {
	block, err := fx.Eth.BlockByNumber(fx.Context, number)
	if err != nil {
		return nil, fmt.Errorf("block: %w", err)
	}

	receipts, err := fx.blockReceipts(block.Number())
	if err != nil {
		return nil, err
	}

	ethTxs := block.Transactions()
	signer := types.MakeSigner(fx.Chain, block.Number(), block.Time())
	txs := make([]*Receipt, len(ethTxs))
	for i, tx := range ethTxs {
		r := receipts[i]

		from, _ := types.Sender(signer, tx)

		var contractAddr *common.Address
		if r.ContractAddress != (common.Address{}) {
			contractAddr = &r.ContractAddress
		}

		txs[i] = &Receipt{
			TxHash:            tx.Hash(),
			From:              from,
			To:                tx.To(),
			Value:             tx.Value(),
			Input:             tx.Data(),
			Status:            r.Status,
			GasUsed:           r.GasUsed,
			EffectiveGasPrice: r.EffectiveGasPrice,
			ContractAddress:   contractAddr,
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
		TxCount:      len(txs),
		Transactions: txs,
	}, nil
}

func (fx *Fx) blockReceipts(number *big.Int) ([]*types.Receipt, error) {
	var receipts []*types.Receipt
	arg := "latest"
	if number != nil {
		arg = fmt.Sprintf("0x%x", number)
	}
	if err := fx.Rpc.CallContext(fx.Context, &receipts, "eth_getBlockReceipts", arg); err != nil {
		return nil, fmt.Errorf("block receipts: %w", err)
	}
	return receipts, nil
}

func (fx *Fx) Logs(raw []*types.Log) []*Log {
	logs := make([]*Log, len(raw))
	for i, l := range raw {
		logs[i] = &Log{
			Address: l.Address,
			Topics:  l.Topics,
			Data:    l.Data,
			Index:   l.Index,
			TxIndex: l.TxIndex,
		}
	}
	return logs
}
