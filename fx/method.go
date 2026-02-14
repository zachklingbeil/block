package fx

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func (fx *Fx) decode(addr common.Address, topics []common.Hash, data []byte) *Event {
	c, ok := fx.Contracts[addr]
	if !ok || c.ABI == nil {
		return nil
	}
	if len(topics) > 0 {
		return fx.decodeEvent(c, addr, topics, data)
	}
	return fx.decodeMethod(c, addr, data)
}

func (fx *Fx) decodeMethod(c *Contract, addr common.Address, data []byte) *Event {
	if len(data) < 4 {
		return nil
	}
	m, err := c.ABI.MethodById(data[:4])
	if err != nil {
		return nil
	}
	values := make(map[string]any)
	if len(data) > 4 {
		if err := m.Inputs.UnpackIntoMap(values, data[4:]); err != nil {
			return nil
		}
	}
	var sel [4]byte
	copy(sel[:], data[:4])
	cleanValues(values)
	return &Event{
		Contract: addr,
		Name:     m.Name,
		Sig:      m.Sig,
		Selector: sel,
		Params:   values,
	}
}

func (fx *Fx) decodeEvent(c *Contract, addr common.Address, topics []common.Hash, data []byte) *Event {
	e, err := c.ABI.EventByID(topics[0])
	if err != nil {
		return nil
	}
	values := make(map[string]any)
	indexed := make(abi.Arguments, 0, len(e.Inputs))
	nonIndexed := make(abi.Arguments, 0, len(e.Inputs))
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
	cleanValues(values)
	return &Event{
		Contract: addr,
		Name:     e.Name,
		Sig:      e.Sig,
		Selector: sel,
		Action:   Classify(topics[0], topics, values),
		Params:   values,
	}
}

func cleanValues(values map[string]any) {
	for k, v := range values {
		values[k] = cleanValue(v)
	}
}

func cleanValue(v any) any {
	switch val := v.(type) {
	case [32]byte:
		var zero [12]byte
		if [12]byte(val[:12]) == zero {
			return common.BytesToAddress(val[12:])
		}
		return common.BytesToHash(val[:])
	case []byte:
		if len(val) == 20 {
			return common.BytesToAddress(val)
		}
		if len(val) == 32 {
			return cleanValue([32]byte(val))
		}
		return val
	case common.Hash:
		var zero [12]byte
		if [12]byte(val[:12]) == zero {
			return common.BytesToAddress(val[12:])
		}
		return val
	case []any:
		out := make([]any, len(val))
		for i, elem := range val {
			out[i] = cleanValue(elem)
		}
		return out
	default:
		return v
	}
}

func matchOutcome(a *Outcome, events []Event) bool {
	if len(a.Events) != len(events) {
		return false
	}
	for i := range a.Events {
		if a.Events[i].Selector != events[i].Selector {
			return false
		}
		if a.Events[i].Contract != events[i].Contract {
			return false
		}
		if a.Events[i].Action != events[i].Action {
			return false
		}
	}
	return true
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
