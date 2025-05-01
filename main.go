package main

import (
	"github.com/zachklingbeil/block/loopring"
	"github.com/zachklingbeil/block/value"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble()
	// value.NewValue(factory)
	v := value.NewValue(factory)
	loopring.Connect(factory, v)
	// loop.BlockByBlock(55555)
	select {}
}
