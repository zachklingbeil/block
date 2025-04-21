package loopring

import (
	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Block   *Block
}

type Block struct {
	In  *Input
	Out map[Coordinate]*Tx `json:"transactions"`
}

func Connect(factory *factory.Factory) *Loopring {
	loop := &Loopring{
		Factory: factory,
		Block: &Block{
			In:  &Input{},
			Out: make(map[Coordinate]*Tx),
		},
	}
	return loop
}
