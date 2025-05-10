package universe

import (
	"github.com/zachklingbeil/factory"
)

type Zero struct {
	Factory *factory.Factory
	One     []*One
	Map     map[string]*One
	Maps    *Maps
}

type Maps struct {
	LoopringId map[int64]string
	TokenId    map[int64]string
}

func NewZero(factory *factory.Factory) *Zero {
	v := &Zero{
		Factory: factory,
		One:     make([]*One, 0),
		Map:     make(map[string]*One),
		Maps: &Maps{
			LoopringId: make(map[int64]string),
			TokenId:    make(map[int64]string),
		},
	}

	// v.LoadTokens()
	// v.LoadPeers()

	// for _, p := range v.Peers {
	// 	v.Map[p.Address] = p
	// }

	// for _, t := range v.Tokens {
	// 	v.Map[t.Address] = t
	// }
	// fmt.Printf("Map: %d\n", len(v.Map))
	return v
}

func (v *Zero) Source(address string) any {
	return v.Map[address]
}
