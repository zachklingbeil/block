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
	Symbol    string `json:"symbol,omitempty"`
	Address   string `json:"address,omitempty"`
	TokenId   int64  `json:"tokenId,omitempty"`
	Decimals  int    `json:"decimals,omitempty"`
	AccountID int64  `json:"accountId,omitempty"`
}

func NewTokens(factory *factory.Factory, circuit *circuit.Circuit) {
	var in []Token
	if err := json.Unmarshal(tokens, &in); err != nil {
		log.Fatalf("Failed to unmarshal tokens: %v", err)
	}

	for _, token := range in {
		circuit.AddString(token.Address, token)
	}
	fmt.Printf("%d tokens\n", len(in))
}
