package value

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
)

const (
	byAddress = "https://api3.loopring.io/api/v3/account?owner=%s"
	byId      = "https://api3.loopring.io/api/v3/account?accountId=%s"
	dotLoop   = "https://api3.loopring.io/api/wallet/v3/resolveName?owner=%s"
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
		log.Fatalf("Failed to fetch peers from Redis: %v", err)
	}
	v.Peers = make([]Peer, 0, len(source))
	for _, peerJSON := range source {
		var peer Peer
		if err := json.Unmarshal([]byte(peerJSON), &peer); err != nil {
			log.Printf("Skipping invalid peer: %v", err)
			continue
		}
		v.Peers = append(v.Peers, peer)
		v.Map[peer.Address] = &peer
		v.Map[peer.LoopringENS] = &peer
		v.Map[peer.LoopringID] = &peer
		v.Map[peer.ENS] = &peer
	}
	return nil
}

func (v *Value) HelloUniverse(value string) {
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	peer, exists := v.Map[value]
	if !exists {
		peer = &Peer{
			ENS:         "",
			LoopringENS: "",
		}
		if common.IsHexAddress(value) {
			peer.Address = v.Format(value)
		} else {
			peer.LoopringID = value
			v.GetLoopringAddress(peer)
		}
		v.Peers = append(v.Peers, *peer)
	}

	if peer.ENS == "" && peer.ENS != "!" {
		v.GetENS(peer)
	}
	if peer.LoopringENS == "" && peer.LoopringENS != "!" {
		v.GetLoopringENS(peer)
	}
	if peer.LoopringID == "" && peer.LoopringID != "!" {
		v.GetLoopringID(peer)
	}
	if peer.Address == "" && peer.Address != "!" {
		v.GetLoopringAddress(peer)
	}
	fmt.Printf("%s %s %s\n", peer.ENS, peer.LoopringENS, peer.LoopringID)
}

func (v *Value) Save(peer *Peer) {
	peerJSON, err := json.Marshal(peer)
	if err != nil {
		log.Printf("Failed to serialize peer: %v", err)
		return
	}
	err = v.Factory.Data.RB.SAdd(v.Factory.Ctx, "peer", peerJSON).Err()
	if err != nil {
		log.Printf("Failed to store peer in Redis: %v", err)
	}
}
