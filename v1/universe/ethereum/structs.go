package ethereum

// type Transaction struct {
// 	BlockNumber      uint64      `json:"block"`
// 	Timestamp        uint64      `json:"timestamp"`
// 	TransactionIndex uint        `json:"transactionIndex"`
// 	Hash             string      `json:"hash"`
// 	From             string      `json:"from"`
// 	To               string      `json:"to"`
// 	Value            string      `json:"value,omitempty"`
// 	Gas              uint64      `json:"gas"`
// 	GasUsed          uint64      `json:"gasUsed"`
// 	Inputs           string      `json:"inputs,omitempty"`
// 	GasPrice         string      `json:"gasPrice"`
// 	Signature        string      `json:"signature,omitempty"`
// 	LogsIn           []LogTopics `json:"logsIn,omitempty"`
// 	Logs             []LogTopics `json:"logs,omitempty"`
// }

// type LogTopics struct {
// 	Topic0 string `json:"topic0,omitempty"`
// 	Topic1 string `json:"topic1,omitempty"`
// 	Topic2 string `json:"topic2,omitempty"`
// 	Topic3 string `json:"topic3,omitempty"`
// 	Topic4 string `json:"topic4,omitempty"`
// 	From   string `json:"from,omitempty"`
// 	To     string `json:"to,omitempty"`
// 	Value  string `json:"value,omitempty"`
// }

// // BlockByBlock fetches and processes the latest block, storing all decoded transactions.
// func (e *Ethereum) BlockByBlock() error {
// 	block, err := e.Factory.Eth.BlockByNumber(e.Factory.Ctx, e.Header)
// 	if err != nil {
// 		log.Printf("Error fetching block %d: %v", e.Header.Uint64(), err)
// 		return err
// 	}

// 	var txs []Transaction
// 	for idx, tx := range block.Transactions() {
// 		txs = append(txs, e.ProcessTransaction(tx, block, uint(idx)))
// 	}

// 	// Store or use txs as needed
// 	if err := e.StoreTransactions(txs); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (e *Ethereum) StoreTransactions(txs []Transaction) error {
// 	txsJSON, err := json.Marshal(txs)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal transactions: %w", err)
// 	}
// 	err = e.Factory.Data.RB.SAdd(e.Factory.Ctx, "ethereum:txs", txsJSON).Err()
// 	if err != nil {
// 		return fmt.Errorf("failed to store transactions in Redis: %w", err)
// 	}
// 	return nil
// }

// func (e *Ethereum) ProcessTransaction(tx *types.Transaction, block *types.Block, txIndex uint) Transaction {
// 	signer := e.Signer(block.Number(), block.Time())
// 	from := ""
// 	if addr, err := types.Sender(signer, tx); err == nil {
// 		from = strings.ToLower(addr.Hex())
// 	}
// 	to := ""
// 	if tx.To() != nil {
// 		to = strings.ToLower(tx.To().Hex())
// 	}

// 	// input := common.Bytes2Hex(tx.Data())
// 	// if tx.To() != nil && len(tx.Data()) >= 4 {
// 	// 	methodID := common.Bytes2Hex(tx.Data()[:4])
// 	// 	if methodName, ok := e.GetSignature(methodID); ok {
// 	// 		input = methodName
// 	// 	} else if parsedABI, ok := e.ABIs[tx.To().Hex()]; ok {
// 	// 		if method, err := parsedABI.MethodById(tx.Data()[:4]); err == nil {
// 	// 			args := make(map[string]any)
// 	// 			_ = method.Inputs.UnpackIntoMap(args, tx.Data()[4:])
// 	// 			if len(args) > 0 {
// 	// 				if b, err := json.Marshal(args); err == nil {
// 	// 					input = string(b)
// 	// 				}
// 	// 			}
// 	// 		}
// 	// 	}
// 	// }

// 	receipt, _ := e.Factory.Eth.TransactionReceipt(e.Factory.Ctx, tx.Hash())

// 	var logsDecoded []LogTopics
// 	if receipt != nil {
// 		for _, l := range receipt.Logs {
// 			logsDecoded = append(logsDecoded, e.decodeLogTopics(l))
// 		}
// 	}

// 	value := tx.Value().String()
// 	if value == "0" {
// 		value = ""
// 	}

// 	var signature string
// 	if tx.Data() != nil && len(tx.Data()) >= 4 {
// 		sig := fmt.Sprintf("0x%x", tx.Data()[:4])
// 		if sigText, ok := e.GetSignature(sig); ok {
// 			signature = sigText
// 		}
// 	}

