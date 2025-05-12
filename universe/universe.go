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
	}
	z.LoadOnes()
	return z
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

func (z *Zero) Who(id int64) string {
	peer := z.LoopringId(id)
	if peer == nil {
		return ""
	}
	if peer.ENS != "" && peer.ENS != "." {
		return peer.ENS
	}
	if peer.LoopringENS != "" && peer.LoopringENS != "." && peer.LoopringENS != "!" {
		return peer.LoopringENS
	}
	return peer.Address
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

	left := valueStr[:len(valueStr)-dec]
	right := valueStr[len(valueStr)-dec:]
	result := left + "." + right
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
		z.One = append(ones, one)
		z.Map[one.Address] = one
		z.Maps.LoopringId[one.LoopringID] = one
		z.Maps.TokenId[one.TokenId] = one
		z.Maps.ABI[one.ABI] = one
	}
	return nil
}
