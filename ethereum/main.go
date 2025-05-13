package ethereum

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
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
