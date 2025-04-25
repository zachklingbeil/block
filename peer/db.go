package peer

import (
	"encoding/json"
	"fmt"
)

func (p *Peers) LoadPeers() error {
	peerJSONs, err := p.Factory.Redis.SMembers(p.Factory.Ctx, "peers").Result()
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
		p.Slice = append(p.Slice, peer)
	}
	fmt.Printf("%d peers\n", len(p.Slice))
	p.SavePeers()
	return nil
}

func (p *Peers) SavePeers() error {
	for _, peer := range p.Slice {
		peerJSON, err := json.Marshal(peer)
		if err != nil {
			return fmt.Errorf("failed to serialize peer: %w", err)
		}
		err = p.Factory.Redis.SAdd(p.Factory.Ctx, "peers", peerJSON).Err()
		if err != nil {
			return fmt.Errorf("failed to save peer to Redis: %w", err)
		}
	}
	fmt.Printf("%d peers stored in Redis\n", len(p.Slice))
	return nil
}
