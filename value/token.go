package value

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
)

func (v *Value) LoadTokens() error {
	// Step 1: Fetch all tokens from the Redis hash
	hashKey := "token" // The Redis hash key used in MigrateTokens
	source, err := v.Factory.Data.RB.HGetAll(v.Factory.Ctx, hashKey).Result()
	if err != nil {
		return fmt.Errorf("failed to fetch tokens from Redis hash: %v", err)
	}

	// Step 2: Clear existing tokens and initialize the token map
	v.Tokens = make([]*Token, 0, len(source))
	v.TokenMap = make(map[any]*Token)

	// Step 3: Deserialize each token and populate the token map
	for _, tokenJSON := range source {
		var token Token
		if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
			log.Printf("Skipping invalid token: %v (data: %s)", err, tokenJSON)
			continue
		}

		// Add the token to the list and map
		v.Tokens = append(v.Tokens, &token)
		v.TokenMap[token.TokenInt] = &token
		v.TokenMap[token.TokenId] = &token
		v.TokenMap[token.Token] = &token
	}

	fmt.Printf("%d tokens", len(v.Tokens))
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
