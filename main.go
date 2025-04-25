package main

import (
	"github.com/zachklingbeil/block/circuit"
	"github.com/zachklingbeil/block/peer"
	"github.com/zachklingbeil/block/token"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble("timefactory", 0)
	circuit := circuit.NewCircuit(factory)

	token.NewTokens(factory, circuit)
	peer.HelloPeers(factory, circuit)

	select {}
	// loop := loopring.Connect(factory, circuit)
	// loop.Loop()
	// loop.BlockByBlock(55555)
}
