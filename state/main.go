package state

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/zachklingbeil/factory"
)

type Tokens struct {
	Factory *factory.Factory
	Slice   []Token
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
	}

	source, err := factory.Data.RB.SMembers(factory.Ctx, "tokens").Result()
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
	}
	factory.State.Add("tokens", len(t.Slice))
	return t
}

func NewToken[T any](ctx context.Context, rb *redis.Client) []T {
	var tokens []T
	source, err := rb.SMembers(ctx, "tokens").Result()
	if err != nil {
		log.Fatalf("Failed to fetch tokens from Redis: %v", err)
	}

	for _, s := range source {
		var t T
		if err := json.Unmarshal([]byte(s), &t); err != nil {
			log.Printf("Skipping invalid token: %v (data: %s)", err, s)
			continue
		}
		tokens = append(tokens, t)
	}
	return tokens
}

func SliceIn[T any](ctx context.Context, rb *redis.Client, key string) []T {
	var items []T
	source, err := rb.SMembers(ctx, key).Result()
	if err != nil {
		log.Fatalf("Failed to fetch items from Redis set '%s': %v", key, err)
	}

	for _, s := range source {
		var item T
		if err := json.Unmarshal([]byte(s), &item); err != nil {
			log.Printf("Skipping invalid item: %v (data: %s)", err, s)
			continue
		}
		items = append(items, item)
	}
	return items
}
