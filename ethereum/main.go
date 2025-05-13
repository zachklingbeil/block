package ethereum

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/zachklingbeil/block/universe"
	"github.com/zachklingbeil/factory"
)

type Ethereum struct {
	Factory        *factory.Factory
	Zero           *universe.Zero
	Chain          *params.ChainConfig
	Signature      map[string]string
	EventSignature map[string]string
	EventABI       map[string]abi.Event
	Header         *big.Int
}

func NewEthereum(factory *factory.Factory, zero *universe.Zero) *Ethereum {
	eth := &Ethereum{
		Factory:        factory,
		Zero:           zero,
		Chain:          params.MainnetChainConfig,
		Signature:      make(map[string]string),
		EventSignature: make(map[string]string),
		EventABI:       make(map[string]abi.Event),
	}
	eth.LoadSignatures()
	eth.PopulateEventABI()
	// eth.BlockByBlock()
	return eth
}

func (e *Ethereum) PopulateEventABI() {
	for addr, abiJSON := range e.Zero.Maps.ABI {
		if abiJSON.ABI == "" || abiJSON.ABI == "." {
			continue
		}
		parsedABI, err := abi.JSON(strings.NewReader(abiJSON.ABI))
		if err != nil {
			log.Printf("Failed to parse ABI for %s: %v", addr, err)
			continue
		}
		for _, event := range parsedABI.Events {
			sigHash := event.ID.Hex()
			e.EventABI[sigHash] = event
			e.EventSignature[sigHash] = event.String()
		}
	}
}

type Signature struct {
	Hex  string `json:"hex"`
	Text string `json:"text"`
}

// Signer returns a signer for Ethereum mainnet at the given block number and time.
func (e *Ethereum) Signer(blockNumber *big.Int, blockTime uint64) types.Signer {
	return types.MakeSigner(e.Chain, blockNumber, blockTime)
}

func (e *Ethereum) LoadSignatureSet() error {
	source, err := e.Factory.Data.RB.SMembers(e.Factory.Ctx, "signature").Result()
	if err != nil {
		return fmt.Errorf("failed to fetch signature from Redis set: %v", err)
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
	e.Signature = m
	e.Factory.Rw.Unlock()
	return nil
}

// LoadEventSignatureSet loads the "event_signature" set from Redis into the EventSignature map.
func (e *Ethereum) LoadEventSignatureSet() error {
	source, err := e.Factory.Data.RB.SMembers(e.Factory.Ctx, "event_signature").Result()
	if err != nil {
		return fmt.Errorf("failed to fetch event_signature from Redis set: %v", err)
	}
	m := make(map[string]string)
	for _, sigJSON := range source {
		var sig Signature
		if err := json.Unmarshal([]byte(sigJSON), &sig); err != nil {
			log.Printf("Skipping invalid event_signature: %v (data: %s)", err, sigJSON)
			continue
		}
		m[sig.Hex] = sig.Text
	}
	e.Factory.Rw.Lock()
	e.EventSignature = m
	e.Factory.Rw.Unlock()
	return nil
}
func (e *Ethereum) LoadSignatures() error {
	if err := e.LoadSignatureSet(); err != nil {
		return err
	}
	if err := e.LoadEventSignatureSet(); err != nil {
		return err
	}
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
