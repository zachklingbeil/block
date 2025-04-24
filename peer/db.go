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

	// Use factory.Math.Down to count down as peers are added
	go p.Factory.Math.Down(int64(len(peerJSONs)), func(count int64) {
		index := len(peerJSONs) - int(count) // Calculate the index from the count
		if index < 0 || index >= len(peerJSONs) {
			return // Skip invalid indices
		}

		peerJSON := peerJSONs[index]
		var peer Peer
		if err := json.Unmarshal([]byte(peerJSON), &peer); err != nil {
			fmt.Printf("Failed to deserialize peer JSON: %v\n", err)
			return
		}
		p.Circuit.AddString(peer.Address, peer)
		p.Circuit.AddString(peer.ENS, peer)
		p.Circuit.AddString(peer.LoopringENS, peer)
		p.Circuit.AddInt(peer.LoopringID, peer)

		fmt.Printf("Processed peer %d %s\n", count, peer.Address)
	})

	fmt.Printf("%d peers\n", len(peerJSONs))
	return nil
}

func (p *Peers) SavePeers() error {
	for address, peer := range p.Map {
		peerJSON, err := json.Marshal(peer)
		if err != nil {
			return fmt.Errorf("failed to serialize peer (address: %s): %w", address, err)
		}

		err = p.Factory.Redis.SAdd(p.Factory.Ctx, "peers", peerJSON).Err()

		if err != nil {
			return fmt.Errorf("failed to store peer in Redis (address: %s): %w", address, err)
		}
	}
	fmt.Printf("%d peers\n", len(p.Map))
	return nil
}

func (p *Peers) LoadPeer() error {
	query := `
        SELECT address, ens, loopringEns, loopringId FROM peers
    `
	rows, err := p.Factory.Pg.Query(query)
	if err != nil {
		return fmt.Errorf("failed to load peers from database: %w", err)
	}
	defer rows.Close()
	p.Factory.Mu.Lock()
	defer p.Factory.Mu.Unlock()
	for rows.Next() {
		var peer Peer
		if err := rows.Scan(&peer.Address, &peer.ENS, &peer.LoopringENS, &peer.LoopringID); err != nil {
			return fmt.Errorf("failed to scan peer row: %w", err)
		}
		p.Circuit.AddString(peer.Address, peer)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over peer rows: %w", err)
	}
	// p.SavePeers()
	return nil
}
