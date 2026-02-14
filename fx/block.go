package fx

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Raw holds unmodified geth types from a single block fetch.
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
	TxHash            common.Hash     `json:"hash"`
	TxIndex           uint            `json:"index"`
	Type              uint8           `json:"type,omitempty"`
	From              common.Address  `json:"from"`
	To                *common.Address `json:"to,omitempty"`
	Value             *big.Int        `json:"value,omitempty"`
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
	Selector [4]byte        `json:"-"`
	Values   map[string]any `json:"values"`
}

// Block fetches, resolves, transforms, and records a finalized block.
func (fx *Fx) Block(number *big.Int) (*Block, error) {
	raw, err := fx.Collect(number)
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

// Collect fetches the block and receipts from the node. Pure I/O.
func (fx *Fx) Collect(number *big.Int) (*Raw, error) {
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

// Resolve ensures ABIs are loaded for every contract address in the raw block.
func (fx *Fx) Resolve(raw *Raw) {
	seen := make(map[common.Address]struct{})
	for i, tx := range raw.Block.Transactions() {
		if tx.To() != nil && len(tx.Data()) >= 4 {
			seen[*tx.To()] = struct{}{}
		}
		if i < len(raw.Receipts) {
			r := raw.Receipts[i]
			if r.ContractAddress != (common.Address{}) {
				seen[r.ContractAddress] = struct{}{}
			}
			for _, l := range r.Logs {
				seen[l.Address] = struct{}{}
			}
		}
	}
	for addr := range seen {
		fx.GetContract(addr)
	}
}

// Transform converts raw geth types into domain types. Decoding uses fx.Contracts.
func (fx *Fx) Transform(raw *Raw) *Block {
	signer := types.MakeSigner(fx.Chain, raw.Block.Number(), raw.Block.Time())
	txs := make([]*Transaction, len(raw.Block.Transactions()))

	for i, tx := range raw.Block.Transactions() {
		from, _ := types.Sender(signer, tx)

		t := &Transaction{
			TxHash:  tx.Hash(),
			Type:    tx.Type(),
			TxIndex: uint(i),
			From:    from,
			To:      tx.To(),
			Value:   tx.Value(),
		}

		if i < len(raw.Receipts) {
			r := raw.Receipts[i]
			t.Status = r.Status
			t.Gas = r.GasUsed
			t.EffectiveGasPrice = r.EffectiveGasPrice
			if r.ContractAddress != (common.Address{}) {
				addr := r.ContractAddress
				t.ContractAddress = &addr
			}

			if t.To != nil && len(tx.Data()) >= 4 {
				t.Method = fx.method(*t.To, tx.Data())
			}

			for _, l := range r.Logs {
				if len(l.Topics) == 0 {
					continue
				}
				if d := fx.event(l.Address, l.Topics, l.Data); d != nil {
					t.Events = append(t.Events, *d)
				}
			}

			if t.Status == 0 && t.To != nil {
				t.Error = fx.revert(*t.To, tx.Data())
			}
		}

		txs[i] = t
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

	o := Outcome{Status: t.Status, Count: 1}

	if t.Status == 1 {
		for _, e := range t.Events {
			paramMap := make(map[string]string)
			for ep, ev := range e.Values {
				for mp, mv := range t.Method.Values {
					if valuesEqual(ev, mv) {
						paramMap[ep] = mp
					}
				}
			}

			o.Events = append(o.Events, OutcomeEvent{
				Contract: e.Contract,
				Selector: e.Selector,
				ParamMap: paramMap,
			})
		}
	}

	if t.Status == 0 && t.Error != nil {
		o.Error = &OutcomeEvent{
			Contract: t.Error.Contract,
			Selector: t.Error.Selector,
		}
	}

	tmpl, exists := c.Templates[sel]
	if !exists {
		c.Templates[sel] = &Template{
			Method:   sel,
			Outcomes: []Outcome{o},
		}
		return
	}

	for i := range tmpl.Outcomes {
		if matchOutcome(&tmpl.Outcomes[i], &o) {
			tmpl.Outcomes[i].Count++
			return
		}
	}

	tmpl.Outcomes = append(tmpl.Outcomes, o)
}

func matchOutcome(a, b *Outcome) bool {
	if a.Status != b.Status {
		return false
	}

	if a.Status == 0 {
		if a.Error == nil && b.Error == nil {
			return true
		}
		if a.Error == nil || b.Error == nil {
			return false
		}
		return a.Error.Selector == b.Error.Selector
	}

	if len(a.Events) != len(b.Events) {
		return false
	}
	for i := range a.Events {
		if a.Events[i].Selector != b.Events[i].Selector {
			return false
		}
		if a.Events[i].Contract != b.Events[i].Contract {
			return false
		}
	}
	return true
}

func (fx *Fx) method(addr common.Address, input []byte) *Decoded {
	if len(input) < 4 {
		return nil
	}
	c, ok := fx.Contracts[addr]
	if !ok || c.ABI == nil {
		return nil
	}
	m, err := c.ABI.MethodById(input[:4])
	if err != nil {
		return nil
	}
	values := make(map[string]any)
	if err := m.Inputs.UnpackIntoMap(values, input[4:]); err != nil {
		return nil
	}
	var sel [4]byte
	copy(sel[:], input[:4])
	return &Decoded{
		Contract: addr,
		Name:     m.Name,
		Sig:      m.Sig,
		Selector: sel,
		Values:   values,
	}
}
func (fx *Fx) event(addr common.Address, topics []common.Hash, data []byte) *Decoded {
	if len(topics) == 0 {
		return nil
	}
	c, ok := fx.Contracts[addr]
	if !ok || c.ABI == nil {
		return nil
	}
	e, err := c.ABI.EventByID(topics[0])
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

	var sel [4]byte
	copy(sel[:], topics[0][:4])
	return &Decoded{
		Contract: addr,
		Name:     e.Name,
		Sig:      e.Sig,
		Selector: sel,
		Values:   values,
	}
}

func (fx *Fx) revert(addr common.Address, data []byte) *Decoded {
	if len(data) < 4 {
		return nil
	}
	c, ok := fx.Contracts[addr]
	if !ok || c.ABI == nil {
		return nil
	}
	var sel [4]byte
	copy(sel[:], data[:4])
	e, err := c.ABI.ErrorByID(sel)
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
		Selector: sel,
		Values:   values,
	}
}

func valuesEqual(a, b any) bool {
	switch av := a.(type) {
	case common.Address:
		bv, ok := b.(common.Address)
		return ok && av == bv
	case *big.Int:
		bv, ok := b.(*big.Int)
		return ok && bv != nil && av.Cmp(bv) == 0
	case common.Hash:
		bv, ok := b.(common.Hash)
		return ok && av == bv
	case bool:
		bv, ok := b.(bool)
		return ok && av == bv
	case string:
		bv, ok := b.(string)
		return ok && av == bv
	default:
		return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
	}
}
