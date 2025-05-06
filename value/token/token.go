package token

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zachklingbeil/factory"
)

type Tokens struct {
	Factory *factory.Factory
	Tokens  []*Token
	Map     map[common.Address]*Token // map by Address
	IdMap   map[int64]*Token          // map by TokenId
}

func NewTokens(factory *factory.Factory) *Tokens {
	t := &Tokens{
		Factory: factory,
	}
	t.LoadTokens()
	return t
}

type Token struct {
	Token    string         `json:"token,omitempty"`
	Address  common.Address `json:"address,omitempty"`
	Decimals string         `json:"decimals,omitempty"`
	TokenId  int64          `json:"tokenId,omitempty"`
}

func (t *Tokens) LoadTokens() error {
	source, err := t.Factory.Data.RB.SMembers(t.Factory.Ctx, "token").Result()
	if err != nil {
		return fmt.Errorf("failed to fetch tokens from Redis set: %v", err)
	}
	t.Tokens = make([]*Token, 0, len(source))
	t.Map = make(map[common.Address]*Token, len(source))
	t.IdMap = make(map[int64]*Token, len(source))
	for _, tokenJSON := range source {
		var token Token
		if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
			log.Printf("Skipping invalid token: %v (data: %s)", err, tokenJSON)
			continue
		}
		t.Tokens = append(t.Tokens, &token)
		t.Map[token.Address] = &token
		t.IdMap[token.TokenId] = &token
	}
	fmt.Printf("%d tokens loaded\n", len(t.Tokens))
	return nil
}

// GetAddress returns the common.Address for a given tokenId.
func (t *Tokens) GetAddress(tokenId int64) string {
	if tokenId >= 500 {
		return strconv.FormatInt(tokenId, 10)
	}

	t.Factory.Rw.RLock()
	defer t.Factory.Rw.RUnlock()

	token, exists := t.IdMap[tokenId]
	if !exists {
		log.Printf("Token not found for ID: %d", tokenId)
		return strconv.FormatInt(tokenId, 10)
	}
	return strings.ToLower(token.Address.Hex())
}

// Format formats a string input as a decimal string based on the token's decimals, using address.
func (t *Tokens) Format(input string, address string) string {
	addr := common.HexToAddress(address)
	t.Factory.Rw.RLock()
	token, exists := t.Map[addr]
	t.Factory.Rw.RUnlock()
	if !exists {
		return input
	}
	return format(input, token)
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
