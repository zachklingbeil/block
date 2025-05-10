package universe

import (
	"encoding/json"
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Zero struct {
	Factory *factory.Factory
	One     []*One
	Map     map[string]*One
	Maps    *Maps
}

type One struct {
	ENS         string `json:"ens,omitempty"`
	LoopringENS string `json:"loopringEns,omitempty"`
	LoopringID  int64  `json:"loopringId,omitempty"`
	Address     string `json:"address"`
	Token       string `json:"token,omitempty"`
	Decimals    int64  `json:"decimals,omitempty"`
	TokenId     int64  `json:"tokenId,omitempty"`
	ABI         string `json:"abi,omitempty"`
}

type Maps struct {
	LoopringId map[int64]string
	TokenId    map[int64]string
}

func NewZero(factory *factory.Factory) *Zero {
	z := &Zero{
		Factory: factory,
		Map:     make(map[string]*One),
		Maps:    &Maps{},
	}
	z.LoadOnes()
	return z
}

func (z *Zero) Source(address string) *One {
	z.Factory.Rw.RLock()
	defer z.Factory.Rw.RUnlock()
	return z.Map[address]
}

// GetOneByLoopringId returns the *One for a given LoopringId, or an error if not found.
func (z *Zero) LoopringId(id int64) *One {
	z.Factory.Rw.RLock()
	defer z.Factory.Rw.RUnlock()
	addr, ok := z.Maps.LoopringId[id]
	if !ok || addr == "" {
		// Return a minimal One with only LoopringID set
		return &One{LoopringID: id}
	}
	one, ok := z.Map[addr]
	if !ok {
		return &One{LoopringID: id}
	}
	return one
}

func (z *Zero) TokenId(id int64) *One {
	z.Factory.Rw.RLock()
	defer z.Factory.Rw.RUnlock()
	addr, ok := z.Maps.TokenId[id]
	if !ok || addr == "" {
		// Return a minimal One with only TokenId set
		return &One{TokenId: id}
	}
	one, ok := z.Map[addr]
	if !ok {
		return &One{TokenId: id}
	}
	return one
}

func (z *Zero) LoadOnes() error {
	source, err := z.Factory.Data.RB.SMembers(z.Factory.Ctx, "one").Result()
	if err != nil {
		return fmt.Errorf("failed to fetch peers from Redis set: %v", err)
	}

	ones := make([]*One, 0, len(source))
	for _, s := range source {
		one := &One{}
		if err := json.Unmarshal([]byte(s), one); err != nil {
			return fmt.Errorf("failed to unmarshal One: %v", err)
		}
		ones = append(ones, one)
	}

	z.Factory.Rw.Lock()
	defer z.Factory.Rw.Unlock()

	z.One = ones
	z.Map = make(map[string]*One, len(ones))
	z.Maps.LoopringId = make(map[int64]string, len(ones))
	z.Maps.TokenId = make(map[int64]string, len(ones))
	for _, one := range ones {
		z.Map[one.Address] = one
		z.Maps.LoopringId[one.LoopringID] = one.Address
		z.Maps.TokenId[one.TokenId] = one.Address
	}
	return nil
}
