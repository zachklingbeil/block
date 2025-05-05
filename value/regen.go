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
