package main

import (
	"github.com/zachklingbeil/block/circuit"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble("timefactory", 0)
	circuit := circuit.NewCircuit(factory)

	circuit.Get("zachklingbeil.eth")
	circuit.Get("eTH")
	circuit.Get("LRC")
	circuit.Get(35773)
	circuit.Get(69)
	circuit.Get(-1)
	circuit.Get("friction")
	select {}
	// loop := loopring.Connect(factory, circuit)
	// loop.Loop()
	// loop.BlockByBlock(55555)
}
