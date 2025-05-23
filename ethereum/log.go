package ethereum

// type LogInfo struct {
// 	Address    string         `json:"address,omitempty"`
// 	Topics     []string       `json:"topics,omitempty"`
// 	DataLength int            `json:"dataLength,omitempty"`
// 	EventType  string         `json:"eventType,omitempty"`
// 	Zero       string         `json:"zero,omitempty"`
// 	One        string         `json:"one,omitempty"`
// 	Value      *big.Int       `json:"value,omitempty"`
// 	Fields     map[string]any `json:"fields,omitempty"`
// }

// const (
// 	transferEvent           = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
// 	transfer1155SingleEvent = "0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62"
// 	transfer1155BatchEvent  = "0x4a39dc06d4c0dbc64b70b1b5fdcf9a43c3b840ecb9c7aafb5c62c0124c6a16e3"
// )

// func (e *Ethereum) ParseLogSafe(log *types.Log) *LogInfo {
// 	if len(log.Topics) == 0 {
// 		return nil
// 	}
// 	sighash := log.Topics[0].Hex()

// 	// ERC1155 special handling
// 	if len(log.Topics) == 4 && sighash == transfer1155SingleEvent && len(log.Data) == 64 {
// 		return e.parseERC1155Transfer(log)
// 	}

// 	// Use ABI for all other events, including ERC20/ERC721
// 	return e.parseEventWithABI(log)
// }

// func (e *Ethereum) parseEventWithABI(log *types.Log) *LogInfo {
// 	if len(log.Topics) == 0 {
// 		return nil
// 	}
// 	sighash := log.Topics[0].Hex()
// 	event, ok := e.EventABI[sighash]
// 	if !ok {
// 		return nil
// 	}

// 	fields := make(map[string]any)
// 	indexedIdx := 1

// 	// Decode indexed fields from topics
// 	for _, input := range event.Inputs {
// 		if input.Indexed {
// 			if indexedIdx >= len(log.Topics) {
// 				break
// 			}
// 			switch input.Type.String() {
// 			case "address":
// 				fields[input.Name] = common.HexToAddress(log.Topics[indexedIdx].Hex()).Hex()
// 			case "uint256", "uint":
// 				fields[input.Name] = new(big.Int).SetBytes(log.Topics[indexedIdx].Bytes())
// 			default:
// 				fields[input.Name] = log.Topics[indexedIdx].Hex()
// 			}
// 			indexedIdx++
// 		}
// 	}

// 	// Decode non-indexed fields from data
// 	nonIndexed := event.Inputs.NonIndexed()
// 	if len(nonIndexed) > 0 && len(log.Data) > 0 {
// 		values, err := nonIndexed.Unpack(log.Data)
// 		if err == nil {
// 			for i, arg := range nonIndexed {
// 				switch arg.Type.String() {
// 				case "address":
// 					if v, ok := values[i].(common.Address); ok {
// 						fields[arg.Name] = v.Hex()
// 					} else {
// 						fields[arg.Name] = values[i]
// 					}
// 				case "uint256", "uint":
// 					if v, ok := values[i].(*big.Int); ok {
// 						fields[arg.Name] = v
// 					} else {
// 						fields[arg.Name] = values[i]
// 					}
// 				default:
// 					fields[arg.Name] = values[i]
// 				}
// 			}
// 		}
// 	}

// 	return &LogInfo{
// 		Address: toLowerHex(log.Address.Hex()),
// 		Fields:  fields,
// 	}
// }

// // // decodeLogFields extracts indexed and non-indexed fields from the log based on the event ABI.
// // func (e *Ethereum) decodeLogFields(event abi.Event, log *types.Log) map[string]any {
// // 	fields := make(map[string]any, len(event.Inputs))
// // 	idx := 1
// // 	for _, arg := range event.Inputs {
// // 		fields[arg.Name] = nil
// // 		if arg.Indexed && idx < len(log.Topics) {
// // 			switch arg.Type.String() {
// // 			case "address":
// // 				fields[arg.Name] = common.HexToAddress(log.Topics[idx].Hex()).Hex()
// // 			case "uint256", "uint":
// // 				fields[arg.Name] = new(big.Int).SetBytes(log.Topics[idx].Bytes())
// // 			case "bool":
// // 				fields[arg.Name] = log.Topics[idx].Big().Cmp(big.NewInt(0)) != 0
// // 			default:
// // 				fields[arg.Name] = log.Topics[idx].Hex()
// // 			}
// // 			idx++
// // 		}
// // 	}

