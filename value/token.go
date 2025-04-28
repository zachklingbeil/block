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

	tokens := make([]Token, 0, len(source))
	for _, tokenJSON := range source {
		var token Token
		if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
			log.Printf("Skipping invalid token: %v (data: %s)", err, tokenJSON)
			continue
		}
		tokens = append(tokens, token)
	}
	v.Tokens = tokens
	v.Factory.State.Add("tokens", len(v.Tokens))
	return nil
}
