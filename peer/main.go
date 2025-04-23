package peer

import (
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/zachklingbeil/factory"
)

type Peers struct {
	Factory        *factory.Factory
	Map            map[any]*Peer
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

func (p *Peers) LoadPeers() error {
	peerJSONs, err := p.Factory.Db.Rdb.SMembers(p.Factory.Ctx, "peers").Result()
	if err != nil {
		return fmt.Errorf("failed to retrieve peers from Redis: %w", err)
	}

	p.Factory.Mu.Lock()
	defer p.Factory.Mu.Unlock()

	for _, peerJSON := range peerJSONs {
		var peer Peer
		if err := json.Unmarshal([]byte(peerJSON), &peer); err != nil {
			return fmt.Errorf("failed to deserialize peer JSON: %w", err)
		}
		p.Map[peer.Address] = &peer
		p.Map[peer.ENS] = &peer
		p.Map[peer.LoopringENS] = &peer
		p.Map[peer.LoopringID] = &peer
	}

	fmt.Printf("%d peers\n", len(p.Map))
	return nil
}

func (p *Peers) SavePeers() error {
	for address, peer := range p.Map {
		peerJSON, err := json.Marshal(peer)
		if err != nil {
			return fmt.Errorf("failed to serialize peer (address: %s): %w", address, err)
		}

		err = p.Factory.Db.Rdb.SAdd(p.Factory.Ctx, "peers", peerJSON).Err()

		if err != nil {
			return fmt.Errorf("failed to store peer in Redis (address: %s): %w", address, err)
		}
	}
	fmt.Printf("%d peers\n", len(p.Map))
	return nil
}
