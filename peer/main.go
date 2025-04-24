package peer

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/zachklingbeil/factory"
)

type Peers struct {
	Factory        *factory.Factory
	Map            map[any]*Peer
	Slice          []Peer
	LoopringApiKey string
}

type Peer struct {
	Address     string `json:"address"`
	ENS         string `json:"ens"`
	LoopringENS string `json:"loopringEns"`
	LoopringID  int64  `json:"loopringId"`
}

func HelloPeers(factory *factory.Factory) *Peers {
	peers := &Peers{
		Factory:        factory,
		Map:            make(map[any]*Peer),
		LoopringApiKey: os.Getenv("LOOPRING_API_KEY"),
	}

	if err := peers.LoadPeers(); err != nil {
		fmt.Printf("Error loading peers: %v\n", err)
	}
	return peers
}

func (p *Peers) HelloUniverse(peers []any) {
	for _, peers := range peers {
		switch v := peers.(type) {
		case string:
			p.Factory.Mu.Lock()
			if _, exists := p.Map[v]; !exists {
				p.Map[v] = &Peer{Address: v, ENS: v, LoopringENS: v}
			}
			p.Factory.Mu.Unlock()
		case int64:
			p.Factory.Mu.Lock()
			if _, exists := p.Map[v]; exists {
				p.Map[v] = &Peer{LoopringID: v}
			}
		}
		p.Factory.Mu.Unlock()
	}

	for {
		p.Factory.Mu.Lock()
		for _, peer := range p.Map {
			p.processPeer(peer)
			fmt.Printf("%s %s %s %d\n", peer.Address, peer.ENS, peer.LoopringENS, peer.LoopringID)
		}
		p.Factory.Mu.Unlock()
	}
}

func (p *Peers) processPeer(peer *Peer) {
	p.GetENS(peer, peer.Address)
	p.GetLoopringENS(peer, peer.Address)
	p.GetLoopringID(peer, peer.Address)
}
