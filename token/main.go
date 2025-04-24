package token

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"

	"github.com/zachklingbeil/factory"
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

func NewTokens(factory *factory.Factory) {
	var in []Token
	if err := json.Unmarshal(tokens, &in); err != nil {
		log.Fatalf("Failed to unmarshal tokens: %v", err)
	}

	for _, token := range in {
		tokenJSON, _ := json.Marshal(token)
		factory.Redis.SAdd(factory.Ctx, "tokens", tokenJSON).Err()
	}
	fmt.Printf("%d tokens\n", (len(in)))
}
