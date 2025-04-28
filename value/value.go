package value

import (
	"encoding/json"
	"log"

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
	return v
}

func (v *Value) LoadPeers() error {
	source, err := v.Factory.Data.RB.SMembers(v.Factory.Ctx, "peers").Result()
	if err != nil {
		return err
	}

	for _, peerJSON := range source {
		var peer Peer
		if err := json.Unmarshal([]byte(peerJSON), &peer); err != nil {
			log.Printf("Skipping invalid peer: %v", err)
			continue
		}
		v.Peers = append(v.Peers, peer)
	}
	v.Factory.State.Add("peers", len(v.Peers))
	return nil
}
