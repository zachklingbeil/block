package value

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/zachklingbeil/factory"
)

type Value struct {
	Factory  *factory.Factory
	Peers    []*Peer
	Tokens   []*Token
	Universe map[common.Address]any
	Maps     *Maps
}

type Maps struct {
	LoopringId map[int64]*Peer
	TokenId    map[int64]*Token
}

func NewValue(factory *factory.Factory) *Value {
	v := &Value{
		Factory: factory,

		Universe: make(map[common.Address]any),
	}
	v.LoadTokens()
	v.LoadPeers()
	v.populateMap()
	fmt.Printf("Universe: %d\n", len(v.Universe))
	return v
}

func (v *Value) populateMap() {
	for _, p := range v.Peers {
		v.Universe[common.HexToAddress(strings.ToLower(p.Address))] = p
	}

	for _, t := range v.Tokens {
		v.Universe[t.Address] = t
	}
}

func (v *Value) Source(address common.Address) any {
	return v.Universe[address]
}
