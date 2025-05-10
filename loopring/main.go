package loopring

import (
	"github.com/zachklingbeil/block/universe"
	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Zero    *universe.Zero
}

func Connect(factory *factory.Factory, zero *universe.Zero) *Loopring {
	loop := &Loopring{
		Factory: factory,
		Zero:    zero,
	}
	loop.CurrentBlock()
	go loop.Listen()
	go loop.Loop()
	return loop
}
