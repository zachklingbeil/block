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

type Peer struct {
	ENS         string `json:"ens"`
	LoopringENS string `json:"loopringEns"`
	LoopringID  string `json:"loopringId"`
	Address     string `json:"address"`
}

type Token struct {
	Token    string `json:"token,omitempty"`
	Address  string `json:"address,omitempty"`
	Decimals string `json:"decimals,omitempty"`
	TokenId  string `json:"tokenId,omitempty"`
	TokenInt int64  `json:"tokenInt,omitempty"`
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

func (v *Value) LoadTokens() error {
	source, err := v.Factory.Data.RB.SMembers(v.Factory.Ctx, "token").Result()
	if err != nil {
		log.Fatalf("Failed to fetch tokens from Redis: %v", err)
	}

	tokens := make([]Token, 0, len(source))
	for _, tokenJSON := range source {
		var token Token
		if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
			log.Printf("Skipping invalid token: %v (data: %s)", err, tokenJSON)
			continue
		}
		tokens = append(tokens, token)
	}
	v.Tokens = tokens
	v.Factory.State.Add("tokens", len(v.Tokens))
	return nil
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
