package token

import (
	_ "embed"
	"encoding/json"
)

//go:embed tokens.json
var tokens []byte

type Token struct {
	Symbol   string `json:"symbol,omitempty"`
	Address  string `json:"address,omitempty"`
	TokenId  int    `json:"tokenId,omitempty"`
	Decimals int    `json:"decimals,omitempty"`
	Zero     int    `json:"accountId,omitempty"`
}

func NewTokens() []*Token {
	var source []*Token
	if err := json.Unmarshal(tokens, &source); err != nil {
		return nil
	}
	return source
}
