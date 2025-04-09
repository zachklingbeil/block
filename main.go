package main

import (
	"log"

	"github.com/zachklingbeil/block/loopring"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory, err := factory.NewFactory("block")
	if err != nil {
		log.Fatalf("Error creating factory: %v", err)
	}

	loop := loopring.NewLoopring(factory)
	loop.FetchBlocks()
	loop.QualityControl()
	loop.LoadBlocks()
	select {}
}
