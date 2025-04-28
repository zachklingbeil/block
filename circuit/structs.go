package circuit

import (
	"encoding/json"

	"github.com/zachklingbeil/block/value"
)

// Loopring
type Raw struct {
	Number       int64 `json:"blockId"`
	Timestamp    int64 `json:"createdAt"`
	Size         int64 `json:"blockSize"`
	Transactions []any `json:"transactions"`
}

type Block struct {
	Number int64      `json:"block"`
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
	Depth       uint16 `json:"depth,omitempty"`
}

type Tx struct {
	Zero     any             `json:"zero,omitempty"`
	One      any             `json:"one,omitempty"`
	Value    any             `json:"value,omitempty"`
	Token    any             `json:"token,omitempty"`
	Fee      any             `json:"fee,omitempty"`
	FeeToken any             `json:"feeToken,omitempty"`
	Type     any             `json:"type,omitempty"`
	Index    any             `json:"index"`
	Raw      json.RawMessage `json:"raw,omitempty"`
}

type One struct {
	Peer       *value.Peer  `json:"peer,omitempty"`
	Token      *value.Token `json:"token,omitempty"`
	Block      *Block       `json:"block,omitempty"`
	Coordinate *Coordinate  `json:"coordinate,omitempty"`
	Tx         *Tx          `json:"tx,omitempty"`
}
