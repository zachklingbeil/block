package value

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
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
		v.TokenMap[token.TokenId] = &token
		v.TokenMap[token.Token] = &token
	}
	return nil
}

func (v *Value) GetTokenById(tokenId any) *Token {
	v.Factory.Rw.RLock()
	defer v.Factory.Rw.RUnlock()

	var token *Token
	var exists bool

	switch id := tokenId.(type) {
	case int64:
		token, exists = v.TokenMap[id]
	case string:
		token, exists = v.TokenMap[id]
	default:
		// log.Printf("Unsupported token ID type: %T", tokenId)
		return &Token{Token: fmt.Sprintf("%v", tokenId)} // Return a default token with the ID as a string
	}

	if !exists {
		// log.Printf("Token not found for ID: %v", tokenId)
		return &Token{Token: fmt.Sprintf("%v", tokenId)} // Return a default token with the ID as a string
	}
	return token
}

// FormatValue formats a string input as a decimal string based on the token's decimals.
func (v *Value) FormatValue(input string, key any) string {
	v.Factory.Rw.RLock()
	token, exists := v.TokenMap[key]
	v.Factory.Rw.RUnlock()
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
