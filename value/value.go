package value

import (
	"github.com/zachklingbeil/factory"
)

type Value struct {
	Factory  *factory.Factory
	Peers    []Peer
	Tokens   []Token
	Map      map[string]*Peer
	TokenMap map[any]*Token
}

func NewValue(factory *factory.Factory) *Value {
	v := &Value{
		Factory:  factory,
		Map:      make(map[string]*Peer),
		TokenMap: make(map[any]*Token),
	}
	v.LoadTokens()
	// v.LoadPeers()
	// v.ReprocessPeers()
	return v
}
