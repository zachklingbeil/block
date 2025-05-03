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
	// e := ethereum.NewEthereum(factory, v)

	// go e.ProcessBlocks(10)
	go loopring.Connect(factory, v)
	select {}
}
