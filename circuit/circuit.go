package circuit

import (
	"encoding/json"
	"time"

	"github.com/zachklingbeil/factory"
)

type Circuit struct {
	Factory *factory.Factory
}

type Raw struct {
	Number       int64 `json:"blockId"`
	Timestamp    int64 `json:"createdAt"`
	Size         int64 `json:"blockSize"`
	Transactions []any `json:"transactions"`
}

type Block struct {
	Number int64      `json:"block"`
	Depth  uint16     `json:"depth"`
	Zero   Coordinate `json:"zero"`
	Ones   []Tx       `json:"one"`
}

type Coordinate struct {
	Year        uint8  `json:"year"`
	Month       uint8  `json:"month"`
	Day         uint8  `json:"day"`
	Hour        uint8  `json:"hour"`
	Minute      uint8  `json:"minute"`
	Second      uint8  `json:"second"`
	Millisecond uint16 `json:"millisecond"`
	Index       uint16 `json:"index"`
}

type Tx struct {
	Zero     any             `json:"zero,omitempty"`
	One      any             `json:"one,omitempty"`
	Value    any             `json:"value,omitempty"`
	Token    any             `json:"token,omitempty"`
	Fee      any             `json:"fee,omitempty"`
	FeeToken any             `json:"feeToken,omitempty"`
	Type     string          `json:"type,omitempty"`
	Index    uint16          `json:"index"`
	Raw      json.RawMessage `json:"raw,omitempty"`
}

func NewCircuit(factory *factory.Factory) *Circuit {
	circuit := &Circuit{
		Factory: factory,
	}
	return circuit
}

func (c *Circuit) Coordinates(input *Raw) ([]any, *Block) {
	for i := range input.Transactions {
		if tx, ok := input.Transactions[i].(map[string]any); ok {
			tx["index"] = i + 1
		}
	}
	transactions := c.Factory.Json.Simplify(input.Transactions, "")

	t := time.UnixMilli(input.Timestamp)
	coordinate := Coordinate{
		Year:        uint8(t.Year() - 2015),
		Month:       uint8(t.Month()),
		Day:         uint8(t.Day()),
		Hour:        uint8(t.Hour()),
		Minute:      uint8(t.Minute()),
		Second:      uint8(t.Second()),
		Millisecond: uint16(t.Nanosecond() / 1e6),
		Index:       0,
	}
	depth := uint16(len(transactions))
	block := &Block{
		Number: input.Number,
		Depth:  depth,
		Zero:   coordinate,
		Ones:   make([]Tx, depth),
	}
	return transactions, block
}
