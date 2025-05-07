package value

import (
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Value struct {
	Factory  *factory.Factory
	Peers    []*Peer
	Tokens   []*Token
	Universe map[string]any
	Maps     *Maps
}

type Maps struct {
	LoopringId map[int64]string
	TokenId    map[int64]string
}

func NewValue(factory *factory.Factory) *Value {
	v := &Value{
		Factory:  factory,
		Universe: make(map[string]any),
		Maps: &Maps{
			LoopringId: make(map[int64]string),
			TokenId:    make(map[int64]string),
		},
	}

	v.LoadTokens()
	v.LoadPeers()

	for _, p := range v.Peers {
		v.Universe[p.Address] = p
	}

	for _, t := range v.Tokens {
		v.Universe[t.Address] = t
	}
	fmt.Printf("Universe: %d\n", len(v.Universe))
	return v
}

func (v *Value) Source(address string) any {
	return v.Universe[address]
}
