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

	io := &Input{
		Loopring: loopring,
		Factory:  factory,
	}
	loopring.FetchBlocks()
	return io
}
