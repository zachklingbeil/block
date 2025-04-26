package token

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/zachklingbeil/factory"
)

//go:embed tokens.json
var tokens []byte

type TokenIn struct {
	Symbol     string `json:"symbol"`
	Address    string `json:"address"`
	Token      int64  `json:"tokenId"`
	LoopringID string `json:"accountId,omitempty"`
	Decimals   int    `json:"decimals"`
}
type Token struct {
	Symbol     string `json:"symbol"`
	Address    string `json:"address"`
	Token      int64  `json:"tokenId"`
	LoopringID int64  `json:"accountId,omitempty"`
	Decimals   int64  `json:"decimals"`
}

func NewTokens(factory *factory.Factory) {
	var in []TokenIn
	if err := json.Unmarshal(tokens, &in); err != nil {
		log.Fatalf("Failed to unmarshal tokens: %v", err)
	}

	// Separate maps for LoopringID keys and TokenID keys
	loopringIDMap := make(map[int64]Token)
	tokenIDMap := make(map[int64]Token)

	for _, tokenIn := range in {
		// Convert TokenIn to Token
		token := Token{
			Symbol:   tokenIn.Symbol,
			Address:  tokenIn.Address,
			Token:    tokenIn.Token,
			Decimals: int64(tokenIn.Decimals),
		}

		// Convert LoopringID (string) to int64 if valid
		if tokenIn.LoopringID != "" {
			loopringID, err := strconv.ParseInt(tokenIn.LoopringID, 10, 64)
			if err != nil {
				log.Printf("Failed to parse LoopringID '%s' for token '%s': %v", tokenIn.LoopringID, tokenIn.Symbol, err)
			} else {
				token.LoopringID = loopringID
				// Add to LoopringID map
				loopringIDMap[loopringID] = token
			}
		}

		// Add to TokenID map
		tokenIDMap[token.Token] = token
	}

	// Store LoopringID map in Redis
	for key, token := range loopringIDMap {
		tokenJSON, err := json.Marshal(token)
		if err != nil {
			log.Printf("Failed to marshal token for LoopringID key %d: %v", key, err)
			continue
		}

		// Use HSet to store the key-value pair in the Redis hash "loopringIDMap"
		err = factory.Redis.HSet(factory.Ctx, "loopringIDMap", fmt.Sprintf("%d", key), tokenJSON).Err()
		if err != nil {
			log.Printf("Failed to store token in Redis for LoopringID key %d: %v", key, err)
		}
	}

	// Store TokenID map in Redis
	for key, token := range tokenIDMap {
		tokenJSON, err := json.Marshal(token)
		if err != nil {
			log.Printf("Failed to marshal token for TokenID key %d: %v", key, err)
			continue
		}

		// Use HSet to store the key-value pair in the Redis hash "tokenIDMap"
		err = factory.Redis.HSet(factory.Ctx, "tokenIDMap", fmt.Sprintf("%d", key), tokenJSON).Err()
		if err != nil {
			log.Printf("Failed to store token in Redis for TokenID key %d: %v", key, err)
		}
	}

	fmt.Printf("%d tokens processed\n", len(in))
}
