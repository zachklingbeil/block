package token

import (
	"encoding/json"
	"log"

	"github.com/zachklingbeil/factory"
)

type Tokens struct {
	Factory *factory.Factory
	Slice   []Token
	ID      map[int64]*Token
	LP      map[int64]*Token
}

type Token struct {
	Token      string `json:"token"`
	TokenId    int64  `json:"tokenId"`
	LoopringID int64  `json:"loopringId,omitempty"`
	Decimals   int64  `json:"decimals"`
	Address    string `json:"address"`
}

func NewTokens(factory *factory.Factory) *Tokens {
	t := &Tokens{
		Factory: factory,
		Slice:   make([]Token, 270),
		ID:      make(map[int64]*Token),
		LP:      make(map[int64]*Token),
	}

	source, err := factory.Redis.SMembers(factory.Ctx, "tokens").Result()
	if err != nil {
		log.Fatalf("Failed to fetch tokens from Redis: %v", err)
	}

	for _, tokenJSON := range source {
		var token Token
		if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
			log.Printf("Skipping invalid token: %v (data: %s)", err, tokenJSON)
			continue
		}
		t.Slice = append(t.Slice, token)
		t.ID[token.TokenId] = &t.Slice[len(t.Slice)-1]
		t.LP[token.LoopringID] = &t.Slice[len(t.Slice)-1]
	}
	log.Printf("%d tokens", len(t.Slice))
	return t
}

func (t *Tokens) GetToken(tokenId int64) *Token {
	if token, exists := t.ID[tokenId]; exists {
		return token
	}
	return nil
}

func (t *Tokens) GetLP(loopringId int64) *Token {
	if token, exists := t.LP[loopringId]; exists {
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
