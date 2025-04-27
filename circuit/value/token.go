package value

import (
	"encoding/json"
	"log"

	"github.com/zachklingbeil/factory"
)

type Token struct {
	Token    string `json:"token,omitempty"`
	Address  string `json:"address,omitempty"`
	Decimals string `json:"decimals,omitempty"`
	TokenId  string `json:"tokenId,omitempty"`
	TokenInt int64  `json:"tokenInt,omitempty"`
}

func NewTokens(factory *factory.Factory) []Token {
	var tokens []Token
	source, err := factory.Data.RB.SMembers(factory.Ctx, "token").Result()
	if err != nil {
		log.Fatalf("Failed to fetch tokens from Redis: %v", err)
	}
	for _, tokenJSON := range source {
		var token Token
		if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
			log.Printf("Skipping invalid token: %v (data: %s)", err, tokenJSON)
			continue
		}
		tokens = append(tokens, token)
	}
	return tokens
}
