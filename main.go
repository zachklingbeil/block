package main

import (
	"github.com/zachklingbeil/block/input"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble()

	input.NewSignatures(factory)
	// value.NewValue(factory)
	// v := value.NewValue(factory)
	// e := ethereum.NewEthereum(factory, v)
	// go e.ProcessBlocks(10)
	// go loopring.Connect(factory, v)
	select {}
}
