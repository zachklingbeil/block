package main

import (
	"github.com/zachklingbeil/block/loopring"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble("timefactory")
	// p := peer.HelloPeers(factory)

	go loopring.Connect(factory)
	// go p.HelloUniverse()
	select {}
}
