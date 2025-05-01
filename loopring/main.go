package loopring

import (
	"github.com/zachklingbeil/block/value"
	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Value   *value.Value
}

func Connect(factory *factory.Factory, value *value.Value) *Loopring {
	loop := &Loopring{
		Factory: factory,
		Value:   value,
	}
	loop.CurrentBlock()
	go loop.Listen()
	go loop.Loop()
	return loop
}
