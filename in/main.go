package input

import (
	"log"

	"github.com/zachklingbeil/block/in/ethereum"
	"github.com/zachklingbeil/block/in/loopring"
	"github.com/zachklingbeil/factory"
)

type Input struct {
	Ethereum *ethereum.Ethereum
	Loopring *loopring.Loopring
	Factory  *factory.Factory
}

func NewInput(factory *factory.Factory) *Input {
	loopring, err := loopring.NewLoopring(factory)
	if err != nil {
		log.Fatalf("Error creating Loopring instance: %v", err)
	}

	ethereum, err := ethereum.NewEthereum(factory)
	if err != nil {
		log.Fatalf("Error creating Ethereum instance: %v", err)
	}
	loopring.FetchBlocks()
	loopring.EnsureTransactions()
	return &Input{
		Ethereum: ethereum,
		Loopring: loopring,
		Factory:  factory,
	}
}
