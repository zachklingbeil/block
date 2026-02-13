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
	Method            *Decoded        `json:"method,omitempty"`
	Events            []Decoded       `json:"events,omitempty"`
	Error             *Decoded        `json:"error,omitempty"`
}

type Decoded struct {
	Contract common.Address `json:"contract"`
	Name     string         `json:"name"`
	Sig      string         `json:"sig"`
	Values   map[string]any `json:"values"`
}

func (fx *Fx) Block(number *big.Int) (*Block, error) {
	block, receipts, addrs, err := fx.Input(number)
	if err != nil {
		return nil, fmt.Errorf("input: %w", err)
	}

	for _, addr := range addrs {
		fx.GetABI(addr)
	}
	return fx.Output(block, receipts), nil
}

func (fx *Fx) Input(number *big.Int) (*types.Block, []*types.Receipt, []common.Address, error) {
	block, err := fx.Eth.BlockByNumber(fx.Context, number)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("block: %w", err)
	}

	var receipts []*types.Receipt
	if err := fx.Rpc.CallContext(fx.Context, &receipts, "eth_getBlockReceipts", fmt.Sprintf("0x%x", block.Number())); err != nil {
		return nil, nil, nil, fmt.Errorf("block receipts: %w", err)
	}

	seen := make(map[common.Address]struct{})
	for i, tx := range block.Transactions() {
		if tx.To() != nil && len(tx.Data()) >= 4 {
			seen[*tx.To()] = struct{}{}
		}
		if i < len(receipts) {
			r := receipts[i]
			if r.ContractAddress != (common.Address{}) {
				seen[r.ContractAddress] = struct{}{}
			}
			for _, l := range r.Logs {
				seen[l.Address] = struct{}{}
			}
		}
	}
	addrs := make([]common.Address, 0, len(seen))
	for addr := range seen {
		addrs = append(addrs, addr)
	}

	return block, receipts, addrs, nil
}

func (fx *Fx) Output(block *types.Block, receipts []*types.Receipt) *Block {
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

			if t.To != nil && len(t.Input) >= 8 {
				input, _ := hex.DecodeString(t.Input)
				t.Method = fx.method(*t.To, input)
			}

			for _, l := range r.Logs {
				if len(l.Topics) == 0 {
					continue
				}
				if d := fx.event(l.Address, l.Topics, l.Data); d != nil {
					t.Events = append(t.Events, *d)
				}
			}
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
	}
}

func (fx *Fx) method(addr common.Address, input []byte) *Decoded {
	a, ok := fx.Contracts[addr]
	if !ok || len(input) < 4 {
		return nil
	}
	m, err := a.MethodById(input[:4])
	if err != nil {
		return nil
	}
	values := make(map[string]any)
	if err := m.Inputs.UnpackIntoMap(values, input[4:]); err != nil {
		return nil
	}
	return &Decoded{
		Contract: addr,
		Name:     m.Name,
		Sig:      m.Sig,
		Values:   values,
	}
}

func (fx *Fx) event(addr common.Address, topics []common.Hash, data []byte) *Decoded {
	a, ok := fx.Contracts[addr]
	if !ok || len(topics) == 0 {
		return nil
	}
	e, err := a.EventByID(topics[0])
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

	if err := abi.ParseTopicsIntoMap(values, indexed, topics[1:]); err != nil {
		return nil
	}

	if len(data) > 0 {
		if err := nonIndexed.UnpackIntoMap(values, data); err != nil {
			return nil
		}
	}

	return &Decoded{
		Contract: addr,
		Name:     e.Name,
		Sig:      e.Sig,
		Values:   values,
	}
}

func (fx *Fx) error(addr common.Address, data []byte) *Decoded {
	a, ok := fx.Contracts[addr]
	if !ok || len(data) < 4 {
		return nil
	}
	var sig [4]byte
	copy(sig[:], data[:4])
	e, err := a.ErrorByID(sig)
	if err != nil {
		return nil
	}
	values := make(map[string]any)
	if len(data) > 4 {
		if err := e.Inputs.UnpackIntoMap(values, data[4:]); err != nil {
			return nil
		}
	}
	return &Decoded{
		Contract: addr,
		Name:     e.Name,
		Sig:      e.Sig,
		Values:   values,
	}
}
