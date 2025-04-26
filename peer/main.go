package peer

import (
	"encoding/json"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wealdtech/go-ens/v3"
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
	source, err := p.Factory.Redis.SMembers(p.Factory.Ctx, "peers").Result()
	if err != nil {
		log.Fatalf("Failed to fetch peers from Redis: %v", err)
	}
	for _, peers := range source {
		var peer Peer
		if err := json.Unmarshal([]byte(peers), &peer); err != nil {
			log.Printf("Skipping invalid peer: %v", err)
			continue
		}
		p.Slice = append(p.Slice, peer)
		p.Map[peer.Address] = &peer
		p.Map[peer.ENS] = &peer
		p.Map[peer.LoopringENS] = &peer
		p.Map[peer.LoopringID] = &peer
	}
	return p
}

func (p *Peers) Process(address string) {
	peer, exists := p.GetPeer(address)
	if !exists {
		return // Exit if the peer doesn't exist
	}

	if peer.ENS != "" && peer.ENS != "." {
		return // Skip processing if ENS is already set or marked as checked
	}

	ensName, err := ens.ReverseResolve(p.Factory.Eth, common.HexToAddress(address))

	p.Factory.Rw.Lock()
	defer p.Factory.Rw.Unlock()

	if err != nil || ensName == "" {
		peer.ENS = "." // Mark no ENS
	} else {
		peer.ENS = p.Format(ensName)
	}
}

// hex -> ENS [.eth] or "."
func (p *Peers) ProcessENS(address string) {
	peer, exists := p.GetPeer(address)
	if !exists {
		return // Exit if the peer doesn't exist
	}

	if peer.ENS != "" && peer.ENS != "." {
		return // Skip processing if ENS is already set or marked as checked
	}

	ensName, err := ens.ReverseResolve(p.Factory.Eth, common.HexToAddress(address))

	p.Factory.Rw.Lock()
	defer p.Factory.Rw.Unlock()

	if err != nil || ensName == "" {
		peer.ENS = "." // Mark no ENS
	} else {
		peer.ENS = p.Format(ensName)
	}
}

// GetPeer safely retrieves a Peer from the map using the given key.
func (p *Peers) GetPeer(key string) (*Peer, bool) {
	p.Factory.Rw.RLock()         // Acquire read lock
	defer p.Factory.Rw.RUnlock() // Release read lock when done

	peer, exists := p.Map[key]
	return peer, exists
}

// SetPeer safely sets a Peer in the map with the given key.
func (p *Peers) SetPeer(key string, peer *Peer) {
	p.Factory.Rw.Lock()
	defer p.Factory.Rw.Unlock()

	p.Map[key] = peer
}
