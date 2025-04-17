package main

import (
	"github.com/zachklingbeil/block/loop"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble("timefactory")
	// p := peer.HelloPeers(factory)

	go loop.Connect(factory)
	// go p.HelloUniverse()

	// loop.Connect(factory)
	select {}
}
