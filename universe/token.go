package universe

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
)

type Token struct {
	Token    string `json:"token,omitempty"`
	Address  string `json:"address,omitempty"`
	Decimals int64  `json:"decimals,omitempty"`
	TokenId  int64  `json:"tokenId,omitempty"`
	ABI      string `json:"abi,omitempty"`
}

func (o *One) LoadTokens() error {
	source, err := o.Factory.Data.RB.SMembers(o.Factory.Ctx, "token").Result()
	if err != nil {
		return fmt.Errorf("failed to fetch tokens from Redis set: %v", err)
	}
	o.Tokens = make([]*Token, 0, len(source))
	o.Maps.TokenId = make(map[int64]string)
	for _, tokenJSON := range source {
		var token Token
		if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
			log.Printf("Skipping invalid token: %v (data: %s)", err, tokenJSON)
			continue
		}
		o.Tokens = append(o.Tokens, &token)
		o.Map[token.Address] = &token
		o.Maps.TokenId[token.TokenId] = token.Address
	}
	fmt.Printf("%d tokens loaded\n", len(o.Tokens))
	o.SaveTokensToFile("tokens.json")
	return nil
}

func (o *One) SaveTokensToFile(filename string) error {
	o.Factory.Rw.RLock()
	defer o.Factory.Rw.RUnlock()
	data, err := json.MarshalIndent(o.Tokens, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tokens: %v", err)
	}
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write tokens to file: %v", err)
	}
	return nil
}

// AddToken adds a new token to the Tokens struct and updates all maps.
func (o *One) AddToken(token *Token) {
	o.Factory.Rw.Lock()
	defer o.Factory.Rw.Unlock()
	o.Tokens = append(o.Tokens, token)
	o.Map[token.Address] = token
	o.Maps.TokenId[token.TokenId] = token.Address
}

// GetAddress returns the common.Address for a given tokenId.
func (o *One) GetAddress(tokenId int64) string {
	if tokenId >= 500 {
		return strconv.FormatInt(tokenId, 10)
	}

	o.Factory.Rw.RLock()
	defer o.Factory.Rw.RUnlock()

	token, exists := o.Maps.TokenId[tokenId]
	if !exists {
		log.Printf("Token not found for ID: %d", tokenId)
		return strconv.FormatInt(tokenId, 10)
	}
	return token
}

// Format formats a string input as a decimal string based on the token's decimals, using address.
func (o *One) Format(input string, address string) string {

	o.Factory.Rw.RLock()
	token, exists := o.Map[address]
	o.Factory.Rw.RUnlock()
	if !exists {
		return input
	}
	return format(input, token.(*Token))
}

// Helper function to format value with token decimals.
func format(input string, token *Token) string {
	value := new(big.Int)
	_, ok := value.SetString(input, 10)
	if !ok {
		return input
	}

	decimals := int(token.Decimals)

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
