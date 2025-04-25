package circuit

import (
	"time"
)

func (c *Circuit) Coordinates(input *Raw) ([]any, *Block) {
	for i := range input.Transactions {
		if tx, ok := input.Transactions[i].(map[string]any); ok {
			tx["index"] = i + 1
		}
	}

	transactions := c.Factory.Json.Simplify(input.Transactions, "")
	depth := uint16(len(transactions))

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
		Depth:       depth,
	}

	block := &Block{
		Number: input.Number,
		Zero:   coordinate,
		Ones:   make([]Tx, depth),
	}
	return transactions, block
}