// // 	if nonIndexed := event.Inputs.NonIndexed(); len(nonIndexed) > 0 && len(log.Data) > 0 {
// // 		if values, err := nonIndexed.Unpack(log.Data); err == nil {
// // 			for i, arg := range nonIndexed {
// // 				switch arg.Type.String() {
// // 				case "address":
// // 					if v, ok := values[i].(common.Address); ok {
// // 						fields[arg.Name] = v.Hex()
// // 					} else {
// // 						fields[arg.Name] = values[i]
// // 					}
// // 				case "uint256", "uint":
// // 					if v, ok := values[i].(*big.Int); ok {
// // 						fields[arg.Name] = v
// // 					} else {
// // 						fields[arg.Name] = values[i]
// // 					}
// // 				case "bool":
// // 					if v, ok := values[i].(bool); ok {
// // 						fields[arg.Name] = v
// // 					} else {
// // 						fields[arg.Name] = values[i]
// // 					}
// // 				default:
// // 					fields[arg.Name] = values[i]
// // 				}
// // 			}
// // 		}
// // 	}
// // 	return fields
// // }

// // func (e *Ethereum) parseERC20ERC721Transfer(log *types.Log) *LogInfo {
// // 	token := e.Zero.Source(toLowerHex(log.Address.Hex()))
// // 	var resolvedAddr string
// // 	if token != nil && token.Token != "" {
// // 		resolvedAddr = token.Token
// // 	} else {
// // 		resolvedAddr = toLowerHex(log.Address.Hex())
// // 		token = e.Zero.FetchERC20(common.HexToAddress(resolvedAddr))
// // 		if token != nil && token.Token != "" {
// // 			resolvedAddr = token.Token
// // 		}
// // 	}
// // 	return &LogInfo{
// // 		Address: resolvedAddr,
// // 		Zero:    extractAddr(log.Topics[1]),
// // 		One:     extractAddr(log.Topics[2]),
// // 		Value:   bigIntFromBytes(log.Data),
// // 	}
// // }

// func (e *Ethereum) parseERC1155Transfer(log *types.Log) *LogInfo {
// 	if len(log.Topics) == 4 && len(log.Data) == 64 {
// 		return &LogInfo{
// 			Address:   toLowerHex(log.Address.Hex()),
// 			EventType: "ERC1155 TransferSingle",
// 			Fields: map[string]any{
// 				"operator": extractAddr(log.Topics[1]),
// 				"from":     extractAddr(log.Topics[2]),
// 				"to":       extractAddr(log.Topics[3]),
// 				"id":       bigIntFromBytes(log.Data[:32]),
// 				"value":    bigIntFromBytes(log.Data[32:]),
// 			},
// 		}
// 	}

// 	// ERC1155 TransferBatch
// 	if len(log.Topics) == 4 && len(log.Data) >= 64 {
// 		ids, values := e.decode1155Batch(log.Data)
// 		return &LogInfo{
// 			Address:   toLowerHex(log.Address.Hex()),
// 			EventType: "ERC1155 TransferBatch",
// 			Fields: map[string]any{
// 				"operator": extractAddr(log.Topics[1]),
// 				"from":     extractAddr(log.Topics[2]),
// 				"to":       extractAddr(log.Topics[3]),
// 				"ids":      ids,
// 				"values":   values,
// 			},
// 		}
// 	}

// 	return nil
// }

// func toLowerHex(s string) string           { return strings.ToLower(s) }
// func extractAddr(topic common.Hash) string { return toLowerHex("0x" + topic.Hex()[26:]) }
// func bigIntFromBytes(b []byte) *big.Int    { return new(big.Int).SetBytes(b) }

// func (e *Ethereum) decode1155Batch(data []byte) ([]*big.Int, []*big.Int) {
// 	if len(data) < 128 {
// 		return nil, nil
// 	}
// 	idsOffset := int(new(big.Int).SetBytes(data[:32]).Int64())
// 	valuesOffset := int(new(big.Int).SetBytes(data[32:64]).Int64())
// 	ids := e.decodeBigIntArray(data, idsOffset)
// 	values := e.decodeBigIntArray(data, valuesOffset)
// 	return ids, values
// }

// func (e *Ethereum) decodeBigIntArray(data []byte, offset int) []*big.Int {
// 	length := int(new(big.Int).SetBytes(data[offset : offset+32]).Int64())
// 	result := make([]*big.Int, length)
// 	for i := range length {
// 		start := offset + 32 + i*32
// 		end := start + 32
// 		result[i] = new(big.Int).SetBytes(data[start:end])
// 	}
// 	return result
// }
