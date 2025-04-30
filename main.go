package main

import (
	"github.com/zachklingbeil/block/value"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble()
	// v := value.NewValue(factory)
	// loopring.Connect(factory, v)
	value.NewValue(factory)
	// v.ProcessBlocks()
	// loop.BlockByBlock(55555)
	select {}
}
