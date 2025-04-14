package loopring

import (
	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Txs     []any
}

func Connect(factory *factory.Factory) *Loopring {
	loop := &Loopring{
		Factory: factory,
		Txs:     []any{},
	}
	go loop.Listen()
	go loop.FetchBlocks()
	return loop
}
