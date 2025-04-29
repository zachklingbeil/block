package value

import (
	"encoding/json"
	"log"
)

type Token struct {
	Token    string `json:"token,omitempty"`
	Address  string `json:"address,omitempty"`
	Decimals string `json:"decimals,omitempty"`
	TokenId  string `json:"tokenId,omitempty"`
	TokenInt int64  `json:"tokenInt,omitempty"`
}

func (v *Value) LoadTokens() error {
	source, err := v.Factory.Data.RB.SMembers(v.Factory.Ctx, "token").Result()
	if err != nil {
		log.Fatalf("Failed to fetch tokens from Redis: %v", err)
	}

	v.Tokens = make([]Token, 0, len(source))
	for _, tokenJSON := range source {
		var token Token
		if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
			log.Printf("Skipping invalid token: %v (data: %s)", err, tokenJSON)
			continue
		}
		v.Tokens = append(v.Tokens, token)
		v.TokenMap[token.TokenInt] = &token
	}
	return nil
}

func (v *Value) TokenIdToString(tokenInt int64) string {
	v.Factory.Rw.RLock()
	defer v.Factory.Rw.RUnlock()

	token, exists := v.TokenMap[tokenInt]
	if exists {
		return token.TokenId
	}
	return "!"
}

func (v *Value) GetTokenById(tokenInt int64) *Token {
	v.Factory.Rw.RLock()
	defer v.Factory.Rw.RUnlock()

	token, exists := v.TokenMap[tokenInt]
	if !exists {
		log.Printf("Token not found for ID: %d", tokenInt) // Log the missing tokenInt
		return nil
	}
	return token
}
