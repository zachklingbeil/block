package circuit

import (
	"encoding/json"
	"maps"
	"time"

	"github.com/zachklingbeil/factory"
)

type Circuit struct {
	One     map[any]any
	Factory *factory.Factory
}

type Block struct {
	Number      int64  `json:"block"`
	Year        uint8  `json:"year"`
	Month       uint8  `json:"month"`
	Day         uint8  `json:"day"`
	Hour        uint8  `json:"hour"`
	Minute      uint8  `json:"minute"`
	Second      uint8  `json:"second"`
	Millisecond uint16 `json:"millisecond"`
	Index       uint16 `json:"index"`
	Count       uint16 `json:"count"`
	Txs         []Tx   `json:"transactions"`
}

type Tx struct {
	Zero        any             `json:"zero,omitempty"`
	One         any             `json:"one,omitempty"`
	Value       string          `json:"value,omitempty"`
	Token       any             `json:"token,omitempty"`
	Fee         any             `json:"fee,omitempty"`
	FeeToken    int64           `json:"feeToken,omitempty"`
	OneValue    string          `json:"oneValue,omitempty"`
	OneToken    int64           `json:"oneToken,omitempty"`
	OneFee      any             `json:"oneFee,omitempty"`
	OneFeeToken int64           `json:"oneFeeToken,omitempty"`
	Type        string          `json:"type,omitempty"`
	Index       uint16          `json:"index"`
	Raw         json.RawMessage `json:"raw,omitempty"`
}

func NewCircuit(factory *factory.Factory) *Circuit {
	circuit := &Circuit{
		One:     make(map[any]any),
		Factory: factory,
	}
	return circuit
}

// Add safely adds one to Circuit.One
func (c *Circuit) Add(one map[any]any) {
	c.Factory.Mu.Lock()
	defer c.Factory.Mu.Unlock()
	maps.Copy(c.One, one)
}

func (c *Circuit) Read(zero Block) any {
	c.Factory.Mu.Lock()
	defer c.Factory.Mu.Unlock()
	value := c.One[zero]
	return value
}

func (c *Circuit) Coordinates(blockNumber, timestamp int64, ones []any) (*Block, error) {
	for i := range ones {
		if tx, ok := ones[i].(map[string]any); ok {
			tx["index"] = i + 1
		}
	}

	count := len(ones)
	c.Factory.Json.Simplify(ones, "")

	t := time.UnixMilli(timestamp)
	block := &Block{
		Year:        uint8(t.Year() - 2015),
		Month:       uint8(t.Month()),
		Day:         uint8(t.Day()),
		Hour:        uint8(t.Hour()),
		Minute:      uint8(t.Minute()),
		Second:      uint8(t.Second()),
		Millisecond: uint16(t.Nanosecond() / 1e6),
		Number:      blockNumber,
		Index:       0,
		Count:       uint16(count),
		Txs:         make([]Tx, count),
	}
	return block, nil
}
