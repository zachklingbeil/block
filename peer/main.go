package peer

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zachklingbeil/factory"
)

type Peer struct {
	ENS         string         `json:"ens,omitempty"`
	LoopringENS string         `json:"loopringEns,omitempty"`
	LoopringID  int64          `json:"loopringId,omitempty"`
	Address     common.Address `json:"address,omitempty"`
}

func NewPeers(factory *factory.Factory) {
	// Load peers from Redis set "peers"
	peerStrs, err := factory.Data.RB.SMembers(factory.Ctx, "peers").Result()
	if err != nil {
		log.Fatalf("Failed to load peers from Redis: %v", err)
	}

	var peersData []Peer
	for _, peerStr := range peerStrs {
		var peer Peer
		if err := json.Unmarshal([]byte(peerStr), &peer); err != nil {
			log.Printf("Failed to unmarshal peer: %v", err)
			continue
		}
		peer.Address = common.HexToAddress(strings.ToLower(peer.Address.Hex()))
		peersData = append(peersData, peer)
	}

	// Store peers in two Redis hash sets: by address and by loopringId
	for _, peer := range peersData {
		peerJSON, err := json.Marshal(peer)
		if err != nil {
			log.Printf("Failed to marshal peer: %v", err)
			continue
		}
		// Add peer object to peers:peers set (unhashed)
		if err := factory.Data.RB.SAdd(factory.Ctx, "peer", peerJSON).Err(); err != nil {
			log.Printf("Failed to add peer to Redis set peers:peers: %v", err)
		}
	}
	fmt.Printf("%d peers loaded and indexed\n", len(peersData))
}
