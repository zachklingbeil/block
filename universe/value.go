package universe

import (
	"fmt"

	"github.com/zachklingbeil/factory"
)

type One struct {
	Factory *factory.Factory
	Peers   []*Peer
	Tokens  []*Token
	Map     map[string]any
	Maps    *Maps
}

type Maps struct {
	LoopringId map[int64]string
	TokenId    map[int64]string
}

func NewOne(factory *factory.Factory) *One {
	v := &One{
		Factory: factory,
		Map:     make(map[string]any),
		Maps: &Maps{
			LoopringId: make(map[int64]string),
			TokenId:    make(map[int64]string),
		},
	}

	v.LoadTokens()
	v.LoadPeers()

	for _, p := range v.Peers {
		v.Map[p.Address] = p
	}

	for _, t := range v.Tokens {
		v.Map[t.Address] = t
	}
	fmt.Printf("Map: %d\n", len(v.Map))
	return v
}

func (v *One) Source(address string) any {
	return v.Map[address]
}
