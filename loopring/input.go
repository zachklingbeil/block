package loopring

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Txs     []any
}

func NewLoopring(factory *factory.Factory) *Loopring {
	return &Loopring{Factory: factory}
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

func (l *Loopring) Listen() {
	for {
		wsApiKey, err := l.fetchWsApiKey()
		if err != nil {
			continue
		}

		wsUrl := "wss://ws.api3.loopring.io/v3/ws?wsApiKey=" + wsApiKey
		conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
		if err != nil {
			continue
		}
		defer conn.Close()

		conn.WriteJSON(map[string]any{
			"op":             "sub",
			"sequence":       30006,
			"unsubscribeAll": false,
			"topics":         []map[string]string{{"topic": "blockgen"}},
		})

		var response map[string]any
		conn.ReadJSON(&response)
		if result, ok := response["result"].(map[string]any); !ok || result["status"] != "OK" {
			conn.Close()
			continue
		}

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			if messageType == websocket.TextMessage && string(message) == "ping" {
				conn.WriteMessage(websocket.TextMessage, []byte("pong"))
				continue
			}

			var notification map[string]any
			if json.Unmarshal(message, &notification) == nil {
				if topic, ok := notification["topic"].(string); ok && topic == "blockgen" {
					l.FetchBlocks()
				}
			}
		}
		conn.Close()
	}
}

func (l *Loopring) FetchBlocks() {
	current := l.currentBlock()
	blockHeight := l.blockHeight()

	if blockHeight >= current {
		return
	}

	for i := blockHeight + 1; i <= current; i++ {
		fmt.Printf("%d\n", i)
		if err := l.ProcessBlock(i); err != nil {
			fmt.Printf("Failed to process block %d: %v\n", i, err)
			continue
		}
	}
}

// Helper function to fetch the highest block ID
func (l *Loopring) blockHeight() int64 {
	var blockHeight int64
	err := l.Factory.Db.QueryRow(`SELECT COALESCE(MAX(block), 0) FROM loopring`).Scan(&blockHeight)
	if err != nil {
		fmt.Printf("Failed to fetch the highest block ID: %v\n", err)
		return 0
	}
	return blockHeight
}

// Simplified GetCurrentBlockNumber
func (l *Loopring) currentBlock() int64 {
	data, err := l.Factory.Json.In("https://api3.loopring.io/api/v3/block/getBlock", "")
	if err != nil {
		fmt.Printf("Failed to fetch block data: %v\n", err)
		return 0
	}

	var blockData struct {
		BlockId int64 `json:"blockId"`
	}

	err = json.Unmarshal(data, &blockData)
	if err != nil {
		fmt.Printf("Failed to parse block data: %v\n", err)
		return 0
	}

	return blockData.BlockId
}
