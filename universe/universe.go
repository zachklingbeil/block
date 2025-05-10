package universe

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

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
	addr, ok := z.Maps.LoopringId[id]
	z.Factory.Rw.RUnlock()

	if !ok {
		peer := &One{LoopringID: id}
		z.GetLoopringAddress(peer)
		return peer
	}

	z.Factory.Rw.RLock()
	one := z.Map[addr]
	z.Factory.Rw.RUnlock()
	return one
}

func (z *Zero) TokenId(id int64) (string, int64) {
	z.Factory.Rw.RLock()
	defer z.Factory.Rw.RUnlock()
	addr := z.Maps.TokenId[id]
	one := z.Map[addr]
	if one == nil {
		return "", 0
	}
	return one.Address, one.Decimals
}

// Format formats a string input as a decimal string based on the given decimals.
func (z *Zero) Format(input string, decimals int64) string {
	value := new(big.Int)
	_, ok := value.SetString(input, 10)
	if !ok {
		return input
	}
	valueStr := value.String()
	dec := int(decimals)
	if len(valueStr) <= dec {
		paddedValue := strings.Repeat("0", dec-len(valueStr)+1) + valueStr
		result := "0." + paddedValue
		return strings.TrimRight(result, "0")
	}

	intPart := valueStr[:len(valueStr)-dec]
	fracPart := valueStr[len(valueStr)-dec:]
	result := intPart + "." + fracPart
	result = strings.TrimRight(result, "0")
	result = strings.TrimSuffix(result, ".")
	return result
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
