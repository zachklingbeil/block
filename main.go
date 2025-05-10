package main

import (
	"github.com/zachklingbeil/block/loopring"
	"github.com/zachklingbeil/block/universe"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble()
	one := universe.NewZero(factory)

	// go ethereum.NewEthereum(factory, one)
	go loopring.Connect(factory, one)
	select {}
}
