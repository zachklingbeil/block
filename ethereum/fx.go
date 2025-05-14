package ethereum

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
)

// ListenForNewBlocksIPC subscribes to new block headers using IPC and updates e.Header.
func (e *Ethereum) Listen() error {
	headers := make(chan *types.Header)
	sub, err := e.Factory.Eth.SubscribeNewHead(e.Factory.Ctx, headers)
	if err != nil {
		return fmt.Errorf("failed to subscribe to new heads: %w", err)
	}

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case <-e.Factory.Ctx.Done():
				return
			case err := <-sub.Err():
				log.Printf("Subscription error: %v", err)
				return
			case header := <-headers:
				e.Header = header.Number
				e.BlockByBlock()
			}
		}
	}()
	return nil
}

// // BlockByBlock processes and stores the latest block, one at a time.
// func (e *Ethereum) BlockByBlock() error {
// 	block, err := e.Factory.Eth.BlockByNumber(e.Factory.Ctx, e.Header)
// 	if err != nil {
// 		log.Printf("Error fetching block %d: %v", e.Header.Uint64(), err)
// 		return err
// 	}

// 	// Process all transactions in the block
// 	if err := e.TransactionByTransaction(block.Transactions()); err != nil {
// 		log.Printf("Error processing transactions in block %d: %v", block.NumberU64(), err)
// 	}

// 	return nil
// }

// // TransactionByTransaction processes a slice of transactions and their logs.
// func (e *Ethereum) TransactionByTransaction(txs types.Transactions) error {
// 	for _, tx := range txs {
// 		receipt, err := e.Factory.Eth.TransactionReceipt(e.Factory.Ctx, tx.Hash())
// 		if err != nil {
// 			log.Printf("Error fetching receipt for tx %s: %v", tx.Hash().Hex(), err)
// 			continue
// 		}

// 		decodedLogs := make([]map[string]interface{}, 0, len(receipt.Logs))
// 		for _, logEntry := range receipt.Logs {
// 			decoded, _ := e.LogByLog(logEntry, tx)
// 			if decoded != nil {
// 				decodedLogs = append(decodedLogs, decoded)
// 			}
// 		}

// 		// Here you would store or attach decodedLogs to your transaction struct.
// 		// For example, if you have a custom Transactions struct:
// 		txInfo.Logs = decodedLogs

// 		log.Printf("Decoded logs for tx %s: %+v", tx.Hash().Hex(), decodedLogs)
// 	}
// 	return nil
// }

// // LogByLog processes a single log entry and decodes it using the loaded ABIs.
// func (e *Ethereum) LogByLog(logEntry *types.Log, tx *types.Transaction) (map[string]interface{}, error) {
// 	contractAddr := logEntry.Address.Hex()
// 	parsedABI, ok := e.ABIs[contractAddr]
// 	if !ok {
// 		log.Printf("No ABI found for contract %s, raw log: %+v", contractAddr, logEntry)
// 		return nil, nil
// 	}

// 	if len(logEntry.Topics) == 0 {
// 		log.Printf("No topics in log for tx %s", tx.Hash().Hex())
// 		return nil, nil
// 	}

// 	event, err := parsedABI.EventByID(logEntry.Topics[0])
// 	if err != nil {
// 		log.Printf("No matching event for topic %s in contract %s: %v", logEntry.Topics[0].Hex(), contractAddr, err)
// 		return nil, nil
// 	}

// 	// Decode event fields
// 	dataMap := make(map[string]interface{})
// 	err = parsedABI.UnpackIntoMap(dataMap, event.Name, logEntry.Data)
// 	if err != nil {
// 		log.Printf("Failed to unpack log data for event %s: %v", event.Name, err)
// 		return nil, nil
// 	}

// 	// Decode indexed fields from topics
// 	topicIndex := 1 // Topic 0 is the event signature
// 	for _, input := range event.Inputs {
// 		if input.Indexed {
// 			if len(logEntry.Topics) > topicIndex {
// 				arg, err := abi.ParseTopics([]interface{}{input.Type}, []common.Hash{logEntry.Topics[topicIndex]})
// 				if err == nil && len(arg) > 0 {
// 					dataMap[input.Name] = arg[0]
// 				}
// 			}
// 			topicIndex++
// 		}
// 	}

// 	return dataMap, nil
// }

// func (e *Ethereum) StoreBlock(blockNumber int64, block *Raw) error {
// 	blockJSON, err := json.Marshal(block)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal block: %w", err)
// 	}
// 	err = e.Factory.Data.RB.SAdd(e.Factory.Ctx, "ethereum", blockJSON).Err()
// 	if err != nil {
// 		return fmt.Errorf("failed to store block in Redis hash: %w", err)
// 	}
// 	return nil
// }
