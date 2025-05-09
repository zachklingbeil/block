package ethereum

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

type Signature struct {
	Hex  string `json:"hex"`
	Text string `json:"text"`
}

// Signer returns a signer for Ethereum mainnet at the given block number and time.
func (e *Ethereum) Signer(blockNumber *big.Int, blockTime uint64) types.Signer {
	return types.MakeSigner(e.Chain, blockNumber, blockTime)
}

// LoadSignatures loads both "signature" and "event_signature" sets from Redis into the struct maps.
func (e *Ethereum) LoadSignatures() error {
	// Helper function to load a Redis set into a map
	loadSet := func(setName string) (map[string]string, error) {
		source, err := e.Factory.Data.RB.SMembers(e.Factory.Ctx, setName).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch %s from Redis set: %v", setName, err)
		}
		m := make(map[string]string)
		for _, sigJSON := range source {
			var sig Signature
			if err := json.Unmarshal([]byte(sigJSON), &sig); err != nil {
				log.Printf("Skipping invalid %s: %v (data: %s)", setName, err, sigJSON)
				continue
			}
			m[sig.Hex] = sig.Text
		}
		return m, nil
	}

	sigMap, err := loadSet("signature")
	if err != nil {
		return err
	}
	eventSigMap, err := loadSet("event_signature")
	if err != nil {
		return err
	}

	e.Factory.Rw.Lock()
	e.Signature = sigMap
	e.EventSignature = eventSigMap
	e.Factory.Rw.Unlock()

	fmt.Printf("%d Signature, %d EventSignature loaded\n", len(e.Signature), len(e.EventSignature))
	return nil
}

// Concurrent read methods
func (e *Ethereum) GetSignature(hex string) (string, bool) {
	e.Factory.Rw.RLock()
	defer e.Factory.Rw.RUnlock()
	val, ok := e.Signature[hex]
	return val, ok
}

func (e *Ethereum) GetEventSignature(hex string) (string, bool) {
	e.Factory.Rw.RLock()
	defer e.Factory.Rw.RUnlock()
	val, ok := e.EventSignature[hex]
	return val, ok
}
