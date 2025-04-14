package process

import (
	"time"

	"github.com/zachklingbeil/factory"
)

type Process struct {
	Factory *factory.Factory
	Blocks  []Block
	Txs     []Txs
}

type Coordinate struct {
	Block       int64 `json:"block"`
	Year        int64 `json:"year"`
	Month       int64 `json:"month"`
	Day         int64 `json:"day"`
	Hour        int64 `json:"hour"`
	Minute      int64 `json:"minute"`
	Second      int64 `json:"second"`
	Millisecond int64 `json:"millisecond"`
	Index       int64 `json:"index"`
}

func (p *Process) Coordinates(block int64, timestamp int64, index int64) Coordinate {
	t := time.UnixMilli(timestamp)
	coordinates := Coordinate{
		Block:       block,
		Year:        int64(t.Year() - 2015),
		Month:       int64(t.Month()),
		Day:         int64(t.Day()),
		Hour:        int64(t.Hour()),
		Minute:      int64(t.Minute()),
		Second:      int64(t.Second()),
		Millisecond: int64(t.Nanosecond() / 1e6),
		Index:       index,
	}
	return coordinates
}
