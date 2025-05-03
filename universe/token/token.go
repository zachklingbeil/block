package token

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/zachklingbeil/factory"
)

type Tokens struct {
	Factory  *factory.Factory
	Tokens   []*Token
	TokenMap map[any]*Token
}
type Token struct {
	Token    string `json:"token,omitempty"`
	Address  string `json:"address,omitempty"`
	Decimals string `json:"decimals,omitempty"`
	TokenId  string `json:"tokenId,omitempty"`
	TokenInt int64  `json:"tokenInt,omitempty"`
}

func (t *Tokens) LoadTokens(ctx context.Context, redis *redis.Client) error {
	hashKey := "token"
	source, err := redis.HGetAll(ctx, hashKey).Result()
	if err != nil {
		return fmt.Errorf("failed to fetch tokens from Redis hash: %v", err)
	}

	t.Tokens = make([]*Token, 0, len(source))
	t.TokenMap = make(map[any]*Token)

	for _, tokenJSON := range source {
		var token Token
		if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
			log.Printf("Skipping invalid token: %v (data: %s)", err, tokenJSON)
			continue
		}
		t.Tokens = append(t.Tokens, &token)
		t.TokenMap[token.TokenInt] = &token
	}
	return nil
}

func (t *Tokens) GetAddress(tokenId int64) string {
	t.Factory.Rw.RLock()
	defer t.Factory.Rw.RUnlock()

	token, exists := t.TokenMap[tokenId]
	if !exists {
		log.Printf("Token not found for ID: %d", tokenId)
		return ""
	}
	return token.Address
}

func (t *Tokens) GetTokenById(tokenId int64) *Token {
	// Lock the read mutex to ensure thread-safe access
	t.Factory.Rw.RLock()
	defer t.Factory.Rw.RUnlock()

	// Attempt to retrieve the token from the map using the int64 tokenId
	token, exists := t.TokenMap[tokenId]
	if !exists {
		// Log and return a default token if the tokenId is not found
		log.Printf("Token not found for ID: %d", tokenId)
		return &Token{Token: fmt.Sprintf("%d", tokenId)}
	}

	return token
}

// FormatValue formats a string input as a decimal string based on the token's decimals.
func (t *Tokens) FormatValue(input string, tokenInt int64) string {
	t.Factory.Rw.RLock()
	token, exists := t.TokenMap[tokenInt]
	t.Factory.Rw.RUnlock()
	if !exists {
		return input
	}

	value := new(big.Int)
	_, ok := value.SetString(input, 10)
	if !ok {
		return input
	}

	decimals, err := strconv.Atoi(token.Decimals)
	if err != nil {
		return input
	}

	valueStr := value.String()
	if len(valueStr) <= decimals {
		paddedValue := strings.Repeat("0", decimals-len(valueStr)+1) + valueStr
		result := "0." + paddedValue
		return strings.TrimRight(result, "0")
	}

	intPart := valueStr[:len(valueStr)-decimals]
	fracPart := valueStr[len(valueStr)-decimals:]
	result := intPart + "." + fracPart
	result = strings.TrimRight(result, "0")
	result = strings.TrimSuffix(result, ".")
	return result
}

// //go:embed token.json
// var tokens []byte

// func NewTokens(factory *factory.Factory) {
// 	var tokensData []Token
// 	if err := json.Unmarshal(tokens, &tokensData); err != nil {
// 		log.Fatalf("Failed to unmarshal tokens: %v", err)
// 	}
// 	for _, token := range tokensData {
// 		tokenJSON, err := json.Marshal(token)
// 		if err != nil {
// 			log.Printf("Failed to marshal token: %v", err)
// 			continue
// 		}

// 		if err := factory.Data.RB.SAdd(factory.Ctx, "token", tokenJSON).Err(); err != nil {
// 			log.Printf("Failed to add token to Redis: %v", err)
// 		}
// 	}
// 	// factory.State.Add("tokens", len(tokensData))
// 	fmt.Printf("%d tokens\n", len(tokensData))
// }
