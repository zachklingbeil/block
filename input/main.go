package input

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zachklingbeil/factory"
)

type Token struct {
	Token    string         `json:"token,omitempty"`
	Address  common.Address `json:"address,omitempty"`
	Decimals int64          `json:"decimals,omitempty"`
	TokenId  int64          `json:"tokenId,omitempty"`
	ABI      string         `json:"abi,omitempty"`
}

type rawToken struct {
	Token    string          `json:"token,omitempty"`
	Address  common.Address  `json:"address,omitempty"`
	Decimals json.RawMessage `json:"decimals,omitempty"`
	TokenId  int64           `json:"tokenId,omitempty"`
	ABI      string          `json:"abi,omitempty"`
}

type Signature struct {
	ID      int    `json:"id"`
	TextSig string `json:"text"`
	HexSig  string `json:"hex"`
}

//go:embed token.json
var tokens []byte

//go:embed signatures.json
var signatures []byte

//go:embed event_signatures.json
var event_signatures []byte

func NewTokens(factory *factory.Factory) {
	var rawTokens []rawToken
	if err := json.Unmarshal(tokens, &rawTokens); err != nil {
		log.Fatalf("Failed to unmarshal tokens: %v", err)
	}
	var tokensData []Token
	for _, rt := range rawTokens {
		var dec int64
		// Try to unmarshal as int
		if err := json.Unmarshal(rt.Decimals, &dec); err != nil {
			// Try to unmarshal as string
			var decStr string
			if err := json.Unmarshal(rt.Decimals, &decStr); err == nil {
				parsed, err := strconv.ParseInt(decStr, 10, 64)
				if err == nil {
					dec = parsed
				}
			}
		}
		tokensData = append(tokensData, Token{
			Token:    rt.Token,
			Address:  rt.Address,
			Decimals: dec,
			TokenId:  rt.TokenId,
			ABI:      rt.ABI,
		})
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

func NewSignatures(factory *factory.Factory) {
	var signaturesData []Signature
	if err := json.Unmarshal(signatures, &signaturesData); err != nil {
		log.Fatalf("Failed to unmarshal signatures: %v", err)
	}
	for _, sig := range signaturesData {
		sigJSON, err := json.Marshal(sig)
		if err != nil {
			log.Printf("Failed to marshal signature: %v", err)
			continue
		}

		if err := factory.Data.RB.SAdd(factory.Ctx, "signature", sigJSON).Err(); err != nil {
			log.Printf("Failed to add signature to Redis: %v", err)
		}
	}
	fmt.Printf("%d signatures\n", len(signaturesData))
}
func NewEventSignatures(factory *factory.Factory) {
	var signaturesData []Signature
	if err := json.Unmarshal(event_signatures, &signaturesData); err != nil {
		log.Fatalf("Failed to unmarshal event signatures: %v", err)
	}
	for _, sig := range signaturesData {
		sigJSON, err := json.Marshal(sig)
		if err != nil {
			log.Printf("Failed to marshal event signature: %v", err)
			continue
		}

		if err := factory.Data.RB.SAdd(factory.Ctx, "event_signature", sigJSON).Err(); err != nil {
			log.Printf("Failed to add event signature to Redis: %v", err)
		}
	}
	fmt.Printf("%d event signatures\n", len(signaturesData))
}
