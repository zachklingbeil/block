package fx

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Block struct {
	Number       *big.Int       `json:"number"`
	Hash         common.Hash    `json:"hash"`
	Timestamp    uint64         `json:"timestamp"`
	GasLimit     uint64         `json:"gasLimit"`
	GasUsed      uint64         `json:"gasUsed"`
	BaseFee      *big.Int       `json:"baseFeePerGas,omitempty"`
	Transactions []*Transaction `json:"transactions"`
}

type Transaction struct {
	TxHash            common.Hash     `json:"hash"`
	TxIndex           uint            `json:"index"`
	From              common.Address  `json:"from"`
	To                *common.Address `json:"to,omitempty"`
	Value             *big.Int        `json:"value,omitempty"`
	Input             string          `json:"input,omitempty"`
	Status            uint64          `json:"status"`
	Gas               uint64          `json:"gas"`
	EffectiveGasPrice *big.Int        `json:"gasPrice"`
	ContractAddress   *common.Address `json:"contractAddress,omitempty"`
	Logs              []*types.Log    `json:"logs,omitempty"`
	Decoded           *Decoded        `json:"decoded,omitempty"`
}

type Contract struct {
	abi.ABI
}

type Decoded struct {
	Call   *Call  `json:"call,omitempty"`
	Events []Call `json:"events,omitempty"`
	Error  *Call  `json:"error,omitempty"`
}

type Call struct {
	Contract common.Address `json:"contract"`
	Name     string         `json:"name"`
	Sig      string         `json:"sig"`
	Values   map[string]any `json:"values"`
}

func (fx *Fx) Block(number *big.Int) (*Block, error) {
	block, err := fx.Eth.BlockByNumber(fx.Context, number)
	if err != nil {
		return nil, fmt.Errorf("block: %w", err)
	}

	var receipts []*types.Receipt
	if err := fx.Rpc.CallContext(fx.Context, &receipts, "eth_getBlockReceipts", fmt.Sprintf("0x%x", block.Number())); err != nil {
		return nil, fmt.Errorf("block receipts: %w", err)
	}

	signer := types.MakeSigner(fx.Chain, block.Number(), block.Time())
	txs := make([]*Transaction, len(block.Transactions()))

	for i, tx := range block.Transactions() {
		from, _ := types.Sender(signer, tx)

		t := &Transaction{
			TxHash:  tx.Hash(),
			TxIndex: uint(i),
			From:    from,
			To:      tx.To(),
			Value:   tx.Value(),
			Input:   hex.EncodeToString(tx.Data()),
		}

		if i < len(receipts) {
			r := receipts[i]
			t.Status = r.Status
			t.Gas = r.GasUsed
			t.EffectiveGasPrice = r.EffectiveGasPrice
			if r.ContractAddress != (common.Address{}) {
				addr := r.ContractAddress
				t.ContractAddress = &addr
			}
			t.Logs = r.Logs
		}
		txs[i] = t
	}

	return &Block{
		Number:       block.Number(),
		Hash:         block.Hash(),
		Timestamp:    block.Time(),
		GasLimit:     block.GasLimit(),
		GasUsed:      block.GasUsed(),
		BaseFee:      block.BaseFee(),
		Transactions: txs,
	}, nil
}

func (fx *Fx) Output(b *Block) *Block {
	for _, tx := range b.Transactions {
		d := &Decoded{}

		if tx.To != nil && len(tx.Input) >= 8 {
			input, _ := hex.DecodeString(tx.Input)
			d.Call = fx.Method(*tx.To, input)
		}

		for _, l := range tx.Logs {
			if len(l.Topics) == 0 {
				continue
			}
			if c := fx.Event(l.Address, l); c != nil {
				d.Events = append(d.Events, *c)
			}
		}

		if d.Call != nil || len(d.Events) > 0 {
			tx.Decoded = d
		}
	}
	return b
}

func (fx *Fx) Method(addr common.Address, input []byte) *Call {
	contract, ok := fx.Contracts[addr]
	if !ok || len(input) < 4 {
		return nil
	}
	m, err := contract.MethodById(input[:4])
	if err != nil {
		return nil
	}
	values := make(map[string]any)
	if err := m.Inputs.UnpackIntoMap(values, input[4:]); err != nil {
		return nil
	}
	return &Call{
		Contract: addr,
		Name:     m.Name,
		Sig:      m.Sig,
		Values:   values,
	}
}

func (fx *Fx) Event(addr common.Address, l *types.Log) *Call {
	contract, ok := fx.Contracts[addr]
	if !ok || len(l.Topics) == 0 {
		return nil
	}
	e, err := contract.EventByID(l.Topics[0])
	if err != nil {
		return nil
	}
	values := make(map[string]any)

	indexed := make(abi.Arguments, 0)
	nonIndexed := make(abi.Arguments, 0)
	for _, arg := range e.Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		} else {
			nonIndexed = append(nonIndexed, arg)
		}
	}

	if err := abi.ParseTopicsIntoMap(values, indexed, l.Topics[1:]); err != nil {
		return nil
	}

	if len(l.Data) > 0 {
		if err := nonIndexed.UnpackIntoMap(values, l.Data); err != nil {
			return nil
		}
	}

	return &Call{
		Contract: addr,
		Name:     e.Name,
		Sig:      e.Sig,
		Values:   values,
	}
}

func (fx *Fx) Error(addr common.Address, data []byte) *Call {
	contract, ok := fx.Contracts[addr]
	if !ok || len(data) < 4 {
		return nil
	}
	var sig [4]byte
	copy(sig[:], data[:4])
	e, err := contract.ErrorByID(sig)
	if err != nil {
		return nil
	}
	values := make(map[string]any)
	if len(data) > 4 {
		if err := e.Inputs.UnpackIntoMap(values, data[4:]); err != nil {
			return nil
		}
	}
	return &Call{
		Contract: addr,
		Name:     e.Name,
		Sig:      e.Sig,
		Values:   values,
	}
}
