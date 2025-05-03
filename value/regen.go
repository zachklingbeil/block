package value

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
)

//go:embed lp.json
var peers []byte

func (v *Value) LoadEmbeddedPeers() ([]Peer, error) {
	// Parse the embedded lp.json data into []Peer
	var embeddedPeers []Peer
	if err := json.Unmarshal(peers, &embeddedPeers); err != nil {
		return nil, fmt.Errorf("failed to parse embedded lp.json: %v", err)
	}

	return embeddedPeers, nil
}
func (v *Value) UpdatePeersFromEmbeddedData() error {
	// Load embedded peers
	embeddedPeers, err := v.LoadEmbeddedPeers()
	if err != nil {
		return err
	}

	// Lock the factory for writing
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	// Iterate through the embedded data
	for _, embeddedPeer := range embeddedPeers {
		// Format the address using v.Format
		formattedAddress := v.Format(embeddedPeer.Address)

		// Check if the peer exists in the current map
		if peer, exists := v.Map[formattedAddress]; exists {
			// Update the peer's LoopringENS and LoopringID
			peer.LoopringENS = embeddedPeer.LoopringENS
			peer.LoopringID = embeddedPeer.LoopringID

			// Save the updated peer to Redis
			if err := v.Save(peer); err != nil {
				log.Printf("Failed to save updated peer to Redis: %v", err)
				continue
			}

			// Print the entire updated peer
			fmt.Printf("%+v\n", peer)
		} else {
			fmt.Printf("Peer with address %s not found in the map\n", formattedAddress)
		}
	}

	return nil
}
func (v *Value) LoadPeers() error {
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	hashKey := "peer"
	source, err := v.Factory.Data.RB.HGetAll(v.Factory.Ctx, hashKey).Result()
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
		v.Map[peer.Address] = &peer
		v.Map[peer.ENS] = &peer
		v.Map[peer.LoopringENS] = &peer
		v.Map[peer.LoopringID] = &peer
	}
	v.Peers = peers
	fmt.Printf("%d peers\n", len(v.Peers))
	return nil
}

func (v *Value) Save(peer *Peer) error {
	peerJSON, err := json.Marshal(peer)
	if err != nil {
		return fmt.Errorf("failed to serialize peer: %v", err)
	}
	hashKey := "peer"
	if err := v.Factory.Data.RB.HSet(v.Factory.Ctx, hashKey, peer.Address, peerJSON).Err(); err != nil {
		return fmt.Errorf("failed to store peer in Redis hash: %v", err)
	}
	return nil
}
