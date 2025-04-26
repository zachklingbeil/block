package circuit

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Block struct {
	Number int64      `json:"block"`
	Zero   Coordinate `json:"zero"`
	Ones   []Tx       `json:"one"`
}

// Loopring
type Raw struct {
	Number       int64 `json:"blockId"`
	Timestamp    int64 `json:"createdAt"`
	Size         int64 `json:"blockSize"`
	Transactions []any `json:"transactions"`
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
	Symbol     string `json:"symbol"`
	Address    string `json:"address"`
	LoopringID string `json:"accountId,omitempty"`
	TokenId    int64  `json:"tokenId"`
	Decimals   int    `json:"decimals"`
}

type Peer struct {
	ENS         string `json:"ens"`
	LoopringENS string `json:"loopringEns"`
	LoopringID  string `json:"loopringId"`
	Address     string `json:"address"`
}

type Value struct {
	Symbol      string `json:"symbol,omitempty"`
	ENS         string `json:"ens,omitempty"`
	LoopringENS string `json:"loopringEns,omitempty"`
	Address     string `json:"address,omitempty"`
	LoopringID  string `json:"loopringId,omitempty"`
	Token       string `json:"token,omitempty"`
	Decimals    string `json:"decimals,omitempty"`
}

func (c *Circuit) ConvertToValues() {
	values := make([]Value, 0)
	for _, token := range c.Tokens {
		values = append(values, Value{
			Symbol:     token.Symbol,
			Address:    token.Address,
			LoopringID: token.LoopringID,
			Token:      strconv.FormatInt(token.TokenId, 10),
			Decimals:   strconv.Itoa(token.Decimals),
		})
	}

	// Convert Peers to Values and append to the temporary slice
	for _, peer := range c.Peers {
		values = append(values, Value{
			ENS:         peer.ENS,
			LoopringENS: peer.LoopringENS,
			LoopringID:  peer.LoopringID,
			Address:     peer.Address,
		})
	}

	// Marshal the []Value slice into JSON
	jsonData, err := json.Marshal(values)
	if err != nil {
		panic(fmt.Errorf("failed to marshal values to JSON: %w", err))
	}

	// Unmarshal the JSON back into a slice of maps for Simplify
	var rawValues []any
	if err := json.Unmarshal(jsonData, &rawValues); err != nil {
		panic(fmt.Errorf("failed to unmarshal JSON to []any: %w", err))
	}

	simplified := c.Factory.Json.Simplify(rawValues, "")
	c.Values = make([]Value, len(simplified))
	for i, item := range simplified {
		itemJSON, _ := json.Marshal(item)
		_ = json.Unmarshal(itemJSON, &c.Values[i])
	}

	for _, value := range c.Values {
		valueJSON, err := json.Marshal(value)
		if err != nil {
			panic(fmt.Errorf("failed to marshal simplified value to JSON: %w", err))
		}
		if err := c.Factory.Redis.SAdd(c.Factory.Ctx, "values", valueJSON).Err(); err != nil {
			panic(fmt.Errorf("failed to store value in Redis: %w", err))
		}
	}
}
