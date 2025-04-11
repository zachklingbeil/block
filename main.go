package main

import (
	"github.com/zachklingbeil/block/loopring"
	"github.com/zachklingbeil/factory"
	"github.com/zachklingbeil/peer"
)

func main() {
	factory := factory.Assemble("timefactory")
	peer.HelloPeers(factory)
	loopring.Connect(factory)
	select {}
}