// 	return Transaction{
// 		BlockNumber:      block.NumberU64(),
// 		Timestamp:        block.Time(),
// 		TransactionIndex: txIndex,
// 		Hash:             tx.Hash().Hex(),
// 		From:             from,
// 		To:               to,
// 		Value:            value,
// 		Gas:              tx.Gas(),
// 		GasUsed:          receipt.GasUsed,
// 		GasPrice:         tx.GasPrice().String(),
// 		Signature:        signature,
// 		Logs:             logsDecoded,
// 		// Inputs:           input,
// 	}
// }

// func (e *Ethereum) decodeLogTopics(l *types.Log) LogTopics {
// 	address := strings.ToLower(l.Address.Hex())
// 	var decoded LogTopics

// 	// If ABI is available, decode event and all inputs (indexed and non-indexed)
// 	if parsedABI, ok := e.ABIs[address]; ok && len(l.Topics) > 0 {
// 		if event, err := parsedABI.EventByID(l.Topics[0]); err == nil {
// 			decoded.Topic0 = event.Name
// 			topicIndex := 1
// 			// Prepare a map for non-indexed data
// 			dataMap := make(map[string]interface{})
// 			// Unpack non-indexed data
// 			if err := parsedABI.UnpackIntoMap(dataMap, event.Name, l.Data); err == nil {
// 				// Copy non-indexed fields to decoded struct
// 				if v, ok := dataMap["from"].(common.Address); ok {
// 					decoded.From = v.Hex()
// 				}
// 				if v, ok := dataMap["to"].(common.Address); ok {
// 					decoded.To = v.Hex()
// 				}
// 				if v, ok := dataMap["value"]; ok {
// 					decoded.Value = fmt.Sprintf("%v", v)
// 				}
// 			}
// 			// Now decode indexed topics
// 			for i, input := range event.Inputs {
// 				if input.Indexed && len(l.Topics) > topicIndex {
// 					switch input.Type.String() {
// 					case "address":
// 						val := "0x" + l.Topics[topicIndex].Hex()[26:]
// 						if input.Name == "from" {
// 							decoded.From = val
// 						} else if input.Name == "to" {
// 							decoded.To = val
// 						} else {
// 							switch i {
// 							case 0:
// 								decoded.Topic1 = val
// 							case 1:
// 								decoded.Topic2 = val
// 							case 2:
// 								decoded.Topic3 = val
// 							case 3:
// 								decoded.Topic4 = val
// 							}
// 						}
// 					case "uint256", "uint8", "uint":
// 						val := l.Topics[topicIndex].Big().String()
// 						if input.Name == "value" {
// 							decoded.Value = val
// 						} else {
// 							switch i {
// 							case 0:
// 								decoded.Topic1 = val
// 							case 1:
// 								decoded.Topic2 = val
// 							case 2:
// 								decoded.Topic3 = val
// 							case 3:
// 								decoded.Topic4 = val
// 							}
// 						}
// 					default:
// 						// fallback to hex
// 						switch i {
// 						case 0:
// 							decoded.Topic1 = l.Topics[topicIndex].Hex()
// 						case 1:
// 							decoded.Topic2 = l.Topics[topicIndex].Hex()
// 						case 2:
// 							decoded.Topic3 = l.Topics[topicIndex].Hex()
// 						case 3:
// 							decoded.Topic4 = l.Topics[topicIndex].Hex()
// 						}
// 					}
// 					topicIndex++
// 				}
// 			}
// 			return decoded
// 		}
// 	}

// 	// If no ABI, use signature map for topic0, and hex for others
// 	if len(l.Topics) > 0 {
// 		eventSig := l.Topics[0].Hex()
// 		if eventText, ok := e.GetSignature(eventSig); ok {
// 			decoded.Topic0 = eventText
// 		} else {
// 			decoded.Topic0 = eventSig
// 		}
// 	}
// 	if len(l.Topics) > 1 {
// 		decoded.Topic1 = l.Topics[1].Hex()
// 	}
// 	if len(l.Topics) > 2 {
// 		decoded.Topic2 = l.Topics[2].Hex()
// 	}
// 	if len(l.Topics) > 3 {
// 		decoded.Topic3 = l.Topics[3].Hex()
// 	}
// 	if len(l.Topics) > 4 {
// 		decoded.Topic4 = l.Topics[4].Hex()
// 	}
// 	return decoded
// }
