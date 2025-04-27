package value

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zachklingbeil/factory"
)

type State struct {
	Map     map[string]any
	Factory *factory.Factory
}

func NewState(factory *factory.Factory) *State {
	state := &State{
		Map:     make(map[string]any),
		Factory: factory,
	}

	state.Add("t0", time.Now().Format("08:04:05.0000000"))
	return state
}

func (s *State) Add(key string, value any) error {
	s.Factory.Mu.Lock()
	defer s.Factory.Mu.Unlock()

	// Update the map with the new key-value pair
	s.Map[key] = value

	// Serialize the state map to JSON
	state, err := json.Marshal(s.Map)
	if err != nil {
		return fmt.Errorf("failed to marshal state map: %w", err)
	}

	// Use the current timestamp as the score for ZAdd
	score := float64(time.Now().UnixNano()) / 1e9 // Convert nanoseconds to seconds
	member := redis.Z{
		Score:  score,
		Member: state,
	}

	// Add the serialized state to the Redis sorted set
	if err := s.Factory.Data.RB.ZAdd(s.Factory.Ctx, "state", member).Err(); err != nil {
		return fmt.Errorf("failed to add state to Redis: %w", err)
	}

	return nil
}

func (s *State) Get() {
	s.Factory.Rw.RLock()
	defer s.Factory.Rw.RUnlock()
	s.Factory.Json.Print(s.Map)
}
