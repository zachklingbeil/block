package circuit

import (
	"encoding/json"
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Circuit struct {
	Factory        *factory.Factory
	Map            map[string]any
	Values         []Value
	LoopringApiKey string
}

func NewCircuit(factory *factory.Factory) *Circuit {
	circuit := &Circuit{
		Factory: factory,
		Map:     make(map[string]any),
	}

	return circuit
}

func (c *Circuit) Continue() error {
	c.Factory.Mu.Lock()
	defer c.Factory.Mu.Unlock()
	source, err := c.Factory.Redis.SMembers(c.Factory.Ctx, "value").Result()
	if err != nil {
		return fmt.Errorf("failed to load values from Redis: %w", err)
	}

	c.Values = make([]Value, 0, len(source))
	for _, i := range source {
		var value Value
		if err := json.Unmarshal([]byte(i), &value); err != nil {
			return fmt.Errorf("failed to unmarshal value: %w", err)
		}
		c.Map[value.Address] = &value
		c.Map[value.ENS] = &value
		c.Map[value.LoopringENS] = &value
		c.Map[value.LoopringID] = &value
		c.Map[value.Symbol] = &value
		c.Map[value.Address] = &value
		c.Map[value.LoopringID] = &value
		c.Map[value.Token] = &value
		c.Values = append(c.Values, value)
	}
	fmt.Printf("%d\n", len(c.Map))
	return nil
}
