package ethereum

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Transactions struct {
	From              string     `json:"from,omitempty"`
	To                string     `json:"to,omitempty"`
	Value             *big.Int   `json:"value,omitempty"`
	Gas               uint64     `json:"gas,omitempty"`
	GasPrice          *big.Int   `json:"gasPrice,omitempty"`
	Nonce             uint64     `json:"nonce,omitempty"`
	CumulativeGasUsed uint64     `json:"cumulativeGasUsed,omitempty"`
	FunctionSignature string     `json:"functionSignature,omitempty"`
	Logs              []*LogInfo `json:"logs,omitempty"`
}

type LogInfo struct {
	Address    string         `json:"address,omitempty"`
	Topics     []string       `json:"topics,omitempty"`
	DataLength int            `json:"dataLength,omitempty"`
	EventType  string         `json:"eventType,omitempty"`
	Fields     map[string]any `json:"fields,omitempty"`
}

const (
	transferEvent           = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	transfer1155SingleEvent = "0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62"
	transfer1155BatchEvent  = "0x4a39dc06d4c0dbc64b70b1b5fdcf9a43c3b840ecb9c7aafb5c62c0124c6a16e3"
)

func (e *Ethereum) populateReceiptInfo(txInfo *Transactions, receipt *types.Receipt) {
	txInfo.CumulativeGasUsed = receipt.CumulativeGasUsed
	for _, log := range receipt.Logs {
		if logInfo := e.ParseLogSafe(log); logInfo != nil {
			txInfo.Logs = append(txInfo.Logs, logInfo)
		}
	}
}

func (e *Ethereum) ParseLogSafe(log *types.Log) *LogInfo {
	if len(log.Topics) == 0 {
		return nil
	}
	sighash := log.Topics[0].Hex()
	token := e.Zero.Source(toLowerHex(log.Address.Hex()))
	var resolvedAddr string
	if token != nil && token.Token != "" {
		resolvedAddr = token.Token
	} else {
		resolvedAddr = toLowerHex(log.Address.Hex())
	}
	// ERC20/ERC721 Transfer
	if len(log.Topics) == 3 && sighash == transferEvent && len(log.Data) == 32 {
		return &LogInfo{
			Address:   resolvedAddr,
			EventType: "ERC20/ERC721 Transfer",
			Fields: map[string]any{
				"from":  extractAddr(log.Topics[1]),
				"to":    extractAddr(log.Topics[2]),
				"value": bigIntFromBytes(log.Data),
			},
		}
	}
	// ERC1155 TransferSingle
	if len(log.Topics) == 4 && sighash == transfer1155SingleEvent && len(log.Data) == 64 {
		return &LogInfo{
			Address:   toLowerHex(log.Address.Hex()),
			EventType: "ERC1155 TransferSingle",
			Fields: map[string]any{
				"operator": extractAddr(log.Topics[1]),
				"from":     extractAddr(log.Topics[2]),
				"to":       extractAddr(log.Topics[3]),
				"id":       bigIntFromBytes(log.Data[:32]),
				"value":    bigIntFromBytes(log.Data[32:]),
			},
		}
	}

	// ERC1155 TransferBatch
	if len(log.Topics) == 4 && sighash == transfer1155BatchEvent && len(log.Data) >= 64 {
		ids, values := e.decode1155Batch(log.Data)
		return &LogInfo{
			Address:   toLowerHex(log.Address.Hex()),
			EventType: "ERC1155 TransferBatch",
			Fields: map[string]any{
				"operator": extractAddr(log.Topics[1]),
				"from":     extractAddr(log.Topics[2]),
				"to":       extractAddr(log.Topics[3]),
				"ids":      ids,
				"values":   values,
			},
		}
	}
	// Generic ABI-based decoding
	event, ok := e.EventABI[sighash]
	if !ok {
		return nil
	}

	fields := make(map[string]any, len(event.Inputs))
	idx := 1
	for _, arg := range event.Inputs {
		fields[arg.Name] = nil
		if arg.Indexed && len(log.Topics) > idx {
			switch arg.Type.String() {
			case "address":
				fields[arg.Name] = common.HexToAddress(log.Topics[idx].Hex()).Hex()
			case "uint256", "uint":
				fields[arg.Name] = new(big.Int).SetBytes(log.Topics[idx].Bytes())
			case "bool":
				fields[arg.Name] = log.Topics[idx].Big().Cmp(big.NewInt(0)) != 0
			default:
				fields[arg.Name] = log.Topics[idx].Hex()
			}
			idx++
		}
	}

	nonIndexed := event.Inputs.NonIndexed()
	if len(nonIndexed) > 0 && len(log.Data) > 0 {
		if values, err := nonIndexed.Unpack(log.Data); err == nil {
			for i, arg := range nonIndexed {
				switch arg.Type.String() {
				case "address":
					fields[arg.Name] = values[i].(common.Address).Hex()
				case "uint256", "uint":
					fields[arg.Name] = values[i].(*big.Int)
				case "bool":
					fields[arg.Name] = values[i].(bool)
				default:
					fields[arg.Name] = values[i]
				}
			}
		}
	}
	logInfo := &LogInfo{
		Address:   toLowerHex(log.Address.Hex()),
		EventType: event.Name,
		Fields:    fields,
	}

	// Set EventType to signature if available
	if sig, ok := e.GetEventSignature(sighash); ok {
		logInfo.EventType = sig
	}

	return logInfo
}

func toLowerHex(s string) string           { return strings.ToLower(s) }
func extractAddr(topic common.Hash) string { return toLowerHex("0x" + topic.Hex()[26:]) }
func bigIntFromBytes(b []byte) *big.Int    { return new(big.Int).SetBytes(b) }

func (e *Ethereum) decode1155Batch(data []byte) ([]*big.Int, []*big.Int) {
	if len(data) < 128 {
		return nil, nil
	}
	idsOffset := int(new(big.Int).SetBytes(data[:32]).Int64())
	valuesOffset := int(new(big.Int).SetBytes(data[32:64]).Int64())
	ids := e.decodeBigIntArray(data, idsOffset)
	values := e.decodeBigIntArray(data, valuesOffset)
	return ids, values
}

func (e *Ethereum) decodeBigIntArray(data []byte, offset int) []*big.Int {
	length := int(new(big.Int).SetBytes(data[offset : offset+32]).Int64())
	result := make([]*big.Int, length)
	for i := range length {
		start := offset + 32 + i*32
		end := start + 32
		result[i] = new(big.Int).SetBytes(data[start:end])
	}
	return result
}
