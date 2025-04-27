package circuit

import (
	"encoding/json"
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Circuit struct {
	Factory *factory.Factory
	Map     map[string]any
	Value   []Value
	Peers   []Peer  `json:"peers,omitempty"`
	Tokens  []Token `json:"tokens,omitempty"`
	PeerMap map[string]*Peer
	IDMap   map[int64]*Token
	LPMap   map[int64]*Token
	ApiKey  string
}

type Value struct {
	Peer       *Peer       `json:"peer,omitempty"`
	Token      *Token      `json:"token,omitempty"`
	Block      *Block      `json:"block,omitempty"`
	Coordinate *Coordinate `json:"coordinate,omitempty"`
	Tx         *Tx         `json:"tx,omitempty"`
}

func NewCircuit(factory *factory.Factory) *Circuit {
	circuit := &Circuit{
		Factory: factory,
		Map:     make(map[string]any),
		Tokens:  make([]Token, 270),
		Peers:   make([]Peer, 0),
		IDMap:   make(map[int64]*Token),
		LPMap:   make(map[int64]*Token),
		PeerMap: make(map[string]*Peer),
	}
	circuit.Load()
	// fmt.Printf("%d tokens\n", len(circuit.Tokens))

	return circuit
}

func (c *Circuit) Continue() error {
	c.Factory.Mu.Lock()
	defer c.Factory.Mu.Unlock()
	source, err := c.Factory.Data.RB.SMembers(c.Factory.Ctx, "value").Result()
	if err != nil {
		return fmt.Errorf("failed to load values from Redis: %w", err)
	}

	for _, i := range source {
		var value Value
		if err := json.Unmarshal([]byte(i), &value); err != nil {
			return fmt.Errorf("failed to unmarshal value: %w", err)
		}
		// c.Map[value.Address] = &value
		// c.Map[value.ENS] = &value
		// c.Map[value.LoopringENS] = &value
		// c.Map[value.LoopringID] = &value
		// c.Map[value.Symbol] = &value
		// c.Map[value.Address] = &value
		// c.Map[value.LoopringID] = &value
		// c.Map[value.Token] = &value
		// c.Values = append(c.Values, value)
	}
	fmt.Printf("%d\n", len(c.Map))
	return nil
}
