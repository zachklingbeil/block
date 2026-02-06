package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	client, err := rpc.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Failed to connect to geth: %v", err)
	}
	defer client.Close()

	ethClient := ethclient.NewClient(client)
	defer ethClient.Close()

	var raw json.RawMessage
	err = client.Call(&raw, "eth_getBlockByNumber", "latest", true)
	if err != nil {
		log.Fatalf("Failed to fetch latest block: %v", err)
	}

	var block map[string]interface{}
	if err := json.Unmarshal(raw, &block); err != nil {
		log.Fatalf("Failed to unmarshal block: %v", err)
	}

	output, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal block: %v", err)
	}

	filename := "block.json"
	if err := os.WriteFile(filename, output, 0644); err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	fmt.Printf("Latest block written to %s\n", filename)
}
