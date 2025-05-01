package main

import (
	"github.com/zachklingbeil/block/loopring"
	"github.com/zachklingbeil/block/value"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble()
	v := value.NewValue(factory)
	loop := loopring.Connect(factory, v)
	loop.BlockByBlock(52222)
	// e := ethereum.NewEthereum(factory)
	// e.ProcessBlocks(10)

	select {}
}
