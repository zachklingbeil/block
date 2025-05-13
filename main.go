package main

import (
	"github.com/zachklingbeil/block/ethereum"
	"github.com/zachklingbeil/block/universe"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble()
	one := universe.NewZero(factory)

	eth := ethereum.NewEthereum(factory, one)
	eth.Listen()
	// go loopring.Connect(factory, one)
	select {}
}
