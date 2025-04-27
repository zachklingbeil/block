package manual

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"

	"github.com/zachklingbeil/factory"
)

//go:embed peer.json
var peers []byte

type LP struct {
	LoopringENS string `json:"loopringEns"`
	LoopringID  string `json:"loopringId"`
	Address     string `json:"address"`
	ENS         string `json:"ens"`
}

func NewLP(factory *factory.Factory) {
	var data []LP
	if err := json.Unmarshal(peers, &data); err != nil {
		log.Fatalf("Failed to unmarshal peers: %v", err)
	}
	for _, peer := range data {
		tokenJSON, err := json.Marshal(peer)
		if err != nil {
			log.Printf("Failed to marshal peer: %v", err)
			continue
		}

		if err := factory.Data.RB.SAdd(factory.Ctx, "peer", tokenJSON).Err(); err != nil {
			log.Printf("Failed to add token to Redis: %v", err)
		}
	}
	fmt.Printf("%d peers\n", len(data))
}
