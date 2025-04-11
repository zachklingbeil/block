package loopring

import (
	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Txs     []any
	Map     map[string]*Tx
}
type Tx struct {
	Zero     string `json:"zero"`
	One      string `json:"one"`
	Token    string `json:"token"`
	Value    int64  `json:"value"`
	Fee      int64  `json:"fee"`
	FeeToken string `json:"feeToken"`
}

func Connect(factory *factory.Factory) *Loopring {
	loop := &Loopring{
		Factory: factory,
		Txs:     []any{},
		Map:     make(map[string]*Tx),
	}
	go loop.Listen()
	go loop.FetchBlocks()
	return loop
}
