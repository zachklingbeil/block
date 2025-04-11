package loopring

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

func (l *Loopring) Listen() {
	for {
		wsApiKey, err := l.fetchWsApiKey()
		if err != nil {
			fmt.Printf("Error fetching WebSocket API key: %v\n", err)
			continue
		}

		wsUrl := "wss://ws.api3.loopring.io/v3/ws?wsApiKey=" + wsApiKey
		conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
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
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("Error reading WebSocket message: %v\n", err)
				break
			}

			if messageType == websocket.TextMessage && string(message) == "ping" {
				if err := conn.WriteMessage(websocket.TextMessage, []byte("pong")); err != nil {
					fmt.Printf("Error sending pong: %v\n", err)
				}
				continue
			}

			var notification struct {
				Data []struct {
					BlockId int64 `json:"blockId"`
				} `json:"data"`
			}
			if err := json.Unmarshal(message, &notification); err != nil {
				fmt.Printf("Error unmarshaling notification: %v\n", err)
				continue
			}

			for _, block := range notification.Data {
				fmt.Printf("%d via notification\n", block.BlockId)
				if err := l.ProcessBlock(block.BlockId); err != nil {
					fmt.Printf("Error processing block %d: %v\n", block.BlockId, err)
				}
			}
		}
	}
}

func (l *Loopring) fetchWsApiKey() (string, error) {
	data, err := l.Factory.Json.In("https://api3.loopring.io/v3/ws/key", "")
	if err != nil {
		return "", fmt.Errorf("failed to fetch WebSocket API key: %v", err)
	}

	var result struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("failed to parse API key response: %v", err)
	}
	return result.Key, nil
}
