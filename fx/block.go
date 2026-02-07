package fx

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

type Transaction struct {
	From    common.Address
	Tx      *types.Transaction
	Receipt *types.Receipt
}

type Block struct {
	Header       *types.Header
	Transactions []Transaction
	Signer       types.Signer
}

type JSONBlock struct {
	Number       *big.Int        `json:"number"`
	Hash         common.Hash     `json:"hash"`
	ParentHash   common.Hash     `json:"parentHash"`
	Timestamp    uint64          `json:"timestamp"`
	GasUsed      uint64          `json:"gasUsed"`
	GasLimit     uint64          `json:"gasLimit"`
	BaseFee      *big.Int        `json:"baseFeePerGas,omitempty"`
	Miner        common.Address  `json:"miner"`
	Transactions json.RawMessage `json:"transactions"`
}

type JSONTransaction struct {
	Hash              common.Hash      `json:"hash"`
	Nonce             uint64           `json:"nonce"`
	From              common.Address   `json:"from"`
	To                *common.Address  `json:"to"`
	Value             *big.Int         `json:"value"`
	Input             hexutil.Bytes    `json:"input"`
	Type              uint8            `json:"type"`
	Gas               uint64           `json:"gas"`
	GasPrice          *big.Int         `json:"gasPrice"`
	MaxFeePerGas      *big.Int         `json:"maxFeePerGas,omitempty"`
	MaxPriorityFee    *big.Int         `json:"maxPriorityFeePerGas,omitempty"`
	ChainID           *big.Int         `json:"chainId,omitempty"`
	AccessList        types.AccessList `json:"accessList,omitempty"`
	BlobGas           uint64           `json:"blobGas,omitempty"`
	BlobGasFeeCap     *big.Int         `json:"maxFeePerBlobGas,omitempty"`
	BlobHashes        []common.Hash    `json:"blobVersionedHashes,omitempty"`
	V                 *big.Int         `json:"v"`
	R                 *big.Int         `json:"r"`
	S                 *big.Int         `json:"s"`
	Status            uint64           `json:"status"`
	GasUsed           uint64           `json:"gasUsed"`
	EffectiveGasPrice *big.Int         `json:"effectiveGasPrice"`
	CumulativeGasUsed uint64           `json:"cumulativeGasUsed"`
	ContractAddress   common.Address   `json:"contractAddress"`
	Logs              []*types.Log     `json:"logs"`
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	tx := t.Tx
	v, r, s := tx.RawSignatureValues()

	out := JSONTransaction{
		Hash:              tx.Hash(),
		Nonce:             tx.Nonce(),
		From:              t.From,
		To:                tx.To(),
		Value:             tx.Value(),
		Input:             tx.Data(),
		Type:              tx.Type(),
		Gas:               tx.Gas(),
		GasPrice:          tx.GasPrice(),
		ChainID:           tx.ChainId(),
		V:                 v,
		R:                 r,
		S:                 s,
		Status:            t.Receipt.Status,
		GasUsed:           t.Receipt.GasUsed,
		EffectiveGasPrice: t.Receipt.EffectiveGasPrice,
		CumulativeGasUsed: t.Receipt.CumulativeGasUsed,
		ContractAddress:   t.Receipt.ContractAddress,
		Logs:              t.Receipt.Logs,
	}

	switch {
	case tx.Type() == 3:
		out.AccessList = tx.AccessList()
		out.MaxFeePerGas = tx.GasFeeCap()
		out.MaxPriorityFee = tx.GasTipCap()
		out.BlobGas = tx.BlobGas()
		out.BlobGasFeeCap = tx.BlobGasFeeCap()
		out.BlobHashes = tx.BlobHashes()
	case tx.Type() >= 2:
		out.AccessList = tx.AccessList()
		out.MaxFeePerGas = tx.GasFeeCap()
		out.MaxPriorityFee = tx.GasTipCap()
	case tx.Type() >= 1:
		out.AccessList = tx.AccessList()
	}

	if out.Logs == nil {
		out.Logs = []*types.Log{}
	}

	return json.Marshal(out)
}

func (b *Block) MarshalJSON() ([]byte, error) {
	txs := make([]json.RawMessage, len(b.Transactions))
	for i := range b.Transactions {
		raw, err := b.Transactions[i].MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("tx %d: %w", i, err)
		}
		txs[i] = raw
	}

	txsJSON, err := json.Marshal(txs)
	if err != nil {
		return nil, err
	}

	return json.Marshal(JSONBlock{
		Number:       b.Header.Number,
		Hash:         b.Header.Hash(),
		ParentHash:   b.Header.ParentHash,
		Timestamp:    b.Header.Time,
		GasUsed:      b.Header.GasUsed,
		GasLimit:     b.Header.GasLimit,
		BaseFee:      b.Header.BaseFee,
		Miner:        b.Header.Coinbase,
		Transactions: txsJSON,
	})
}

func (fx *Fx) Fetch(number *big.Int) (*types.Header, []*types.Transaction, []*types.Receipt, error) {
	tag := "latest"
	if number != nil {
		tag = hexutil.EncodeBig(number)
	}

	var (
		raw      json.RawMessage
		receipts []*types.Receipt
	)

	batch := []rpc.BatchElem{
		{Method: "eth_getBlockByNumber", Args: []interface{}{tag, true}, Result: &raw},
		{Method: "eth_getBlockReceipts", Args: []interface{}{tag}, Result: &receipts},
	}

	if err := fx.Rpc.BatchCallContext(fx.Context, batch); err != nil {
		return nil, nil, nil, fmt.Errorf("batch: %w", err)
	}
	for _, elem := range batch {
		if elem.Error != nil {
			return nil, nil, nil, fmt.Errorf("%s: %w", elem.Method, elem.Error)
		}
	}

	var head types.Header
	if err := json.Unmarshal(raw, &head); err != nil {
		return nil, nil, nil, fmt.Errorf("header: %w", err)
	}

	var body struct {
		Transactions []*types.Transaction `json:"transactions"`
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, nil, nil, fmt.Errorf("transactions: %w", err)
	}

	return &head, body.Transactions, receipts, nil
}

func (fx *Fx) Block(number *big.Int) (*Block, error) {
	header, rawTxs, receipts, err := fx.Fetch(number)
	if err != nil {
		return nil, err
	}

	if len(rawTxs) != len(receipts) {
		return nil, fmt.Errorf("tx/receipt mismatch: %d txs, %d receipts", len(rawTxs), len(receipts))
	}

	signer := types.MakeSigner(params.MainnetChainConfig, header.Number, header.Time)

	txs := make([]Transaction, len(rawTxs))
	for i, tx := range rawTxs {
		from, _ := types.Sender(signer, tx)
		txs[i] = Transaction{From: from, Tx: tx, Receipt: receipts[i]}
	}

	return &Block{
		Header:       header,
		Transactions: txs,
		Signer:       signer,
	}, nil
}
