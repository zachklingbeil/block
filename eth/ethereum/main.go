package ethereum

import (
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
	Block          *types.Block
	Signature      map[string]string
	EventSignature map[string]string
	ABIs           map[string]abi.ABI // address → ABI
	Header         *big.Int
}

func NewEthereum(factory *factory.Factory, zero *universe.Zero) *Ethereum {
	eth := &Ethereum{
		Factory:   factory,
		Zero:      zero,
		Chain:     params.MainnetChainConfig,
		ABIs:      make(map[string]abi.ABI),
		Signature: make(map[string]string),
	}
	eth.LoadSignatures()
	eth.PopulateABIs()
	return eth
}

// PopulateABIs loads and parses all contract ABIs into the ABIs map.
func (e *Ethereum) PopulateABIs() {
	for addr, abiJSON := range e.Zero.Maps.ABI {
		if abiJSON.ABI == "" || abiJSON.ABI == "." {
			continue
		}
		parsedABI, err := abi.JSON(strings.NewReader(abiJSON.ABI))
		if err != nil {
			log.Printf("Failed to parse ABI for %s: %v", addr, err)
			continue
		}
		e.ABIs[addr] = parsedABI
	}
}

// Signer returns a signer for Ethereum mainnet at the given block number and time.
func (e *Ethereum) Signer(blockNumber *big.Int, blockTime uint64) types.Signer {
	return types.MakeSigner(e.Chain, blockNumber, blockTime)
}
