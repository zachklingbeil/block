package loop

import (
	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Map     map[Coordinate]*Tx
}

func Connect(factory *factory.Factory) *Loopring {
	loop := &Loopring{
		Factory: factory,
		Map:     make(map[Coordinate]*Tx),
	}
	loop.CreateTxTable()
	go loop.Listen()
	go loop.BlockByBlock()
	return loop
}

func (l *Loopring) BlockByBlock() {
	current := l.currentBlock()
}
