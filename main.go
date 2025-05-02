package main

import (
	"github.com/zachklingbeil/block/ethereum"
	"github.com/zachklingbeil/block/value"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble()
	v := value.NewValue(factory)
	e := ethereum.NewEthereum(factory, v)
	e.ProcessBlocks(10)
	// loopring.Connect(factory, v)
	select {}
}
