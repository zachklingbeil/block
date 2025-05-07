package value

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type Token struct {
	Token    string         `json:"token,omitempty"`
	Address  common.Address `json:"address,omitempty"`
	Decimals string         `json:"decimals,omitempty"`
	TokenId  int64          `json:"tokenId,omitempty"`
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

func (v *Value) LoadTokens() error {
	source, err := v.Factory.Data.RB.SMembers(v.Factory.Ctx, "token").Result()
	if err != nil {
		return fmt.Errorf("failed to fetch tokens from Redis set: %v", err)
	}
	v.Tokens = make([]*Token, 0, len(source))
	v.Maps.TokenId = make(map[int64]*Token, len(source))
	for _, tokenJSON := range source {
		var token Token
		if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
			log.Printf("Skipping invalid token: %v (data: %s)", err, tokenJSON)
			continue
		}
		v.Tokens = append(v.Tokens, &token)
		v.Universe[token.Address] = &token
		v.Maps.TokenId[token.TokenId] = &token
	}
	fmt.Printf("%d tokens loaded\n", len(v.Tokens))
	v.SaveTokensToFile("tokens.json")
	return nil
}

func (v *Value) SaveTokensToFile(filename string) error {
	v.Factory.Rw.RLock()
	defer v.Factory.Rw.RUnlock()
	data, err := json.MarshalIndent(v.Tokens, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tokens: %v", err)
	}
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write tokens to file: %v", err)
	}
	return nil
}

// AddToken adds a new token to the Tokens struct and updates all maps.
func (v *Value) AddToken(token *Token) {
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()
	v.Tokens = append(v.Tokens, token)
	v.Universe[token.Address] = token
	v.Maps.TokenId[token.TokenId] = token
}

// GetAddress returns the common.Address for a given tokenId.
func (v *Value) GetAddress(tokenId int64) string {
	if tokenId >= 500 {
		return strconv.FormatInt(tokenId, 10)
	}

	v.Factory.Rw.RLock()
	defer v.Factory.Rw.RUnlock()

	token, exists := v.Maps.TokenId[tokenId]
	if !exists {
		log.Printf("Token not found for ID: %d", tokenId)
		return strconv.FormatInt(tokenId, 10)
	}
	return strings.ToLower(token.Address.Hex())
}

// Format formats a string input as a decimal string based on the token's decimals, using address.
func (v *Value) Format(input string, address string) string {
	addr := common.HexToAddress(address)
	v.Factory.Rw.RLock()
	token, exists := v.Universe[addr]
	v.Factory.Rw.RUnlock()
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
