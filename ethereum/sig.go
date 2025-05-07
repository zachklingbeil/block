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

func (e *Ethereum) LoadHexToText() error {
	source, err := e.Factory.Data.RB.SMembers(e.Factory.Ctx, "signature").Result()
	if err != nil {
		return fmt.Errorf("failed to fetch signatures from Redis set: %v", err)
	}
	m := make(map[string]string)
	for _, sigJSON := range source {
		var sig Signature
		if err := json.Unmarshal([]byte(sigJSON), &sig); err != nil {
			log.Printf("Skipping invalid signature: %v (data: %s)", err, sigJSON)
			continue
		}
		m[sig.Hex] = sig.Text
	}
	e.Factory.Rw.Lock()
	e.HexToText = m
	e.Factory.Rw.Unlock()
	fmt.Printf("%d HexToText\n", len(e.HexToText))
	return nil
}

// Concurrent read method
func (e *Ethereum) GetHexText(hex string) (string, bool) {
	e.Factory.Rw.RLock()
	defer e.Factory.Rw.RUnlock()
	val, ok := e.HexToText[hex]
	return val, ok
}
