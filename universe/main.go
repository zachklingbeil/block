package universe

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zachklingbeil/factory"
)

type Universe struct {
	Factory *factory.Factory
	Map     map[string]map[*common.Address]*any
}

// NewUniverse initializes a new Universe instance
func NewUniverse(factory *factory.Factory) *Universe {
	return &Universe{
		Factory: factory,
		Map:     make(map[string]map[*common.Address]*any),
	}
}

// BuildMapFromValueSet populates the nested map from Redis sets "peer" and "token"
func (u *Universe) BuildMapFromValueSet() error {
	// Load peers into the map
	if err := u.loadPeers(); err != nil {
		return fmt.Errorf("failed to load peers: %v", err)
	}

	// Load tokens into the map
	if err := u.loadTokens(); err != nil {
		return fmt.Errorf("failed to load tokens: %v", err)
	}
	return nil
}

// loadPeers loads Peer objects into the nested map
func (u *Universe) loadPeers() error {
	peers, err := u.Factory.Data.RB.SMembers(u.Factory.Ctx, "peer").Result()
	if err != nil {
		return fmt.Errorf("failed to fetch peers from Redis set: %v", err)
	}

	for _, peerJSON := range peers {
		var peer Peer
		if err := json.Unmarshal([]byte(peerJSON), &peer); err != nil {
			log.Printf("Skipping invalid peer: %v (data: %s)", err, peerJSON)
			continue
		}

		// Add Peer to the map using its fields as keys
		if err := u.addPeerToMap(&peer); err != nil {
			log.Printf("Skipping peer: %v", err)
			continue
		}
	}

	return nil
}

// loadTokens loads Token objects into the nested map
func (u *Universe) loadTokens() error {
	tokens, err := u.Factory.Data.RB.SMembers(u.Factory.Ctx, "token").Result()
	if err != nil {
		return fmt.Errorf("failed to fetch tokens from Redis set: %v", err)
	}

	for _, tokenJSON := range tokens {
		var token Token
		if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
			log.Printf("Skipping invalid token: %v (data: %s)", err, tokenJSON)
			continue
		}

		// Add Token to the map using its fields as keys
		if err := u.addTokenToMap(&token); err != nil {
			log.Printf("Skipping token: %v", err)
			continue
		}
	}

	return nil
}

// addPeerToMap adds a Peer object to the nested map
func (u *Universe) addPeerToMap(peer *Peer) error {
	if !common.IsHexAddress(peer.Address) {
		return fmt.Errorf("invalid Ethereum address in peer: %s", peer.Address)
	}

	// Add Peer to the map using its fields as keys
	address := common.HexToAddress(peer.Address)
	u.addToNestedMap(peer.ENS, &address, peer)
	u.addToNestedMap(peer.LoopringENS, &address, peer)
	u.addToNestedMap(peer.LoopringID, &address, peer)

	return nil
}

// addTokenToMap adds a Token object to the nested map
func (u *Universe) addTokenToMap(token *Token) error {
	if !common.IsHexAddress(token.Address) {
		return fmt.Errorf("invalid Ethereum address in token: %s", token.Address)
	}

	// Add Token to the map using its fields as keys
	address := common.HexToAddress(token.Address)
	u.addToNestedMap(token.Token, &address, token)
	u.addToNestedMap(token.TokenId, &address, token)

	return nil
}

// addToNestedMap adds an object to the nested map
func (u *Universe) addToNestedMap(key string, address *common.Address, obj any) {
	if key == "" {
		return // Skip empty keys
	}

	if _, exists := u.Map[key]; !exists {
		u.Map[key] = make(map[*common.Address]*any)
	}

	// Store a pointer to the object
	u.Map[key][address] = &obj
	log.Printf("Added object to map with key: %s and address: %s", key, address.Hex())
}

// GetAddressByKey retrieves a single address associated with a given key
func (u *Universe) GetAddressByKey(key string) (*common.Address, error) {
	objects, exists := u.Map[key]
	if !exists {
		return nil, fmt.Errorf("no objects found for key: %s", key)
	}

	for address := range objects {
		return address, nil // Return the first address found
	}

	return nil, fmt.Errorf("no addresses found for key: %s", key)
}

// GetObjectByAddress retrieves an object (Peer or Token) using an address
func (u *Universe) GetObjectByAddress(address *common.Address) (any, error) {
	for key, objects := range u.Map {
		if obj, exists := objects[address]; exists {
			log.Printf("Found object for address: %s under key: %s", address.Hex(), key)
			return *obj, nil
		}
	}

	return nil, fmt.Errorf("no object found for address: %s", address.Hex())
}

// Peer struct represents a peer with various properties
type Peer struct {
	ENS         string `json:"ens,omitempty"`
	LoopringENS string `json:"loopringEns,omitempty"`
	LoopringID  string `json:"loopringId,omitempty"`
	Address     string `json:"address,omitempty"`
	FirstBlock  string `json:"firstBlock,omitempty"`
}

// Token struct represents a token with various properties
type Token struct {
	Token    string `json:"token,omitempty"`
	Address  string `json:"address,omitempty"`
	Decimals string `json:"decimals,omitempty"`
	TokenId  string `json:"tokenId,omitempty"`
	TokenInt int64  `json:"tokenInt,omitempty"`
}
