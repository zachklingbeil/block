package value

import (
	"github.com/zachklingbeil/factory"
)

type Value struct {
	Factory  *factory.Factory
	Peers    []Peer
	Tokens   []Token
	Blocks   []Block
	Map      map[string]*Peer
	TokenMap map[int64]*Token
}

func NewValue(factory *factory.Factory) *Value {
	v := &Value{
		Factory:  factory,
		Map:      make(map[string]*Peer),
		TokenMap: make(map[int64]*Token),
	}
	v.LoadPeers()
	v.LoadTokens()
	// v.LoadBlocks()
	// v.HandleNewPeers()
	v.Factory.State.Add("value", "peers", len(v.Peers))
	v.Factory.State.Add("value", "tokens", len(v.Tokens))
	// v.Factory.State.Add("value", "blocks", len(v.Blocks))
	return v
}
