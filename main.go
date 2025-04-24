package main

import (
	"github.com/zachklingbeil/block/circuit"
	"github.com/zachklingbeil/block/loopring"
	"github.com/zachklingbeil/block/peer"
	"github.com/zachklingbeil/block/token"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble("timefactory", 1)
	circuit := circuit.NewCircuit(factory)
	loop := loopring.Connect(factory, circuit)

	peer.HelloPeers(factory)
	token.NewTokens(factory)
	loop.Loop()
	// loop.BlockByBlock(55555)
	select {}
}
