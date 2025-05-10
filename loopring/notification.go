package loopring

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

func (l *Loopring) CurrentBlock() int64 {
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
				fmt.Printf("block %d\n", block.Number)
				l.Factory.State.Count("loop.block", block.Number)
				l.BlockByBlock(block.Number)
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
