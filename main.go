package main

import (
	"time"

	"github.com/zachklingbeil/block/loopring"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble("timefactory", 10*time.Second)
	// peer.HelloPeers(factory)

	// token.NewTokens(factory)
	loop := loopring.Connect(factory)
	loop.Loop()
	select {}
}
