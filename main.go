package main

import (
	"github.com/zachklingbeil/block/loopring"
	"github.com/zachklingbeil/block/universe"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble()
	one := universe.NewZero(factory)
	// eth.New(factory)
	go loopring.Connect(factory, one)
	select {}
}
