package loop

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

func (l *Loopring) Loop() {
	l.Types = &Type{
		Deposits:       make([]DW, 0, 300),
		Withdrawals:    make([]DW, 0, 300),
		Swaps:          make([]Swap, 0, 300),
		Transfers:      make([]Transfer, 0, 300),
		Mints:          make([]Mint, 0, 300),
		NftData:        make([]NftData, 0, 300),
		AmmUpdates:     make([]AmmUpdate, 0, 300),
		AccountUpdates: make([]AccountUpdate, 0, 300),
		TBD:            make([]any, 0, 10),
	}
	l.Transactions = make([]Transaction, 0, 1000)
}

// Simplified GetCurrentBlockNumber
func (l *Loopring) currentBlock() int64 {
	data, err := l.Factory.Json.In("https://api3.loopring.io/api/v3/block/getBlock", "")
	if err != nil {
		fmt.Printf("Failed to fetch block data: %v\n", err)
		return 0
	}
	var block struct {
		Number int64 `json:"blockId"`
	}
	err = json.Unmarshal(data, &block)
	if err != nil {
		fmt.Printf("Failed to parse block data: %v\n", err)
		return 0
	}
	return block.Number
}

func (l *Loopring) Listen() {
	for {
		key := l.fetchWsApiKey()
		if key == "" {
			continue
		}
		url := "wss://ws.api3.loopring.io/v3/ws?wsApiKey=" + key
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			fmt.Printf("Error connecting to WebSocket: %v\n", err)
			continue
		}
		defer conn.Close()

		sub := map[string]any{
			"op":     "sub",
			"topics": []map[string]string{{"topic": "blockgen"}},
		}
		if err := conn.WriteJSON(sub); err != nil {
			fmt.Printf("Error subscribing to topic: %v\n", err)
			continue
		}

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("Error reading WebSocket message: %v\n", err)
				break
			}
			if msgType == websocket.TextMessage && string(msg) == "ping" {
				conn.WriteMessage(websocket.TextMessage, []byte("pong"))
				continue
			}

			var resp struct {
				Data []struct {
					Number int64 `json:"blockId"`
				} `json:"data"`
			}
			if err := json.Unmarshal(msg, &resp); err != nil {
				continue
			}
			for _, block := range resp.Data {
				fmt.Printf("%d\n", block.Number)
				if err := l.fetchBlock(block.Number); err != nil {
					fmt.Printf("Error processing block %d: %v\n", block.Number, err)
				}
			}
		}
	}
}

func (l *Loopring) fetchWsApiKey() string {
	data, err := l.Factory.Json.In("https://api3.loopring.io/v3/ws/key", "")
	if err != nil {
		fmt.Printf("Error fetching WebSocket API key: %v\n", err)
		return ""
	}
	var result struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		fmt.Printf("Error parsing API key response: %v\n", err)
		return ""
	}
	return result.Key
}
