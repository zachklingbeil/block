package circuit

import (
	"encoding/json"
	"log"
)

func (c *Circuit) LoadTokens() error {
	source, err := c.Factory.Redis.SMembers(c.Factory.Ctx, "tokens").Result()
	if err != nil {
		log.Fatalf("Failed to fetch tokens from Redis: %v", err)
	}

	for _, tokenJSON := range source {
		var token Token
		if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
			log.Printf("Skipping invalid token: %v (data: %s)", err, tokenJSON)
			continue
		}
		c.Tokens = append(c.Tokens, token)

		c.IDMap[token.TokenId] = &c.Tokens[len(c.Tokens)-1]
		c.LPMap[token.LoopringID] = &c.Tokens[len(c.Tokens)-1]

		c.Map[token.Token] = &c.Tokens[len(c.Tokens)-1]
		c.Map[token.Address] = &c.Tokens[len(c.Tokens)-1]
	}
	c.Factory.State.Add("tokens", len(c.Tokens))
	return nil
}

func (c *Circuit) GetToken(tokenId int64) *Token {
	if token, exists := c.IDMap[tokenId]; exists {
		return token
	}
	return nil
}

func (c *Circuit) GetLP(loopringId int64) *Token {
	if token, exists := c.LPMap[loopringId]; exists {
		return token
	}
	return nil
}

// package token

// import (
// 	_ "embed"
// 	"encoding/json"
// 	"fmt"
// 	"log"

// 	"github.com/zachklingbeil/factory"
// )

// type Tokens struct {
// 	Factory *factory.Factory
// }

// type Token struct {
// 	Token      string `json:"token"`
// 	TokenId    int64  `json:"tokenId"`
// 	LoopringID int64  `json:"loopringId,omitempty"`
// 	Decimals   int64  `json:"decimals"`
// 	Address    string `json:"address"`
// }

// //go:embed tokens.json
// var tokens []byte

// func NewTokens(factory *factory.Factory) {
// 	var in []Token
// 	if err := json.Unmarshal(tokens, &in); err != nil {
// 		log.Fatalf("Failed to unmarshal tokens: %v", err)
// 	}
// 	for _, token := range in {
// 		tokenJSON, _ := json.Marshal(token)
// 		factory.Redis.SAdd(factory.Ctx, "tokens", tokenJSON).Err()
// 	}

// 	fmt.Printf("%d tokens\n", len(in))
// }
