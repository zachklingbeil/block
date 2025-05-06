package value

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zachklingbeil/block/value/peer"
	"github.com/zachklingbeil/block/value/token"
	"github.com/zachklingbeil/factory"
)

type Value struct {
	Factory *factory.Factory
	Peer    *peer.Peers
	Token   *token.Tokens
	Map     map[common.Address]any
}

func NewValue(factory *factory.Factory) *Value {
	v := &Value{
		Factory: factory,
		Peer:    peer.NewPeers(factory),
		Token:   token.NewTokens(factory),
		Map:     make(map[common.Address]any),
	}
	v.populateMap()
	fmt.Printf("Map: %d\n", len(v.Map))
	return v
}

func (v *Value) populateMap() {
	for _, p := range v.Peer.Peers {
		v.Map[common.HexToAddress(strings.ToLower(p.Address))] = p
	}

	for _, t := range v.Token.Tokens {
		v.Map[t.Address] = t
	}
}

func (v *Value) Source(address common.Address) any {
	return v.Map[address]
}
