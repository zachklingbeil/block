package universe

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Zero struct {
	// Factory *factory.Factory
	One    []*One
	Map    map[string]*One
	Maps   *Maps
	Format *Format
}

type Format struct{}

type One struct {
	ENS      string `json:"ens,omitempty"`
	Address  string `json:"address"`
	Token    string `json:"token,omitempty"`
	Decimals int64  `json:"decimals,omitempty"`
	TokenId  int64  `json:"tokenId,omitempty"`
	ABI      string `json:"abi,omitempty"`
}

type Maps struct {
	TokenId map[int64]*One
	ENS     map[string]*One
	Token   map[string]*One
	ABI     map[string]*One
}

func NewZero() *Zero {
	z := &Zero{
		// Factory: factory,
		Map: make(map[string]*One),
		Maps: &Maps{
			TokenId: make(map[int64]*One),
			ENS:     make(map[string]*One),
			Token:   make(map[string]*One),
			ABI:     make(map[string]*One),
		},
		Format: &Format{},
	}
	return z
}

type Raw struct {
	Number       int64 `json:"blockId"`
	Timestamp    int64 `json:"createdAt"`
	Size         int64 `json:"blockSize"`
	Transactions []any `json:"transactions"`
}

type Block struct {
	Zero Coordinate `json:"zero"`
	Ones []Tx       `json:"one"`
}

type Coordinate struct {
	Number      int64  `json:"block"`
	Year        uint8  `json:"year"`
	Month       uint8  `json:"month"`
	Day         uint8  `json:"day"`
	Hour        uint8  `json:"hour"`
	Minute      uint8  `json:"minute"`
	Second      uint8  `json:"second"`
	Millisecond uint16 `json:"millisecond"`
	Index       uint16 `json:"index"`
	Depth       uint16 `json:"depth,omitempty"`
}

type Tx struct {
	Zero     string          `json:"zero,omitempty"`
	One      string          `json:"one,omitempty"`
	Value    string          `json:"value,omitempty"`
	Token    string          `json:"token,omitempty"`
	Fee      string          `json:"fee,omitempty"`
	For      string          `json:"for,omitempty"`
	ForToken string          `json:"forToken,omitempty"`
	FeeToken string          `json:"feeToken,omitempty"`
	Type     string          `json:"type,omitempty"`
	Index    uint16          `json:"index"`
	Nonce    int64           `json:"nonce,omitempty"`
	Raw      json.RawMessage `json:"raw,omitempty"`
}

func (z *Zero) Coordinates(input *Raw) ([]any, *Coordinate) {
	for i := range input.Transactions {
		if tx, ok := input.Transactions[i].(map[string]any); ok {
			tx["index"] = i + 1
		}
	}
	depth := uint16(len(input.Transactions))

	t := time.UnixMilli(input.Timestamp)
	coordinate := &Coordinate{
		Number:      input.Number,
		Year:        uint8(t.Year() - 2015),
		Month:       uint8(t.Month()),
		Day:         uint8(t.Day()),
		Hour:        uint8(t.Hour()),
		Minute:      uint8(t.Minute()),
		Second:      uint8(t.Second()),
		Millisecond: uint16(t.Nanosecond() / 1e6),
		Index:       0,
		Depth:       depth,
	}
	return input.Transactions, coordinate
}

// Format formats a string input as a decimal string based on the given decimals.
func (f *Format) Value(input string, decimals int64) string {
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

func (f *Format) Peer(address string) string {
	address = strings.ToLower(address)
	if strings.HasPrefix(address, "0x") || strings.HasSuffix(address, ".eth") {
		return address
	}
	return address
}

// // hex -> .eth
// func (z *Zero) GetENS(peer *One) {
// 	if peer.ENS != "" && peer.ENS != "." {
// 		return
// 	}
// 	ensName, err := ens.ReverseResolve(<ethclient>, common.HexToAddress(peer.Address))
// 	if err != nil || ensName == "" {
// 		peer.ENS = "."
// 	} else {
// 		peer.ENS = z.Format.Peer(ensName)
// 	}
// }

func (z *Zero) HelloUniverse(value string) *One {

	peer, exists := z.Map[value]
	if !exists {
		peer = &One{}
		if common.IsHexAddress(value) {
			peer.Address = strings.ToLower(value)
		}
		z.One = append(z.One, peer)
		z.Map[peer.Address] = peer
	} else if peer.Address == "" && common.IsHexAddress(value) {
		peer.Address = strings.ToLower(value)
		z.Map[peer.Address] = peer
	}

	// z.GetENS(peer)
	fmt.Printf("	%s %s\n", peer.Address, peer.ENS)
	return peer
}
