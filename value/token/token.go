package token

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/zachklingbeil/factory"
)

type Tokens struct {
	Factory *factory.Factory
	Tokens  []*Token
	Map     map[any]*Token
}

func NewTokens(factory *factory.Factory) *Tokens {
	t := &Tokens{
		Factory: factory,
		Tokens:  make([]*Token, 0),
		Map:     make(map[any]*Token),
	}
	t.LoadTokens()
	return t
}

type Token struct {
	Token    string `json:"token,omitempty"`
	Address  string `json:"address,omitempty"`
	Decimals string `json:"decimals,omitempty"`
	TokenId  string `json:"tokenId,omitempty"`
	TokenInt int64  `json:"tokenInt,omitempty"`
}

func (t *Tokens) LoadTokens() error {
	source, err := t.Factory.Data.RB.HGetAll(t.Factory.Ctx, "token").Result()
	if err != nil {
		return fmt.Errorf("failed to fetch tokens from Redis hash: %v", err)
	}
	t.Tokens = make([]*Token, 0, len(source))
	t.Map = make(map[any]*Token, len(source))
	for _, tokenJSON := range source {
		var token Token
		if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
			log.Printf("Skipping invalid token: %v (data: %s)", err, tokenJSON)
			continue
		}
		t.Tokens = append(t.Tokens, &token)
		t.Map[token.TokenInt] = &token
	}
	return nil
}

// Int64 tokenId in and desired return:
// 0 - Token 1 - Address 2 - Decimals 3 - TokenId
func (t *Tokens) Get(tokenId int64, field uint8) string {
	if tokenId >= 500 {
		return strconv.FormatInt(tokenId, 10)
	}

	t.Factory.Rw.RLock()
	defer t.Factory.Rw.RUnlock()

	token, exists := t.Map[tokenId]
	if !exists {
		log.Printf("Token not found for ID: %d", tokenId)
		return ""
	}

	switch field {
	case 0:
		return token.Token
	case 1:
		return token.Address
	case 2:
		return token.Decimals
	case 3:
		return token.TokenId
	default:
		log.Printf("Invalid field selector: %d", field)
		return ""
	}
}

// FormatValue formats a string input as a decimal string based on the token's decimals.
func (t *Tokens) Format(input string, key any) string {
	t.Factory.Rw.RLock()
	token, exists := t.Map[key]
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
