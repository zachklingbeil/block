package fx

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Block struct {
	Number       *big.Int       `json:"number"`
	Hash         common.Hash    `json:"hash"`
	Timestamp    uint64         `json:"timestamp"`
	GasLimit     uint64         `json:"gasLimit"`
	GasUsed      uint64         `json:"gasUsed"`
	BaseFee      *big.Int       `json:"baseFeePerGas,omitempty"`
	Transactions []*Transaction `json:"transactions"`
}

type Transaction struct {
	TxHash            common.Hash     `json:"hash"`
	TxIndex           uint            `json:"index"`
	From              common.Address  `json:"from"`
	To                *common.Address `json:"to,omitempty"`
	Value             *big.Int        `json:"value,omitempty"`
	Input             string          `json:"input,omitempty"`
	Status            uint64          `json:"status"`
	Gas               uint64          `json:"gas"`
	EffectiveGasPrice *big.Int        `json:"gasPrice"`
	ContractAddress   *common.Address `json:"contractAddress,omitempty"`
	Logs              []*Event        `json:"logs,omitempty"`
	Decoded           []*Decoded      `json:"decoded,omitempty"`
}

type Event struct {
	Address common.Address `json:"contract"`
	Topics  []string       `json:"topics"`
	Data    string         `json:"data,omitempty"`
}

type Decoded struct {
	Contract  common.Address `json:"contract"`
	Signature string         `json:"signature"`
	Values    string         `json:"values"`
}

func (fx *Fx) Block(number *big.Int) (*Block, error) {
	block, err := fx.Eth.BlockByNumber(fx.Context, number)
	if err != nil {
		return nil, fmt.Errorf("block: %w", err)
	}

	receipts, err := fx.blockReceipts(block.Number())
	if err != nil {
		return nil, err
	}

	// Collect all contracts that appear in logs or as tx targets
	contracts := make(map[common.Address]struct{})
	for _, r := range receipts {
		for _, l := range r.Logs {
			contracts[l.Address] = struct{}{}
		}
	}
	signer := types.MakeSigner(fx.Chain, block.Number(), block.Time())
	for _, tx := range block.Transactions() {
		if tx.To() != nil {
			contracts[*tx.To()] = struct{}{}
		}
	}

	// Fetch ABIs from sourcify for all involved contracts
	abis := fx.fetchABIs(contracts)

	txs := make([]*Transaction, len(block.Transactions()))
	for i, tx := range block.Transactions() {
		r := receipts[i]
		from, _ := types.Sender(signer, tx)

		var contract *common.Address
		if r.ContractAddress != (common.Address{}) {
			contract = &r.ContractAddress
		}

		// Decode logs using ABIs
		var decoded []*Decoded
		for _, l := range r.Logs {
			name, vals := fx.decodeLog(abis, l)
			if name != "" {
				decoded = append(decoded, &Decoded{
					Contract:  l.Address,
					Signature: name,
					Values:    formatValues(vals),
				})
			}
		}

		// Decode transaction input using ABIs
		if tx.To() != nil && len(tx.Data()) >= 4 {
			name, vals := fx.decodeInput(abis, *tx.To(), tx.Data())
			if name != "" {
				decoded = append(decoded, &Decoded{
					Contract:  *tx.To(),
					Signature: name,
					Values:    formatValues(vals),
				})
			}
		}

		txs[i] = &Transaction{
			TxHash:            tx.Hash(),
			TxIndex:           uint(i),
			From:              from,
			To:                tx.To(),
			Value:             tx.Value(),
			Input:             hexEncode(tx.Data()),
			Status:            r.Status,
			Gas:               r.GasUsed,
			EffectiveGasPrice: r.EffectiveGasPrice,
			ContractAddress:   contract,
			Logs:              events(r.Logs),
			Decoded:           decoded,
		}
	}

	return &Block{
		Number:       block.Number(),
		Hash:         block.Hash(),
		Timestamp:    block.Time(),
		GasLimit:     block.GasLimit(),
		GasUsed:      block.GasUsed(),
		BaseFee:      block.BaseFee(),
		Transactions: txs,
	}, nil
}

func (fx *Fx) blockReceipts(number *big.Int) ([]*types.Receipt, error) {
	var receipts []*types.Receipt
	arg := "latest"
	if number != nil {
		arg = fmt.Sprintf("0x%x", number)
	}
	if err := fx.Rpc.CallContext(fx.Context, &receipts, "eth_getBlockReceipts", arg); err != nil {
		return nil, fmt.Errorf("block receipts: %w", err)
	}
	return receipts, nil
}

func events(raw []*types.Log) []*Event {
	out := make([]*Event, len(raw))
	for i, l := range raw {
		topics := make([]string, len(l.Topics))
		for j, t := range l.Topics {
			topics[j] = t.Hex()
		}
		out[i] = &Event{
			Address: l.Address,
			Topics:  topics,
			Data:    hexEncode(l.Data),
		}
	}
	return out
}

func formatValues(vals map[string]any) string {
	if len(vals) == 0 {
		return ""
	}
	b, err := json.Marshal(vals)
	if err != nil {
		return ""
	}
	return string(b)
}

func hexEncode(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}
