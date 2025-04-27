package value

import (
	"encoding/json"
	"log"

	"github.com/ethereum/go-ethereum/common"
)

type Peer struct {
	ENS         string `json:"ens"`
	LoopringENS string `json:"loopringEns"`
	LoopringID  string `json:"loopringId"`
	Address     string `json:"address"`
}

func (v *Value) LoadPeers() error {
	source, err := v.Factory.Data.RB.SMembers(v.Factory.Ctx, "peer").Result()
	if err != nil {
		return err
	}

	for _, peerJSON := range source {
		var peer Peer
		if err := json.Unmarshal([]byte(peerJSON), &peer); err != nil {
			log.Printf("Skipping invalid peer: %v", err)
			continue
		}
		v.Peers = append(v.Peers, peer)
	}
	return nil
}

func (v *Value) HelloUniverse(address string) {
	peer := v.GetPeer(address)
	v.GetENS(peer)
	v.GetLoopringID(peer)
	v.GetLoopringENS(peer)
}

func (v *Value) GetPeer(value string) *Peer {
	v.Factory.Rw.RLock()
	peer, exists := v.Map[value]
	v.Factory.Rw.RUnlock()
	if exists {
		return peer
	}
	return v.CreatePeer(value)
}

func (v *Value) CreatePeer(value string) *Peer {
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	new := &Peer{}

	switch {
	case common.IsHexAddress(value):
		new.Address = value
	case len(value) > 12 && value[len(value)-13:] == ".loopring.eth":
		new.LoopringENS = value
	case len(value) > 4 && value[len(value)-4:] == ".eth":
		new.ENS = value
	default:
		new.LoopringID = value
	}
	v.Map[value] = new
	v.Peers = append(v.Peers, *new)
	return new
}
