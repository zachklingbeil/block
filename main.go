package main

import (
	"github.com/zachklingbeil/block/loopring"
	"github.com/zachklingbeil/block/peer"
	"github.com/zachklingbeil/block/token"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble("timefactory", 1)
	peer.HelloPeers(factory)
	token.NewTokens(factory)

	loop := loopring.Connect(factory)
	// loop.BlockByBlock(55555)
	loop.Loop()
	select {}
}
