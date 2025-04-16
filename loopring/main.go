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
		Txs:     make([]any, 0, 10000),
	}
	go loop.Listen()
	go loop.FetchBlocks()
	return loop
}
