package value

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

//go:embed loopring.json
var dotloop []byte

func (v *Value) Refresh() {
	for i := range v.Peers {
		fmt.Printf("%d\n", i)
		peer := v.Peers[i]
		v.Format(peer.Address)
		v.Save(peer)
	}
}

// get loopring ids
// func (v *Value) Refresh() {
// 	for i := range v.Peers {
// 		fmt.Printf("%d\n", i)
// 		peer := v.Peers[i]
// 		if peer.LoopringID == "." || peer.LoopringID == "!" || peer.LoopringID == "" {
// 			v.GetLoopringID(peer)
// 		}
// 	}
// }

type LoopId struct {
	LoopringId string `json:"loopringId"`
	Address    string `json:"address"`
	FirstBlock string `json:"firstBlock"`
}

func (v *Value) UpdatePeersLoopringIDAndFirstBlock() error {
	var loopringData []LoopId
	if err := json.Unmarshal(dotloop, &loopringData); err != nil {
		return fmt.Errorf("failed to unmarshal loopring.json: %v", err)
	}

	// Build a map for quick lookup by address
	loopringMap := make(map[string]LoopId)
	for _, d := range loopringData {
		loopringMap[strings.ToLower(d.Address)] = d // ensure address is lowercased for matching
	}

	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	updated := 0
	existing := make(map[string]struct{})
	for _, peer := range v.Peers {
		addr := strings.ToLower(peer.Address)
		if info, ok := loopringMap[addr]; ok {
			peer.LoopringID = info.LoopringId
			peer.FirstBlock = info.FirstBlock
			updated++
		}
		existing[addr] = struct{}{}
		v.Save(peer)
	}

	// Add missing peers from loopringMap
	added := 0
	for addr, info := range loopringMap {
		if _, ok := existing[addr]; !ok {
			newPeer := &Peer{
				Address:    info.Address,
				LoopringID: info.LoopringId,
				FirstBlock: info.FirstBlock,
			}
			v.Peers = append(v.Peers, newPeer)
			v.Save(newPeer)
			added++
		}
	}

	log.Printf("Updated LoopringId and FirstBlock for %d peers, added %d new peers from loopring.json", updated, added)
	return nil
}
