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
	var tokensData []Token
	if err := json.Unmarshal(tokens, &tokensData); err != nil {
		log.Fatalf("Failed to unmarshal tokens: %v", err)
	}

	var failed, skipped int

	for _, token := range tokensData {
		if token.Address == "" {
			skipped++
			continue
		}

		tokenJSON, err := json.Marshal(token)
		if err != nil {
			failed++
			continue
		}
		err = factory.Db.Rdb.SAdd(factory.Ctx, "tokens", tokenJSON).Err()
		if err != nil {
			failed++
		}
	}
	total := len(tokensData)
	fmt.Printf("%d tokens\n", total)
}
