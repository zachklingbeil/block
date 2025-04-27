package manual

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
	Token    string `json:"token,omitempty"`
	Address  string `json:"address,omitempty"`
	TokenId  string `json:"tokenId,omitempty"`
	TokenInt int64  `json:"tokenInt,omitempty"`
	Decimals int64  `json:"decimals,omitempty"`
}

func NewTokens(factory *factory.Factory) {
	var tokensData []Token
	if err := json.Unmarshal(tokens, &tokensData); err != nil {
		log.Fatalf("Failed to unmarshal tokens: %v", err)
	}
	for _, token := range tokensData {
		tokenJSON, err := json.Marshal(token)
		if err != nil {
			log.Printf("Failed to marshal token: %v", err)
			continue
		}

		if err := factory.Data.RB.SAdd(factory.Ctx, "token", tokenJSON).Err(); err != nil {
			log.Printf("Failed to add token to Redis: %v", err)
		}
	}
	fmt.Printf("%d tokens\n", len(tokensData))
}
