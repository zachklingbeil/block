package value

import (
	"github.com/zachklingbeil/factory"
)

type Value struct {
	Factory  *factory.Factory
	Peers    []*Peer
	Map      map[string]*Peer
	TokenMap map[any]*Token
}

type Peer struct {
	ENS         string `json:"ens,omitempty"`
	LoopringENS string `json:"loopringEns,omitempty"`
	LoopringID  string `json:"loopringId,omitempty"`
	Address     string `json:"address,omitempty"`
	FirstBlock  string `json:"firstBlock,omitempty"`
}

type Token struct {
	Token    string `json:"token,omitempty"`
	Address  string `json:"address,omitempty"`
	Decimals string `json:"decimals,omitempty"`
	TokenId  string `json:"tokenId,omitempty"`
	TokenInt int64  `json:"tokenInt,omitempty"`
}

func NewValue(factory *factory.Factory) *Value {
	v := &Value{
		Factory:  factory,
		Map:      make(map[string]*Peer),
		TokenMap: make(map[any]*Token),
	}

	v.LoadTokens()
	v.LoadPeers()
	// v.UpdatePeersFromEmbeddedData()
	// v.Factory.State.Count("peers", len(v.Peers), true)
	return v
}

func (v *Value) GetFromMap(key string) *Peer {
	v.Factory.Rw.RLock()
	defer v.Factory.Rw.RUnlock()
	peer, exists := v.Map[key]
	if !exists {
		return nil
	}
	return peer
}

// func (v *Value) SyncPeersAndTokensToRedis() error {
// 	ctx := v.Factory.Ctx
// 	client := v.Factory.Data.RB
// 	hashKey := "universe" // Single hash key for all data

// 	// Iterate over peers and add them to the "universe" hash
// 	for _, peer := range v.Peers {
// 		if peer.Address != "" {
// 			peerJSON, err := json.Marshal(peer)
// 			if err != nil {
// 				return fmt.Errorf("failed to marshal peer: %v", err)
// 			}
// 			err = client.HSet(ctx, hashKey, peer.Address, peerJSON).Err()
// 			if err != nil {
// 				return fmt.Errorf("failed to sync peer to Redis: %v", err)
// 			}
// 		}
// 	}

// 	// Iterate over tokens and add them to the "universe" hash
// 	for _, token := range v.Tokens {
// 		if token.Address != "" {
// 			tokenJSON, err := json.Marshal(token)
// 			if err != nil {
// 				return fmt.Errorf("failed to marshal token: %v", err)
// 			}
// 			err = client.HSet(ctx, hashKey, token.Address, tokenJSON).Err()
// 			if err != nil {
// 				return fmt.Errorf("failed to sync token to Redis: %v", err)
// 			}
// 		}
// 	}

// 	fmt.Println("Successfully synced peers and tokens to Redis under the 'universe' hash.")
// 	return nil
// }
