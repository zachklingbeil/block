package peer

import (
	"github.com/zachklingbeil/factory"
)

type Peers struct {
	Factory *factory.Factory
	Peers   []*Peer `json:"peers,omitempty"`
}

func NewPeers(factory *factory.Factory) *Peers {
	return &Peers{
		Factory: factory,
		Peers:   make([]*Peer, 0),
	}
}

type Peer struct {
	ENS         string `json:"ens,omitempty"`
	LoopringENS string `json:"loopringEns,omitempty"`
	LoopringID  string `json:"loopringId,omitempty"`
	Address     string `json:"address,omitempty"`
	FirstBlock  string `json:"firstBlock,omitempty"`
}
