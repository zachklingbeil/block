package main

import (
	"github.com/zachklingbeil/block/circuit"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble()
	// circuit := circuit.NewCircuit(factory)
	circuit.NewCircuit(factory)

	// circuit.Continue()
	select {}
}

// loop := loopring.Connect(factory, circuit)
// loop.Loop()
// loop.BlockByBlock(55555)

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
