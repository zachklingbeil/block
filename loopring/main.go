package loopring

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Block   *Block
}

type Block struct {
	Number       int64      `json:"blockId"`
	Timestamp    int64      `json:"createdAt"`
	Size         int64      `json:"blockSize"`
	Coord        Coordinate `json:"coordinate"`
	Transactions []any      `json:"transactions"`
}

func Connect(factory *factory.Factory) *Loopring {
	loop := &Loopring{
		Factory: factory,
		Block:   &Block{},
	}

	loop.currentBlock()
	return loop
}

func (l *Loopring) Loop() error {
	l.FetchBlock(l.Block.Number)
	l.Coordinates()

	if err := l.Index(); err != nil {
		log.Error("Failed to index block: %v", err)
	}
	l.ProcessTransactions()
	return nil
}
