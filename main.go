package main

import (
	"log"

	"github.com/zachklingbeil/block/ethereum"
	"github.com/zachklingbeil/block/loopring"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory, err := factory.NewFactory()
	if err != nil {
		log.Fatalf("Error creating factory: %v", err)
	}

	loopring := loopring.NewLoopring(factory)
	ethereum := ethereum.NewEthereum(factory)
	log.Println("Loopring and Ethereum initialized successfully:", loopring, ethereum)
}
