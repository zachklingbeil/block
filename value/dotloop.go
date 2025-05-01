package value

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
)

//go:embed dotloop.json
var dotloop []byte

type DotLoop struct {
	LoopringENS string `json:"loopringEns"`
	Address     string `json:"address"`
}
type Peer struct {
	ENS         string `json:"ens"`
	LoopringENS string `json:"loopringEns"`
	LoopringID  string `json:"loopringId"`
	Address     string `json:"address"`
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
	}
	v.Peers = peers
	log.Printf("Loaded %d peers from Redis hash: %s", len(v.Peers), hashKey)
	return nil
}

// ...existing code...

// UpdatePeersLoopringENS updates the LoopringENS field of existing peers using dotloop.json data.
func (v *Value) UpdatePeersLoopringENS() error {
	// Unmarshal dotloop.json into []DotLoop
	var dotloopData []DotLoop
	if err := json.Unmarshal(dotloop, &dotloopData); err != nil {
		return fmt.Errorf("failed to unmarshal dotloop: %v", err)
	}

	// Build a map for quick lookup by address
	dotloopMap := make(map[string]string)
	for _, d := range dotloopData {
		dotloopMap[d.Address] = d.LoopringENS
	}

	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	updated := 0
	for _, peer := range v.Peers {
		if newENS, ok := dotloopMap[peer.Address]; ok {
			peer.LoopringENS = newENS
			updated++
		}
		v.Save(peer)
	}
	log.Printf("Updated LoopringENS for %d peers from dotloop.json", updated)
	return nil
}

// ...existing code...
