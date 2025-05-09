package ethereum

import (
	"github.com/ethereum/go-ethereum/params"
	"github.com/zachklingbeil/block/value"
	"github.com/zachklingbeil/factory"
)

type Ethereum struct {
	Factory        *factory.Factory
	Value          *value.Value
	Chain          *params.ChainConfig
	Signature      map[string]string
	EventSignature map[string]string
	Header         int64
}

func NewEthereum(factory *factory.Factory, value *value.Value) *Ethereum {
	eth := &Ethereum{
		Factory:        factory,
		Value:          value,
		Chain:          params.MainnetChainConfig,
		Signature:      make(map[string]string),
		EventSignature: make(map[string]string),
	}
	eth.LoadSignatures()
	go eth.Listen()
	eth.ProcessBlocks(10)
	return eth
}
