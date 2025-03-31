package main

import (
	"log"

	"github.com/zachklingbeil/block/loopring"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory, err := factory.NewFactory()
	if err != nil {
		log.Fatalf("Error creating factory: %v", err)
	}

	loopring := loopring.NewLoopring(factory)
	loopring.CurrentBlock()
	loopring.GetBlock(10000)
	loopring.GetBlock(10001)
	loopring.GetBlock(10002)
}
