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

		err = conn.WriteJSON(map[string]any{
			"op":     "sub",
			"topics": []map[string]string{{"topic": "blockgen"}},
		})
		if err != nil {
			fmt.Printf("Error subscribing to topic: %v\n", err)
			conn.Close()
			continue
		}

		var response map[string]any
		err = conn.ReadJSON(&response)
		if err != nil {
			fmt.Printf("Error reading subscription response: %v\n", err)
			conn.Close()
			continue
		}

		if result, ok := response["result"].(map[string]any); !ok || result["status"] != "OK" {
			fmt.Printf("Subscription failed: %v\n", response)
			conn.Close()
			continue
		}

		fmt.Println("Subscribed to blockgen topic successfully")

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("Error reading WebSocket message: %v\n", err)
				break
			}

			if messageType == websocket.TextMessage && string(message) == "ping" {
				err = conn.WriteMessage(websocket.TextMessage, []byte("pong"))
				if err != nil {
					fmt.Printf("Error sending pong: %v\n", err)
				}
				continue
			}

			var notification map[string]any
			err = json.Unmarshal(message, &notification)
			if err != nil {
				fmt.Printf("Error unmarshaling notification: %v\n", err)
				continue
			}

			if topic, ok := notification["topic"].(string); ok && topic == "blockgen" {
				fmt.Println("Received blockgen notification, calling FetchBlocks")
				l.FetchBlocks()
			} else {
				fmt.Printf("Unexpected notification: %v\n", notification)
			}
		}
		conn.Close()
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
