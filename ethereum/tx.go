package ethereum

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	transferEvent           = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	transfer1155SingleEvent = "0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62"
	transfer1155BatchEvent  = "0x4a39dc06d4c0dbc64b70b1b5fdcf9a43c3b840ecb9c7aafb5c62c0124c6a16e3"
)

// Raw holds information about a block.
type Raw struct {
	Number       uint64
	Hash         string
	ParentHash   string
	Time         uint64
	GasUsed      uint64
	GasLimit     uint64
	BaseFee      *big.Int
	Transactions []*Transactions
}

// Transactions holds information about a transaction.
type Transactions struct {
	Hash              string
	From              string
	To                string
	Value             *big.Int
	Gas               uint64
	GasPrice          *big.Int
	Nonce             uint64
	DataLength        int
	Type              uint8
	Status            uint64
	CumulativeGasUsed uint64
	Logs              []*LogInfo `json:"logs,omitempty"`
}

// LogInfo holds information about a transaction log.
type LogInfo struct {
	Address    string   `json:"Address,omitempty"`
	Topics     []string `json:"topics,omitempty"`
	DataLength int      `json:"dataLength,omitempty"`
	EventType  string   `json:"eventType,omitempty"`
	From       string   `json:"from,omitempty"`
	To         string   `json:"to,omitempty"`
	Value      *big.Int `json:"value,omitempty"`
	Operator   string   `json:"operator,omitempty"`
	ID         *big.Int `json:"id,omitempty"`
	IDs        []string `json:"ids,omitempty"`
	Values     []string `json:"values,omitempty"`
	RawTopics  []string `json:"rawTopics,omitempty"`
}

// processBlock processes a single block and returns its information.
func (e *Ethereum) processBlock(ctx context.Context, block *types.Block) *Raw {
	blockInfo := &Raw{
		Number:     block.NumberU64(),
		Hash:       block.Hash().Hex(),
		ParentHash: block.ParentHash().Hex(),
		Time:       block.Time(),
		GasUsed:    block.GasUsed(),
		GasLimit:   block.GasLimit(),
		BaseFee:    block.BaseFee(),
	}

	signer := e.Signer(block.Number(), block.Time())

	for _, tx := range block.Transactions() {
		txInfo := e.processTransaction(ctx, tx, signer)
		blockInfo.Transactions = append(blockInfo.Transactions, txInfo)
	}
	return blockInfo
}

// processTransaction processes a single transaction and returns its information.
func (e *Ethereum) processTransaction(ctx context.Context, tx *types.Transaction, signer types.Signer) *Transactions {
	txInfo := &Transactions{
		Hash:       strings.ToLower(tx.Hash().Hex()), // Ensure Hash is lowercase
		Value:      tx.Value(),
		Gas:        tx.Gas(),
		GasPrice:   tx.GasPrice(),
		Nonce:      tx.Nonce(),
		DataLength: len(tx.Data()),
		Type:       tx.Type(),
	}

	if addr, err := types.Sender(signer, tx); err == nil {
		txInfo.From = strings.ToLower(addr.Hex()) // Ensure From is lowercase
	}

	// Consolidate getToAddress logic here
	if tx.To() == nil {
		txInfo.To = "Contract Creation"
	} else {
		txInfo.To = strings.ToLower(tx.To().Hex()) // Ensure To is lowercase
	}

	if receipt, err := e.Factory.Eth.TransactionReceipt(ctx, tx.Hash()); err == nil {
		e.populateReceiptInfo(txInfo, receipt)
	}

	return txInfo
}

// populateReceiptInfo populates transaction information from the receipt.
func (e *Ethereum) populateReceiptInfo(txInfo *Transactions, receipt *types.Receipt) {
	txInfo.Status = receipt.Status
	txInfo.CumulativeGasUsed = receipt.CumulativeGasUsed
	for _, log := range receipt.Logs {
		if logInfo := e.parseLog(log); logInfo != nil {
			txInfo.Logs = append(txInfo.Logs, logInfo)
		}
	}
}

// parseLog handles all supported log types in one place.
func (e *Ethereum) parseLog(log *types.Log) *LogInfo {
	switch {
	case len(log.Topics) == 3 && log.Topics[0].Hex() == transferEvent && len(log.Data) == 32:
		return &LogInfo{
			Address:   toLowerHex(log.Address.Hex()),
			EventType: "ERC20/ERC721 Transfer",
			From:      extractAddr(log.Topics[1]),
			To:        extractAddr(log.Topics[2]),
			Value:     bigIntFromBytes(log.Data),
		}
	case len(log.Topics) == 4 && log.Topics[0].Hex() == transfer1155SingleEvent && len(log.Data) == 64:
		return &LogInfo{
			Address:   toLowerHex(log.Address.Hex()),
			EventType: "ERC1155 TransferSingle",
			Operator:  extractAddr(log.Topics[1]),
			From:      extractAddr(log.Topics[2]),
			To:        extractAddr(log.Topics[3]),
			ID:        bigIntFromBytes(log.Data[:32]),
			Value:     bigIntFromBytes(log.Data[32:]),
		}
	case len(log.Topics) == 4 && log.Topics[0].Hex() == transfer1155BatchEvent && len(log.Data) >= 64:
		ids, values := e.decode1155Batch(log.Data)
		return &LogInfo{
			Address:   toLowerHex(log.Address.Hex()),
			EventType: "ERC1155 TransferBatch",
			Operator:  extractAddr(log.Topics[1]),
			From:      extractAddr(log.Topics[2]),
			To:        extractAddr(log.Topics[3]),
			IDs:       ids,
			Values:    values,
		}
	default:
		return nil
	}
}

// Helper functions
func toLowerHex(s string) string           { return strings.ToLower(s) }
func extractAddr(topic common.Hash) string { return toLowerHex("0x" + topic.Hex()[26:]) }
func bigIntFromBytes(b []byte) *big.Int    { return new(big.Int).SetBytes(b) }

// decode1155Batch decodes batch transfer data into IDs and values.
func (e *Ethereum) decode1155Batch(data []byte) ([]string, []string) {
	if len(data) < 128 {
		return nil, nil
	}

	idsOffset := int(new(big.Int).SetBytes(data[:32]).Int64())
	valuesOffset := int(new(big.Int).SetBytes(data[32:64]).Int64())

	ids := e.decodeBigIntArray(data, idsOffset)
	values := e.decodeBigIntArray(data, valuesOffset)

	return ids, values
}

// decodeBigIntArray decodes an array of big integers from data at the given offset.
func (e *Ethereum) decodeBigIntArray(data []byte, offset int) []string {
	length := int(new(big.Int).SetBytes(data[offset : offset+32]).Int64())
	result := make([]string, length)

	for i := range length {
		start := offset + 32 + i*32
		end := start + 32
		result[i] = new(big.Int).SetBytes(data[start:end]).String()
	}

	return result
}
