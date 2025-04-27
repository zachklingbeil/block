package circuit

import (
	"time"

	"github.com/zachklingbeil/block/circuit/value"
	"github.com/zachklingbeil/factory"
)

type Circuit struct {
	Factory *factory.Factory
	Map     map[string]any
	Value   *value.Value
}

func NewCircuit(factory *factory.Factory) *Circuit {
	circuit := &Circuit{
		Factory: factory,
		Map:     make(map[string]any),
		Value:   value.NewValue(factory),
	}
	// circuit.Load()
	// fmt.Printf("%d tokens\n", len(circuit.Tokens))
	return circuit
}

// func (c *Circuit) Continue() error {
// 	c.Factory.Mu.Lock()
// 	defer c.Factory.Mu.Unlock()
// 	source, err := c.Factory.Data.RB.SMembers(c.Factory.Ctx, "value").Result()
// 	if err != nil {
// 		return fmt.Errorf("failed to load values from Redis: %w", err)
// 	}

//		for _, i := range source {
//			var value Value
//			if err := json.Unmarshal([]byte(i), &value); err != nil {
//				return fmt.Errorf("failed to unmarshal value: %w", err)
//			}
//			// c.Map[value.Address] = &value
//			// c.Map[value.ENS] = &value
//			// c.Map[value.LoopringENS] = &value
//			// c.Map[value.LoopringID] = &value
//			// c.Map[value.Symbol] = &value
//			// c.Map[value.Address] = &value
//			// c.Map[value.LoopringID] = &value
//			// c.Map[value.Token] = &value
//			// c.Values = append(c.Values, value)
//		}
//		fmt.Printf("%d\n", len(c.Map))
//		return nil
//	}

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
