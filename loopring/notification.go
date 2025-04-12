package loopring

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

func (l *Loopring) Listen() {
	for {
		key := l.fetchWsApiKey()
		url := "wss://ws.api3.loopring.io/v3/ws?wsApiKey=" + key
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			fmt.Printf("Error connecting to WebSocket: %v\n", err)
			continue
		}
		defer conn.Close()

		if err := conn.WriteJSON(map[string]any{
			"op":     "sub",
			"topics": []map[string]string{{"topic": "blockgen"}},
		}); err != nil {
			fmt.Printf("Error subscribing to topic: %v\n", err)
			continue
		}

		for {
			m, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("Error reading WebSocket message: %v\n", err)
				break
			}

			if m == websocket.TextMessage && string(message) == "ping" {
				conn.WriteMessage(websocket.TextMessage, []byte("pong"))
				continue
			}

			var newBlock struct {
				Data []struct {
					Number int64 `json:"blockId"`
				} `json:"data"`
			}
			if err := json.Unmarshal(message, &newBlock); err != nil {
				continue
			}

			for _, block := range newBlock.Data {
				fmt.Printf("%d\n", block.Number)
				if err := l.ProcessBlock(block.Number); err != nil {
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
