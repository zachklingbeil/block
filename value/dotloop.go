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
