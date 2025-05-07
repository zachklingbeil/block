package input

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"

	"github.com/zachklingbeil/factory"
)

type Signature struct {
	HexSig  string `json:"hex"`
	TextSig string `json:"text"`
}

//go:embed sig.json
var signatures []byte

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
