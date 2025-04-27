package value

import (
	"log"

	"github.com/zachklingbeil/factory"
)

type Value struct {
	Factory *factory.Factory
	Map     map[string]*Peer
	Peers   []Peer
	Tokens  []Token
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
