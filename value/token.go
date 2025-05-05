package value

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
)

func (v *Value) SyncTokensToRedis() error {
	ctx := v.Factory.Ctx
	client := v.Factory.Data.RB
	hashKey := "tokenMap" // Redis hash key for tokens

	// Iterate over the TokenMap
	for key, token := range v.TokenMap {
		tokenJSON, err := json.Marshal(token)
		if err != nil {
			log.Printf("Failed to marshal token: %v (key: %v)", err, key)
			continue
		}

		err = client.HSet(ctx, hashKey, fmt.Sprintf("%v", key), tokenJSON).Err()
		if err != nil {
			return fmt.Errorf("failed to sync token to Redis: %v (key: %v)", err, key)
		}
	}

	fmt.Printf("Successfully synced %d tokens to Redis under the '%s' hash.\n", len(v.TokenMap), hashKey)
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
