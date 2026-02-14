package fx

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

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
		Values:   cleanValues(values),
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
		Values:   cleanValues(values),
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
		Values:   cleanValues(values),
	}
}

func cleanValues(values map[string]any) map[string]any {
	for k, v := range values {
		values[k] = cleanValue(v)
	}
	return values
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
		// Attempt to decode as ABI-encoded tuple/value
		return val

	case common.Hash:
		// Indexed address params come back as Hash (left-padded)
		var zero [12]byte
		if [12]byte(val[:12]) == zero {
			return common.BytesToAddress(val[12:])
		}
		return val

	case []any:
		// Tuple/struct arrays from ABI decoding
		out := make([]any, len(val))
		for i, elem := range val {
			out[i] = cleanValue(elem)
		}
		return out

	default:
		return v
	}
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
