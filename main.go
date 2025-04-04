package main

import (
	"log"

	"github.com/zachklingbeil/block/fx"
	"github.com/zachklingbeil/block/loopring"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory, err := factory.NewFactory("block")
	if err != nil {
		log.Fatalf("Error creating factory: %v", err)
	}

	peers, err := fx.HelloUniverse(factory)
	if err != nil {
		log.Fatalf("Error creating peers: %v", err)
	}
	loopring.NewLoopring(factory, peers)
}
