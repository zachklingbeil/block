package value

import (
	"github.com/zachklingbeil/factory"
)

type Value struct {
	Factory *factory.Factory
	Peers   []Peer
	Tokens  []Token
	Map     map[string]*Peer
}

func NewValue(factory *factory.Factory) *Value {
	v := &Value{
		Factory: factory,
		Map:     make(map[string]*Peer),
		Peers:   make([]Peer, 0),
	}
	v.LoadPeers()
	v.LoadTokens()
	v.Factory.State.AddToPackage("value", "peers", len(v.Peers))
	v.Factory.State.AddToPackage("value", "tokens", len(v.Tokens))
	return v
}
