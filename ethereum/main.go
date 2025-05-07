package ethereum

import (
	"github.com/ethereum/go-ethereum/params"
	"github.com/zachklingbeil/block/value"
	"github.com/zachklingbeil/factory"
)

type Ethereum struct {
	Factory   *factory.Factory
	Value     *value.Value
	Chain     *params.ChainConfig
	HexToText map[string]string
}

func NewEthereum(factory *factory.Factory, value *value.Value) *Ethereum {
	eth := &Ethereum{
		Factory: factory,
		Value:   value,
		Chain:   params.MainnetChainConfig,
	}
	eth.LoadHexToText()
	return eth
}
