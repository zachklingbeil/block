package circuit

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/zachklingbeil/factory"
)

type Circuit struct {
	Factory *factory.Factory
	Map     map[any]any
	String  map[string]any
	Int     map[int64]any
	Redis   *redis.Client
}

func NewCircuit(factory *factory.Factory) *Circuit {
	redis, err := factory.Db.ConnectRedis(0, factory.Ctx)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to Redis: %v", err))
	}

	circuit := &Circuit{
		Factory: factory,
		Map:     make(map[any]any),
		String:  make(map[string]any),
		Int:     make(map[int64]any),
		Redis:   redis,
	}
	return circuit
}

func (c *Circuit) AddString(key string, value any) {
	c.Factory.Mu.Lock()
	defer c.Factory.Mu.Unlock()
	c.String[key] = value

	v, err := json.Marshal(value)
	if err != nil {
		fmt.Printf("Failed to marshal value for key %s: %v\n", key, err)
		return
	}
	if err := c.Redis.HSet(c.Factory.Ctx, "strings", key, v).Err(); err != nil {
		fmt.Printf("Failed to save key %s to Redis: %v\n", key, err)
	}
}

func (c *Circuit) AddInt(key int64, value any) {
	c.Factory.Mu.Lock()
	defer c.Factory.Mu.Unlock()
	c.Int[key] = value

	v, err := json.Marshal(value)
	if err != nil {
		fmt.Printf("Failed to marshal value for key %d: %v\n", key, err)
		return
	}
	if err := c.Redis.HSet(c.Factory.Ctx, "ints", key, v).Err(); err != nil {
		fmt.Printf("Failed to save key %d to Redis: %v\n", key, err)
	}
}

func (c *Circuit) Continue() error {
	c.Factory.Mu.Lock()
	defer c.Factory.Mu.Unlock()

	strings, err := c.Redis.HGetAll(c.Factory.Ctx, "strings").Result()
	if err != nil {
		return fmt.Errorf("failed to load string map from Redis: %w", err)
	}
	for key, v := range strings {
		var value any
		if err := json.Unmarshal([]byte(v), &value); err != nil {
			return fmt.Errorf("failed to unmarshal value for key %s: %w", key, err)
		}
		c.String[key] = value
	}

	ints, err := c.Redis.HGetAll(c.Factory.Ctx, "ints").Result()
	if err != nil {
		return fmt.Errorf("failed to load int map from Redis: %w", err)
	}
	for key, v := range ints {
		var value any
		if err := json.Unmarshal([]byte(v), &value); err != nil {
			return fmt.Errorf("failed to unmarshal value for key %s: %w", key, err)
		}
		intKey, err := strconv.Atoi(key)
		if err != nil {
			return fmt.Errorf("failed to convert key %s to int: %w", key, err)
		}
		c.Int[int64(intKey)] = value
	}
	return nil
}
