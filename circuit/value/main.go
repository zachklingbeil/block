package value

import (
	"log"

	"github.com/zachklingbeil/factory"
)

type Value struct {
	Factory *factory.Factory
	Peers   []Peer
	Tokens  []Token
	Map     map[string]*Peer
	State   map[string]any
}

func NewValue(factory *factory.Factory) *Value {
	v := &Value{
		Factory: factory,
		Map:     make(map[string]*Peer),
		Peers:   make([]Peer, 0),
		Tokens:  make([]Token, 0),
	}
	if err := v.LoadPeers(); err != nil {
		log.Fatalf("Failed to load peers from Redis: %v", err)
	}
	return v
}
