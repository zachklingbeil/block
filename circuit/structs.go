package circuit

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
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
	Type     string          `json:"type,omitempty"`
	Index    uint16          `json:"index"`
	Raw      json.RawMessage `json:"raw,omitempty"`
}

type Token struct {
	Token      string `json:"token"`
	TokenId    int64  `json:"tokenId"`
	LoopringID int64  `json:"loopringId,omitempty"`
	Decimals   int64  `json:"decimals"`
	Address    string `json:"address"`
}

type Peer struct {
	ENS         string `json:"ens"`
	LoopringENS string `json:"loopringEns"`
	LoopringID  string `json:"loopringId"`
	Address     string `json:"address"`
}

// UpdateMessage represents an update for a specific slice.
type UpdateMessage[T any] struct {
	Key   string // Redis key for the slice
	Slice []T    // Updated slice
}

// FetchAndInitialize fetches and initializes all data needed at runtime and sends it to the updates channel.
func FetchAndInitialize[T any](ctx context.Context, rb *redis.Client, key string, updates chan<- UpdateMessage[T]) []T {
	var items []T
	source, err := rb.SMembers(ctx, key).Result()
	if err != nil {
		log.Fatalf("Failed to fetch items from Redis set '%s': %v", key, err)
	}

	for _, s := range source {
		var item T
		if err := json.Unmarshal([]byte(s), &item); err != nil {
			log.Printf("Skipping invalid item: %v (data: %s)", err, s)
			continue
		}
		items = append(items, item)
	}

	// Send the initialized slice to the updates channel.
	updates <- UpdateMessage[T]{Key: key, Slice: items}
	return items
}

// StoreToRedis stores a slice of items into a Redis set.
func StoreToRedis[T any](ctx context.Context, rb *redis.Client, key string, items []T) error {
	pipe := rb.Pipeline()

	for _, item := range items {
		data, err := json.Marshal(item)
		if err != nil {
			log.Printf("Failed to serialize item: %v (error: %v)", item, err)
			continue
		}
		pipe.SAdd(ctx, key, data)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("Failed to store items in Redis set '%s': %v", key, err)
		return err
	}
	return nil
}

// ListenAndSyncToRedisMulti listens for updates to multiple slices and updates Redis accordingly.
func ListenAndSyncToRedisMulti[T any](ctx context.Context, rb *redis.Client, updates <-chan UpdateMessage[T]) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping multi-slice sync to Redis")
			return
		case msg := <-updates:
			err := StoreToRedis(ctx, rb, msg.Key, msg.Slice)
			if err != nil {
				log.Printf("Failed to sync updated slice to Redis for key '%s': %v", msg.Key, err)
			} else {
				log.Printf("Successfully synced updated slice to Redis for key '%s'", msg.Key)
			}
		}
	}
}
