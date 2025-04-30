package main

import (
	"github.com/zachklingbeil/block/loopring"
	"github.com/zachklingbeil/block/value"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble()
	v := value.NewValue(factory)
	loopring.Connect(factory, v)

	// v.ProcessBlocks()
	// loop.BlockByBlock(55555)
	select {}
}

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
