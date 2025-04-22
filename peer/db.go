package peer

import (
	"encoding/json"
	"fmt"
)

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
	fmt.Printf("%d peers stored in Redis\n", len(p.Map))
	return nil
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
	}

	fmt.Printf("%d peers loaded from Redis\n", len(p.Map))
	return nil
}
