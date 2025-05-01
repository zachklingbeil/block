package main

import (
	"github.com/zachklingbeil/block/ethereum"
	"github.com/zachklingbeil/block/loopring"
	"github.com/zachklingbeil/block/value"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble()
	// value.NewValue(factory)

	e := ethereum.NewEthereum(factory)
	e.ProcessBlocks(10)

	v := value.NewValue(factory)
	loopring.Connect(factory, v)
	// loop.BlockByBlock(55555)
	select {}
}
