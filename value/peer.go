package value

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

const (
	byAddress = "https://api3.loopring.io/api/v3/account?owner=%s"
	byId      = "https://api3.loopring.io/api/v3/account?accountId=%s"
	dotLoop   = "https://api3.loopring.io/api/wallet/v3/resolveName?owner=%s"
)

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
	fmt.Printf("%d peers", len(v.Peers))
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

func (v *Value) Format(address string) string {
	address = strings.ToLower(address)
	if strings.HasPrefix(address, "0x") || strings.HasSuffix(address, ".eth") {
		return address
	}
	return address
}

// Helper function to rebuild the map from v.Peers
func (v *Value) rebuildMap() {
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	v.Map = make(map[string]*Peer)
	for _, p := range v.Peers {
		if p.Address != "" && p.Address != "." && p.Address != "!" {
			v.Map[p.Address] = p
		}
		if p.ENS != "" && p.ENS != "." && p.ENS != "!" {
			v.Map[p.ENS] = p
		}
		if p.LoopringENS != "" && p.LoopringENS != "." && p.LoopringENS != "!" {
			v.Map[p.LoopringENS] = p
		}
		if p.LoopringID != "" && p.LoopringID != "." && p.LoopringID != "!" {
			v.Map[p.LoopringID] = p
		}
	}
}

func (v *Value) input(url string, response any) error {
	data, err := v.Factory.Json.In(url, "")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, response)
}
