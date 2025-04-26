package main

import (
	"github.com/zachklingbeil/block/token"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble("timefactory", 0)
	// circuit := circuit.NewCircuit(factory)

	token.NewTokens(factory)

	// circuit.Continue()

	// factory.Json.Print(circuit.Get("zachklingbeil.eth"))
	// factory.Json.Print(circuit.Get("35773"))
	// factory.Json.Print(circuit.Get("eTh"))
	// factory.Json.Print(circuit.Get("Lrc"))
	// factory.Json.Print(circuit.Get("friction"))
	// factory.Json.Print(circuit.Get(0))
	// factory.Json.Print(circuit.Get(1))
	// factory.Json.Print(circuit.Get("0")) // empty
	// factory.Json.Print(circuit.Get("1"))
	// factory.Json.Print(circuit.Get("30"))
	// factory.Json.Print(circuit.Get(30))

	select {}
	// loop := loopring.Connect(factory, circuit)
	// loop.Loop()
	// loop.BlockByBlock(55555)

}
