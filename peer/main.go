package peer

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wealdtech/go-ens/v3"
	"github.com/zachklingbeil/factory"
)

//go:embed loopring.json
var peersJSON []byte

type Peer struct {
	ENS         string `json:"ens,omitempty"`
	LoopringENS string `json:"loopringEns,omitempty"`
	LoopringID  string `json:"loopringId,omitempty"`
	Address     string `json:"address,omitempty"`
}

func NewPeers(factory *factory.Factory) {
	var peersData []Peer
	if err := json.Unmarshal(peersJSON, &peersData); err != nil {
		log.Fatalf("Failed to unmarshal peers: %v", err)
	}

	// Resolve ENS for all peers
	// Resolve ENS for all peers
	for i := range peersData {
		ensName, err := ens.ReverseResolve(factory.Eth, common.HexToAddress(peersData[i].Address))
		if err != nil || ensName == "" {
			peersData[i].ENS = "."
		} else {
			peersData[i].ENS = strings.ToLower(ensName)
		}
		fmt.Printf("%d %v %v %v %v\n", i,
			peersData[i].ENS,
			peersData[i].LoopringENS,
			peersData[i].LoopringID,
			peersData[i].Address)
	}
	// Store all peers in Redis after ENS resolution
	for _, peer := range peersData {
		peerJSON, err := json.Marshal(peer)
		if err != nil {
			log.Printf("Failed to marshal peer: %v", err)
			continue
		}
		if err := factory.Data.RB.SAdd(factory.Ctx, "peers", peerJSON).Err(); err != nil {
			log.Printf("Failed to add peer to Redis: %v", err)
		}
	}
	fmt.Printf("%d peers\n", len(peersData))
}
