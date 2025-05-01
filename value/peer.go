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

type Peer struct {
	ENS         string `json:"ens"`
	LoopringENS string `json:"loopringEns"`
	LoopringID  string `json:"loopringId"`
	Address     string `json:"address"`
}

func (v *Value) LoadPeers() error {
	// Lock while updating v.Peers
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	hashKey := "peer" // The Redis hash key used in MigratePeers
	source, err := v.Factory.Data.RB.HGetAll(v.Factory.Ctx, hashKey).Result()
	if err != nil {
		return fmt.Errorf("failed to fetch peers from Redis hash: %v", err)
	}

	// Clear any existing peers before loading fresh ones
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

func (v *Value) Save(peer *Peer) error {
	// Serialize the peer to JSON
	peerJSON, err := json.Marshal(peer)
	if err != nil {
		return fmt.Errorf("failed to serialize peer: %v", err)
	}

	// Define the Redis hash key for storing peers
	hashKey := "peer"

	// Use the peer's Address as the field in the Redis hash
	if err := v.Factory.Data.RB.HSet(v.Factory.Ctx, hashKey, peer.Address, peerJSON).Err(); err != nil {
		return fmt.Errorf("failed to store peer in Redis hash: %v", err)
	}
	return nil
}

func (v *Value) Hello(value string) string {
	v.Factory.Rw.RLock()
	peer, exists := v.Map[value]
	v.Factory.Rw.RUnlock()
	if !exists {
		return ""
	}
	if peer.ENS != "" && peer.ENS != "." && peer.ENS != "!" {
		return peer.ENS
	}
	if peer.LoopringENS != "" && peer.LoopringENS != "." && peer.LoopringENS != "!" {
		return peer.LoopringENS
	}
	return peer.Address
}

func (v *Value) Format(address string) string {
	address = strings.ToLower(address)
	if strings.HasPrefix(address, "0x") || strings.HasSuffix(address, ".eth") {
		return address
	}
	return address
}

func (v *Value) input(url string, response any) error {
	data, err := v.Factory.Json.In(url, "")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, response)
}
