package universe

import (
	"encoding/json"
	"time"
)

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
	transactions := z.Factory.Json.Simplify(input.Transactions, "")
	depth := uint16(len(transactions))

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
	return transactions, coordinate
}
