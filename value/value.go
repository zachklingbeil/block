package value

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/zachklingbeil/factory"
)

type Value struct {
	Factory  *factory.Factory
	Peers    []*Peer
	Tokens   []*Token
	Map      map[string]*Peer
	TokenMap map[any]*Token
}

func NewValue(factory *factory.Factory) *Value {
	v := &Value{
		Factory:  factory,
		Map:      make(map[string]*Peer),
		TokenMap: make(map[any]*Token),
	}

	v.LoadTokens()
	v.LoadPeers()
	// v.MigratePeers()

	// v.rebuildMap()
	// v.Refresh()
	return v
}

func (v *Value) Refresh() {
	remaining := len(v.Peers)
	for i := range v.Peers {
		fmt.Printf("%d", i)
		peer := v.Peers[i]
		// Clear invalid fields if they are "." or "!"
		if peer.ENS == "." || peer.ENS == "!" {
			peer.ENS = ""
		}
		if peer.LoopringID == "." || peer.LoopringID == "!" {
			peer.LoopringID = ""
		}
		if peer.LoopringENS == "." || peer.LoopringENS == "!" {
			peer.LoopringENS = ""
		}
		if peer.Address == "." || peer.Address == "!" {
			peer.Address = ""
		}
		v.HelloUniverse(peer.Address)
		remaining--
		fmt.Printf("\n%d", remaining)
	}
}

func (v *Value) MigrateCombined() error {
	// Define the Redis hash key for the combined table
	hashKey := "universe" // The new Redis hash key

	// Step 1: Migrate Peers
	for _, peer := range v.Peers {
		if peer.Address == "" {
			log.Printf("Skipping peer with empty address: %+v", peer)
			continue
		}

		peerJSON, err := json.Marshal(peer)
		if err != nil {
			log.Printf("Failed to serialize peer: %v", err)
			continue
		}

		// Store the peer in the combined Redis hash
		if err := v.Factory.Data.RB.HSet(v.Factory.Ctx, hashKey, peer.Address, peerJSON).Err(); err != nil {
			log.Printf("Failed to store peer in combined Redis hash %s with key %s: %v", hashKey, peer.Address, err)
			continue
		}
	}

	// Step 2: Migrate Tokens
	for _, token := range v.Tokens {
		if token.Address == "" {
			log.Printf("Skipping token with empty address: %+v", token)
			continue
		}

		tokenJSON, err := json.Marshal(token)
		if err != nil {
			log.Printf("Failed to serialize token: %v", err)
			continue
		}

		// Store the token in the combined Redis hash
		if err := v.Factory.Data.RB.HSet(v.Factory.Ctx, hashKey, token.Address, tokenJSON).Err(); err != nil {
			log.Printf("Failed to store token in combined Redis hash %s with key %s: %v", hashKey, token.Address, err)
			continue
		}
	}

	log.Printf("Migrated peers and tokens to the combined Redis hash: %s", hashKey)
	return nil
}
