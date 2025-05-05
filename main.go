package main

import (
	"github.com/zachklingbeil/block/peer"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory := factory.Assemble()

	peer.NewPeers(factory)
	// value.NewValue(factory)
	// v := value.NewValue(factory)
	// // // e := ethereum.NewEthereum(factory, v)

	// // // go e.ProcessBlocks(10)
	// go loopring.Connect(factory, v)
	select {}
}
