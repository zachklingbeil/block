package value

import (
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

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
	Zero     any `json:"zero,omitempty"`
	One      any `json:"one,omitempty"`
	Value    any `json:"value,omitempty"`
	Token    any `json:"token,omitempty"`
	Fee      any `json:"fee,omitempty"`
	FeeToken any `json:"feeToken,omitempty"`
	Type     any `json:"type,omitempty"`
	Index    any `json:"index"`
	// Raw      json.RawMessage `json:"raw,omitempty"`
}

func (v *Value) LoadBlocks() error {
	source, err := v.Factory.Data.RB.ZRange(v.Factory.Ctx, "blocks", 0, -1).Result()
	if err != nil {
		return fmt.Errorf("failed to retrieve blocks from Redis: %w", err)
	}
	blocks := make([]Block, 0, len(source))
	for _, blockJSON := range source {
		var block Block
		if err := json.Unmarshal([]byte(blockJSON), &block); err != nil {
			continue
		}
		blocks = append(blocks, block)
	}
	v.Blocks = blocks
	return nil
}

func (v *Value) SaveBlocks() error {
	if v.Blocks == nil {
		return fmt.Errorf("no blocks to save")
	}

	zAddArgs := make([]redis.Z, 0, len(v.Blocks))
	for _, block := range v.Blocks {
		blockJSON, err := json.Marshal(block)
		if err != nil {
			return fmt.Errorf("failed to serialize block: %w", err)
		}
		zAddArgs = append(zAddArgs, redis.Z{
			Score:  float64(block.Number),
			Member: blockJSON,
		})
	}
	_, err := v.Factory.Data.RB.ZAdd(v.Factory.Ctx, "blocks2", zAddArgs...).Result()
	if err != nil {
		return fmt.Errorf("failed to save blocks to Redis: %w", err)
	}
	return nil
}

func (v *Value) UpdateBlockOnes(updateFunc func(*Tx)) error {
	for i := range v.Blocks {
		block := &v.Blocks[i]
		for j := range block.Ones {
			updateFunc(&block.Ones[j])
		}
	}
	return nil
}

func (v *Value) HandleNewPeers() error {
	uniqueValues := make(map[string]struct{})
	for _, block := range v.Blocks {
		for _, tx := range block.Ones {
			if zeroStr, ok := tx.Zero.(string); ok {
				uniqueValues[zeroStr] = struct{}{}
			}

			if oneStr, ok := tx.One.(string); ok {
				uniqueValues[oneStr] = struct{}{}
			}
		}
	}
	values := make([]any, 0, len(uniqueValues))
	for value := range uniqueValues {
		values = append(values, value)
	}
	_, err := v.Factory.Data.RB.SAdd(v.Factory.Ctx, "newPeers", values...).Result()
	if err != nil {
		return fmt.Errorf("failed to store unique values in Redis: %w", err)
	}
	return nil
}
