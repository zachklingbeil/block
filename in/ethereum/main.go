package ethereum

import (
	"github.com/zachklingbeil/factory"
)

type Ethereum struct {
	Factory *factory.Factory
}

func NewEthereum(factory *factory.Factory) (*Ethereum, error) {
	return &Ethereum{
		Factory: factory,
	}, nil
}
