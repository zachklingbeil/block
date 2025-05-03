package value

import (
	"encoding/json"
	"fmt"
	"log"
)

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

type Values struct {
	ENS         string `json:"ens,omitempty"`
	Token       string `json:"token,omitempty"`
	LoopringENS string `json:"loopringEns,omitempty"`
	LoopringID  string `json:"loopringId,omitempty"`
	Address     string `json:"address,omitempty"`
	FirstBlock  string `json:"firstBlock,omitempty"`
	TokenId     int64  `json:"tokenId,omitempty"`
	Decimals    string `json:"decimals,omitempty"`
}

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
	return nil
}

func (v *Value) LoadTokens() error {
	hashKey := "token"
	source, err := v.Factory.Data.RB.HGetAll(v.Factory.Ctx, hashKey).Result()
	if err != nil {
		return fmt.Errorf("failed to fetch tokens from Redis hash: %v", err)
	}
	v.Tokens = make([]*Token, 0, len(source))

	for _, tokenJSON := range source {
		var token Token
		if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
			log.Printf("Skipping invalid token: %v (data: %s)", err, tokenJSON)
			continue
		}
		v.Tokens = append(v.Tokens, &token)
	}
	return nil
}

func (v *Value) ConsolidateAndStoreValues() error {
	// Ensure thread safety
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	// Consolidate Peers and Tokens into Values
	values := make([]Values, 0, len(v.Peers)+len(v.Tokens))

	for _, peer := range v.Peers {
		values = append(values, Values{
			ENS:         peer.ENS,
			LoopringENS: peer.LoopringENS,
			LoopringID:  peer.LoopringID,
			Address:     peer.Address,
			FirstBlock:  peer.FirstBlock,
		})
	}

	for _, token := range v.Tokens {
		values = append(values, Values{
			Token:    token.Token,
			Address:  token.Address,
			Decimals: token.Decimals,
			TokenId:  token.TokenInt,
		})
	}

	for _, value := range values {
		valueJSON, err := json.Marshal(value)
		if err != nil {
			log.Printf("Failed to serialize value: %v (error: %v)", value, err)
			continue
		}

		if err := v.Factory.Data.RB.SAdd(v.Factory.Ctx, "value", valueJSON).Err(); err != nil {
			return fmt.Errorf("failed to store value in Redis set: %v", err)
		}
	}

	return nil
}
