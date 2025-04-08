package main

import (
	"log"

	"github.com/zachklingbeil/factory"
)

func main() {
	factory, err := factory.NewFactory("block")
	if err != nil {
		log.Fatalf("Error creating factory: %v", err)
	}

	factory.Peer.HelloUniverse()

	// loop := loopring.NewLoopring(factory)
	// loop.FetchBlocks()
	// select {}
}
