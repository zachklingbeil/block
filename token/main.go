package token

import (
	_ "embed"
	"encoding/json"
)

//go:embed tokens.json
var tokens []byte

type Tokens struct {
	Map map[string]*Token
}

type Token struct {
	Symbol   string `json:"symbol"`
	Address  string `json:"address"`
	ID       int    `json:"id"`
	Decimals int    `json:"decimals"`
}

func NewTokens() *Tokens {
	t := &Tokens{
		Map: make(map[string]*Token),
	}

	var source []*Token
	if err := json.Unmarshal(tokens, &source); err != nil {
		return nil
	}

	for _, token := range source {
		t.Map[token.Address] = token
		t.Map[token.Symbol] = token
	}
	return t
}
