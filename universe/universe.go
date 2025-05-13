package universe

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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

func (z *Zero) Source(hex string) *One {
	address := strings.ToLower(hex)

	z.Factory.Rw.RLock()
	one := z.Map[address]
	z.Factory.Rw.RUnlock()
	if one != nil {
		return one
	}

	z.Factory.Rw.Lock()
	one = z.Map[address]
	if one == nil {
		one = &One{Address: address}
		z.One = append(z.One, one)
		z.Map[address] = one
	}
	z.Factory.Rw.Unlock()
	z.GetENS(one)
	fmt.Printf("	%s %s\n", one.Address, one.ENS)
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
		z.One = append(ones, one)
		z.Map[one.Address] = one
		z.Maps.LoopringId[one.LoopringID] = one
		z.Maps.TokenId[one.TokenId] = one
		z.Maps.ABI[one.ABI] = one
	}
	return nil
}
