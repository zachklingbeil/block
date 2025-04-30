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

// ProcessBlocks orchestrates the loading, processing, and saving of blocks.
func (v *Value) ProcessBlocks() error {
	if err := v.LoadBlocks(); err != nil {
		return fmt.Errorf("failed to load blocks: %w", err)
	}

	if err := v.SaveBlocks(); err != nil {
		return fmt.Errorf("failed to save blocks: %w", err)
	}
	return nil
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
	zAddArgs := make([]redis.Z, 0, len(v.Blocks))

	for _, block := range v.Blocks {
		v.Factory.Json.Simplify([]any{block}, "")

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

// ProcessTxs iterates through all transactions in the blocks and applies a processing function to each Tx.
func (v *Value) ProcessTxs(processFunc func(*Tx) error) error {
	for i := range v.Blocks {
		block := &v.Blocks[i]
		for j := range block.Ones {
			tx := &block.Ones[j]
			if err := processFunc(tx); err != nil {
				return fmt.Errorf("failed to process transaction in block %d: %w", block.Number, err)
			}
		}
	}
	return nil
}

func (v *Value) HandleNewPeers() error {
	uniqueIDs := make(map[string]struct{})
	uniqueStrings := make(map[string]struct{})

	for _, block := range v.Blocks {
		for _, tx := range block.Ones {
			if zeroStr, ok := tx.Zero.(string); ok {
				if isHexadecimal(zeroStr) {
					uniqueStrings[zeroStr] = struct{}{}
				} else {
					uniqueIDs[zeroStr] = struct{}{}
				}
			}
			if oneStr, ok := tx.One.(string); ok {
				if isHexadecimal(oneStr) {
					uniqueStrings[oneStr] = struct{}{}
				} else {
					uniqueIDs[oneStr] = struct{}{}
				}
			}
		}
	}

	ids := make([]any, 0, len(uniqueIDs))
	for id := range uniqueIDs {
		ids = append(ids, id)
	}

	strings := make([]any, 0, len(uniqueStrings))
	for str := range uniqueStrings {
		strings = append(strings, str)
	}

	if _, err := v.Factory.Data.RB.SAdd(v.Factory.Ctx, "newPeersID", ids...).Result(); err != nil {
		return fmt.Errorf("failed to store IDs in Redis: %w", err)
	}

	if _, err := v.Factory.Data.RB.SAdd(v.Factory.Ctx, "newPeersString", strings...).Result(); err != nil {
		return fmt.Errorf("failed to store hexadecimal strings in Redis: %w", err)
	}
	return nil
}

// Helper function to check if a string is hexadecimal
func isHexadecimal(s string) bool {
	return len(s) > 2 && s[:2] == "0x"
}
