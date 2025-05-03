package peer

import (
	"encoding/json"
	"fmt"
	"log"
)

func (p *Peers) LoadPeers() error {
	p.Factory.Rw.Lock()
	defer p.Factory.Rw.Unlock()

	hashKey := "peer"
	source, err := p.Factory.Data.RB.HGetAll(p.Factory.Ctx, hashKey).Result()
	if err != nil {
		return fmt.Errorf("failed to fetch peers from Redis hash: %v", err)
	}
	peers := make([]*Peer, 0, len(source))
	for _, peerJSON := range source {
		var peer Peer
		if err := json.Unmarshal([]byte(peerJSON), &peer); err != nil {
			log.Printf("Skipping invalid peer: %v (data: %s)", err, peerJSON)
			continue
		}
		peers = append(peers, &peer)
	}
	p.Peers = peers
	fmt.Printf("%d peers\n", len(p.Peers))
	return nil
}

func (p *Peers) Save(peer *Peer) error {
	peerJSON, err := json.Marshal(peer)
	if err != nil {
		return fmt.Errorf("failed to serialize peer: %v", err)
	}
	hashKey := "peer"
	if err := p.Factory.Data.RB.HSet(p.Factory.Ctx, hashKey, peer.Address, peerJSON).Err(); err != nil {
		return fmt.Errorf("failed to store peer in Redis hash: %v", err)
	}
	return nil
}
