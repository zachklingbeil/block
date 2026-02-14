package fx

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Raw struct {
	Block    *types.Block
	Receipts []*types.Receipt
}

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
	TxHash          common.Hash     `json:"hash"`
	TxIndex         uint            `json:"index"`
	From            common.Address  `json:"from"`
	To              *common.Address `json:"to,omitempty"`
	Value           *big.Int        `json:"value,omitempty"`
	Status          uint64          `json:"status"`
	Gas             uint64          `json:"gas"`
	GasPrice        *big.Int        `json:"gasPrice"`
	ContractAddress *common.Address `json:"contractAddress,omitempty"`
	Method          *Event          `json:"method,omitempty"`
	Events          []Event         `json:"events,omitempty"`
}

func (fx *Fx) Block(number *big.Int) (*Block, error) {
	raw, err := fx.Source(number)
	if err != nil {
		return nil, fmt.Errorf("collect: %w", err)
	}
	fx.Resolve(raw)
	block := fx.Transform(raw)
	for _, t := range block.Transactions {
		fx.Record(t)
	}
	return block, nil
}

func (fx *Fx) Source(number *big.Int) (*Raw, error) {
	block, err := fx.Eth.BlockByNumber(fx.Context, number)
	if err != nil {
		return nil, fmt.Errorf("block: %w", err)
	}
	var receipts []*types.Receipt
	if err := fx.Rpc.CallContext(fx.Context, &receipts, "eth_getBlockReceipts", fmt.Sprintf("0x%x", block.Number())); err != nil {
		return nil, fmt.Errorf("block receipts: %w", err)
	}
	return &Raw{Block: block, Receipts: receipts}, nil
}

func (fx *Fx) Resolve(raw *Raw) {
	for addr := range fx.addresses(raw) {
		fx.GetContract(addr)
	}
}

func (fx *Fx) Transform(raw *Raw) *Block {
	signer := types.MakeSigner(fx.Chain, raw.Block.Number(), raw.Block.Time())
	txs := make([]*Transaction, len(raw.Block.Transactions()))
	for i, tx := range raw.Block.Transactions() {
		from, _ := types.Sender(signer, tx)
		txs[i] = fx.transaction(from, tx, fx.receipt(raw.Receipts, i), tx.Data(), i)
	}
	return &Block{
		Number:       raw.Block.Number(),
		Hash:         raw.Block.Hash(),
		Timestamp:    raw.Block.Time(),
		GasLimit:     raw.Block.GasLimit(),
		GasUsed:      raw.Block.GasUsed(),
		BaseFee:      raw.Block.BaseFee(),
		Transactions: txs,
	}
}

func (fx *Fx) Record(t *Transaction) {
	if t.Method == nil || t.To == nil {
		return
	}
	c, ok := fx.Contracts[*t.To]
	if !ok {
		return
	}
	sel := t.Method.Selector
	events := fx.templatize(t)
	outcomes := c.Outcomes[sel]
	for i := range outcomes {
		if matchOutcome(&outcomes[i], events) {
			if t.Status == 1 {
				outcomes[i].Success++
			} else {
				outcomes[i].Fail++
			}
			return
		}
	}
	o := Outcome{Events: events}
	if t.Status == 1 {
		o.Success = 1
	} else {
		o.Fail = 1
	}
	c.Outcomes[sel] = append(outcomes, o)
}

// addresses returns every contract address referenced in the raw block.
func (fx *Fx) addresses(raw *Raw) map[common.Address]struct{} {
	seen := make(map[common.Address]struct{})
	for i, tx := range raw.Block.Transactions() {
		if tx.To() != nil && len(tx.Data()) >= 4 {
			seen[*tx.To()] = struct{}{}
		}
		if r := fx.receipt(raw.Receipts, i); r != nil {
			if r.ContractAddress != (common.Address{}) {
				seen[r.ContractAddress] = struct{}{}
			}
			for _, l := range r.Logs {
				seen[l.Address] = struct{}{}
			}
		}
	}
	return seen
}

// receipt safely indexes into the receipts slice.
func (fx *Fx) receipt(receipts []*types.Receipt, i int) *types.Receipt {
	if i < len(receipts) {
		return receipts[i]
	}
	return nil
}

// transaction builds a Transaction from its parts.
func (fx *Fx) transaction(from common.Address, tx *types.Transaction, r *types.Receipt, data []byte, index int) *Transaction {
	var value *big.Int
	if tx.Value() != nil && tx.Value().Sign() != 0 {
		value = tx.Value()
	}
	t := &Transaction{
		TxHash:  tx.Hash(),
		TxIndex: uint(index),
		From:    from,
		To:      tx.To(),
		Value:   value,
	}
	if r == nil {
		return t
	}
	t.Status = r.Status
	t.Gas = r.GasUsed
	t.GasPrice = r.EffectiveGasPrice
	if r.ContractAddress != (common.Address{}) {
		addr := r.ContractAddress
		t.ContractAddress = &addr
	}
	if t.To != nil && len(data) >= 4 {
		t.Method = fx.decode(*t.To, nil, data)
	}
	for _, l := range r.Logs {
		if len(l.Topics) == 0 {
			continue
		}
		if d := fx.decode(l.Address, l.Topics, l.Data); d != nil {
			t.Events = append(t.Events, *d)
		}
	}
	return t
}

// templatize maps event params to method input names for outcome storage.
func (fx *Fx) templatize(t *Transaction) []Event {
	var events []Event
	for _, e := range t.Events {
		params := make(map[string]any, len(e.Params))
		for ep, ev := range e.Params {
			var mapped any
			for mp, mv := range t.Method.Params {
				if valuesEqual(ev, mv) {
					mapped = mp
					break
				}
			}
			params[ep] = mapped
		}
		events = append(events, Event{
			Contract: e.Contract,
			Selector: e.Selector,
			Action:   e.Action,
			Params:   params,
		})
	}
	return events
}
