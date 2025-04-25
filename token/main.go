package token

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"

	"github.com/zachklingbeil/block/circuit"
	"github.com/zachklingbeil/factory"
)

//go:embed tokens.json
var tokens []byte

type Token struct {
	Address   string `json:"address,omitempty"`
	AccountID int64  `json:"accountId,omitempty"`
	Symbol    string `json:"symbol,omitempty"`
	TokenId   int64  `json:"tokenId,omitempty"`
	Decimals  int    `json:"decimals,omitempty"`
}

func NewTokens(factory *factory.Factory, circuit *circuit.Circuit) {
	var in []Token
	if err := json.Unmarshal(tokens, &in); err != nil {
		log.Fatalf("Failed to unmarshal tokens: %v", err)
	}
	for _, token := range in {
		tokenJSON, _ := json.Marshal(token)
		factory.Redis.SAdd(factory.Ctx, "tokens", tokenJSON).Err()
	}
	for _, token := range in {
		circuit.AddString(token.Symbol, token)
		circuit.AddInt(token.AccountID, token)
		circuit.AddInt(token.TokenId, token)
	}
	fmt.Printf("%d tokens\n", len(in))
}
