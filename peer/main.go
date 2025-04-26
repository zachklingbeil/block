package peer

import (
	"encoding/json"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zachklingbeil/factory"
)

type Peers struct {
	Factory *factory.Factory
	ApiKey  string
	Slice   []Peer
	Map     map[string]*Peer
}

type Peer struct {
	ENS         string `json:"ens"`
	LoopringENS string `json:"loopringEns"`
	LoopringID  string `json:"loopringId"`
	Address     string `json:"address"`
}

func NewPeers(factory *factory.Factory) *Peers {
	p := &Peers{
		Factory: factory,
		Slice:   []Peer{},
		Map:     make(map[string]*Peer),
	}

	if err := p.LoadPeers(); err != nil {
		log.Fatalf("Failed to load peers from Redis: %v", err)
	}
	p.BuildPeerMap()
	return p
}

func (p *Peers) GetPeer(value string) *Peer {
	p.Factory.Rw.RLock()
	peer, exists := p.Map[value]
	p.Factory.Rw.RUnlock()
	if exists {
		return peer
	}
	return p.CreatePeer(value)
}

func (p *Peers) CreatePeer(value string) *Peer {
	p.Factory.Rw.Lock()
	defer p.Factory.Rw.Unlock()

	new := &Peer{}

	switch {
	case common.IsHexAddress(value):
		new.Address = value
	case len(value) > 12 && value[len(value)-13:] == ".loopring.eth":
		new.LoopringENS = value
	case len(value) > 4 && value[len(value)-4:] == ".eth":
		new.ENS = value
	default:
		new.LoopringID = value
	}
	p.Map[value] = new
	p.Slice = append(p.Slice, *new)
	return new
}

func (p *Peers) Process(address string) {
	peer := p.GetPeer(address)
	p.GetENS(peer)
	p.GetLoopringID(peer)
	p.GetLoopringENS(peer)
}

func (p *Peers) LoadPeers() error {
	source, err := p.Factory.Redis.SMembers(p.Factory.Ctx, "peers").Result()
	if err != nil {
		return err
	}

	for _, peerJSON := range source {
		var peer Peer
		if err := json.Unmarshal([]byte(peerJSON), &peer); err != nil {
			log.Printf("Skipping invalid peer: %v", err)
			continue
		}
		p.Slice = append(p.Slice, peer)
	}
	return nil
}

func (p *Peers) BuildPeerMap() {
	for i := range p.Slice {
		peer := &p.Slice[i]
		p.Map[peer.Address] = peer
		p.Map[peer.ENS] = peer
		p.Map[peer.LoopringID] = peer
		p.Map[peer.LoopringENS] = peer
	}
}
