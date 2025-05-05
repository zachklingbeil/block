package loopring

import (
	"github.com/zachklingbeil/block/universe"
	"github.com/zachklingbeil/block/value"
	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory  *factory.Factory
	Value    *value.Value
	Universe *universe.Universe
}

func Connect(factory *factory.Factory, value *value.Value, universe *universe.Universe) *Loopring {
	loop := &Loopring{
		Factory:  factory,
		Value:    value,
		Universe: universe,
	}
	loop.CurrentBlock()
	go loop.Listen()
	go loop.Loop()
	return loop
}
