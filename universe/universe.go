package universe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Zero struct {
	Factory *factory.Factory
	One     []*One
	Map     map[string]*One
	Maps    *Maps
	Format  *Format
}

type Format struct{}

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
	LoopringId map[int64]*One
	TokenId    map[int64]*One
	ENS        map[string]*One
	Token      map[string]*One
	ABI        map[string]*One
}

func NewZero(factory *factory.Factory) *Zero {
	z := &Zero{
		Factory: factory,
		Map:     make(map[string]*One),
		Maps: &Maps{
			LoopringId: make(map[int64]*One),
			TokenId:    make(map[int64]*One),
			ENS:        make(map[string]*One),
			Token:      make(map[string]*One),
			ABI:        make(map[string]*One),
		},
		Format: &Format{},
	}
	z.LoadOnes()
	return z
}

func (z *Zero) SyncOnesToRedis(ctx context.Context) error {
	pipe := z.Factory.Data.RB.Pipeline()
	for _, one := range z.One {
		data, err := json.Marshal(one)
		if err != nil {
			return err
		}
		pipe.SAdd(ctx, "one", data)
	}

	_, err := pipe.Exec(ctx)
	return err
}

func (z *Zero) Source(address string) *One {
	z.Factory.Rw.RLock()
	one := z.Map[address]
	defer z.Factory.Rw.RUnlock()
	return one
}

// GetOneByLoopringId returns the *One for a given LoopringId, or an error if not found.
func (z *Zero) LoopringId(id int64) *One {
	z.Factory.Rw.RLock()
	peer := z.Maps.LoopringId[id]
	z.Factory.Rw.RUnlock()
	return peer
}

func (z *Zero) TokenId(id int64) *One {
	if id == 0 {
		return &One{
			Token:    "eth",
			Address:  "0x0000000000000000000000000000000000000000",
			Decimals: 18,
			TokenId:  0,
		}
	}
	z.Factory.Rw.RLock()
	defer z.Factory.Rw.RUnlock()
	token := z.Maps.TokenId[id]
	return token
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
		z.One = append(ones, one)
		z.Map[one.Address] = one
		z.Maps.LoopringId[one.LoopringID] = one
		z.Maps.TokenId[one.TokenId] = one
		z.Maps.ABI[one.ABI] = one
	}
	return nil
}
