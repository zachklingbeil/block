package ethereum

import (
	"github.com/zachklingbeil/factory"
)

type Ethereum struct {
	Factory *factory.Factory
	Map     map[Coordinate]*Tx
}

func NewEthereum(factory *factory.Factory) *Ethereum {
	return &Ethereum{
		Factory: factory,
	}
}
