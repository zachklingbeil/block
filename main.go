package main

import (
	"github.com/zachklingbeil/block/ethereum"
	"github.com/zachklingbeil/block/value"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble()

	// input.NewSignatures(factory)
	// input.NewEventSignatures(factory)
	// value.NewValue(factory)
	v := value.NewValue(factory)
	go ethereum.NewEthereum(factory, v)
	// go loopring.Connect(factory, v)
	select {}
}
