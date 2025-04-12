package loopring

import (
	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Txs     []any
	Map     map[string]*Tx
}

func Connect(factory *factory.Factory) *Loopring {
	loop := &Loopring{
		Factory: factory,
		Txs:     []any{},
		Map:     make(map[string]*Tx),
	}
	go loop.Listen()
	go loop.FetchBlocks()
	return loop
}
