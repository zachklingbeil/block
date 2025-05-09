package loopring

import (
	"github.com/zachklingbeil/block/universe"
	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	One     *universe.One
}

func Connect(factory *factory.Factory, one *universe.One) *Loopring {
	loop := &Loopring{
		Factory: factory,
		One:     one,
	}
	loop.CurrentBlock()
	go loop.Listen()
	go loop.Loop()
	return loop
}
