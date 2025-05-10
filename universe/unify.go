package universe

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type One struct {
	ENS         string `json:"ens,omitempty"`
	LoopringENS string `json:"loopringEns,omitempty"`
	LoopringID  int64  `json:"loopringId,omitempty"`
	Address     string `json:"address"`
	Token       string `json:"token,omitempty"`
	Decimals    int64  `json:"decimals,omitempty"`
	TokenId     int64  `json:"tokenId,omitempty"`
	ABI         string `json:"abi,omitempty"`
}

type Peer struct {
	ENS         string `json:"ens,omitempty"`
	LoopringENS string `json:"loopringEns,omitempty"`
	LoopringID  int64  `json:"loopringId,omitempty"`
	Address     string `json:"address,omitempty"`
}

type Token struct {
	Token    string `json:"token,omitempty"`
	Address  string `json:"address,omitempty"`
	Decimals int64  `json:"decimals,omitempty"`
	TokenId  int64  `json:"tokenId,omitempty"`
	ABI      string `json:"abi,omitempty"`
}

func (z *Zero) MergePeers() ([]Peer, error) {
	z.Factory.Rw.Lock()
	defer z.Factory.Rw.Unlock()

	source, err := z.Factory.Data.RB.SMembers(z.Factory.Ctx, "peer").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch peers from Redis hash: %v", err)
	}
	peers := make([]Peer, 0, len(source))
	for _, peerJSON := range source {
		var peer Peer
		if err := json.Unmarshal([]byte(peerJSON), &peer); err != nil {
			log.Printf("Skipping invalid peer: %v (data: %s)", err, peerJSON)
			continue
		}
		peers = append(peers, peer)
	}
	return peers, nil
}

func (z *Zero) MergeTokens() ([]Token, error) {
	z.Factory.Rw.Lock()
	defer z.Factory.Rw.Unlock()

	source, err := z.Factory.Data.RB.SMembers(z.Factory.Ctx, "token").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch peers from Redis hash: %v", err)
	}
	tokens := make([]Token, 0, len(source))
	for _, tokenJSON := range source {
		var peer Token
		if err := json.Unmarshal([]byte(tokenJSON), &peer); err != nil {
			log.Printf("Skipping invalid peer: %v (data: %s)", err, tokenJSON)
			continue
		}
		tokens = append(tokens, peer)
	}
	return tokens, nil
}

func (z *Zero) ConsolidateOnes() ([]One, error) {
	peers, err := z.MergePeers()
	if err != nil {
		return nil, fmt.Errorf("failed to merge peers: %v", err)
	}
	tokens, err := z.MergeTokens()
	if err != nil {
		return nil, fmt.Errorf("failed to merge tokens: %v", err)
	}

	onesMap := make(map[string]*One)

	for _, p := range peers {
		addr := strings.ToLower(p.Address)
		if addr == "" {
			continue
		}
		one, exists := onesMap[addr]
		if !exists {
			one = &One{Address: addr}
			onesMap[addr] = one
		}
		one.ENS = p.ENS
		one.LoopringENS = p.LoopringENS
		one.LoopringID = p.LoopringID
	}

	for _, t := range tokens {
		addr := strings.ToLower(t.Address)
		if addr == "" {
			continue
		}
		one, exists := onesMap[addr]
		if !exists {
			one = &One{Address: addr}
			onesMap[addr] = one
		}
		one.Token = t.Token
		one.Decimals = t.Decimals
		one.TokenId = t.TokenId
		one.ABI = t.ABI
	}

	ones := make([]One, 0, len(onesMap))
	oneJSONs := make([]any, 0, len(onesMap))
	for _, one := range onesMap {
		ones = append(ones, *one)
		b, err := json.Marshal(one)
		if err != nil {
			log.Printf("Failed to marshal One: %v", err)
			continue
		}
		oneJSONs = append(oneJSONs, string(b))
	}

	// Save to Redis set "one"
	z.Factory.Rw.Lock()
	defer z.Factory.Rw.Unlock()
	if err := z.Factory.Data.RB.SAdd(z.Factory.Ctx, "one", oneJSONs...).Err(); err != nil {
		return nil, fmt.Errorf("failed to save ones to Redis: %v", err)
	}

	return ones, nil
}
