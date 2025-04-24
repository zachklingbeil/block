package circuit

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/zachklingbeil/factory"
)

type Circuit struct {
	Factory *factory.Factory
	Map     map[any]any
	String  map[string]any
	Int     map[int]any
}

func NewCircuit(factory *factory.Factory) *Circuit {
	circuit := &Circuit{
		Factory: factory,
		Map:     make(map[any]any),
		String:  make(map[string]any),
		Int:     make(map[int]any),
	}
	return circuit
}

func (c *Circuit) AddString(key string, value any) {
	c.Factory.Mu.Lock()
	defer c.Factory.Mu.Unlock()
	c.String[key] = value
}

func (c *Circuit) AddInt(key int, value any) {
	c.Factory.Mu.Lock()
	defer c.Factory.Mu.Unlock()
	c.Int[key] = value
}

func (c *Circuit) GetString(key string) any {
	c.Factory.Rw.Lock()
	defer c.Factory.Rw.Unlock()
	if value, ok := c.String[key]; ok {
		return value
	}
	return nil
}

func (c *Circuit) GetInt(key int) any {
	c.Factory.Rw.Lock()
	defer c.Factory.Rw.Unlock()
	if value, ok := c.Int[key]; ok {
		return value
	}
	return nil
}

func (c *Circuit) SaveStrings(stringKey string) error {
	c.Factory.Mu.Lock()
	defer c.Factory.Mu.Unlock()

	for key, value := range c.String {
		valueJSON, err := json.Marshal(value)
		if err != nil {
			return err
		}
		err = c.Factory.Redis.HSet(c.Factory.Ctx, stringKey, key, valueJSON).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Circuit) SaveInts(intKey string) error {
	c.Factory.Mu.Lock()
	defer c.Factory.Mu.Unlock()

	for key, value := range c.Int {
		valueJSON, err := json.Marshal(value)
		if err != nil {
			return err
		}
		err = c.Factory.Redis.HSet(c.Factory.Ctx, intKey, fmt.Sprintf("%d", key), valueJSON).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Circuit) Continue(stringKey, intKey string) error {
	c.Factory.Mu.Lock()
	defer c.Factory.Mu.Unlock()
	stringEntries, err := c.Factory.Redis.HGetAll(c.Factory.Ctx, stringKey).Result()
	if err != nil {
		return err
	}
	for key, valueJSON := range stringEntries {
		var value any
		err := json.Unmarshal([]byte(valueJSON), &value)
		if err != nil {
			return err
		}
		c.String[key] = value
	}

	intEntries, err := c.Factory.Redis.HGetAll(c.Factory.Ctx, intKey).Result()
	if err != nil {
		return err
	}
	for key, valueJSON := range intEntries {
		var value any
		err := json.Unmarshal([]byte(valueJSON), &value)
		if err != nil {
			return err
		}
		intKey, err := strconv.Atoi(key)
		if err != nil {
			return err
		}
		c.Int[intKey] = value
	}
	return nil
}
