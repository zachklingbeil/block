package ethereum

import (
	"github.com/ethereum/go-ethereum/params"
	"github.com/zachklingbeil/block/universe"
	"github.com/zachklingbeil/factory"
)

type Ethereum struct {
	Factory        *factory.Factory
	One            *universe.Zero
	Chain          *params.ChainConfig
	Signature      map[string]string
	EventSignature map[string]string
	Header         int64
}

func NewEthereum(factory *factory.Factory, one *universe.Zero) *Ethereum {
	eth := &Ethereum{
		Factory:        factory,
		One:            one,
		Chain:          params.MainnetChainConfig,
		Signature:      make(map[string]string),
		EventSignature: make(map[string]string),
	}
	eth.LoadSignatures()
	go eth.Listen()
	eth.ProcessBlocks(10)
	return eth
}
