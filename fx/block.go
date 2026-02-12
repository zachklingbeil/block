package fx

import (
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
	Function          string          `json:"function,omitempty"`
	Args              map[string]any  `json:"args,omitempty"`
	Status            uint64          `json:"status"`
	Gas               uint64          `json:"gas"`
	EffectiveGasPrice *big.Int        `json:"gasPrice"`
	ContractAddress   *common.Address `json:"contractAddress,omitempty"`
	Logs              []*Event        `json:"logs,omitempty"`
}

type Event struct {
	Address common.Address `json:"contract"`
	Name    string         `json:"name,omitempty"`
	Args    map[string]any `json:"args,omitempty"`
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

	// Collect unique contract addresses
	contracts := make(map[common.Address]struct{})
	for i, tx := range block.Transactions() {
		if tx.To() != nil {
			contracts[*tx.To()] = struct{}{}
		}
		for _, l := range receipts[i].Logs {
			contracts[l.Address] = struct{}{}
		}
	}

	// Fetch ABIs for all unique contracts
	abis := fx.fetchABIs(contracts)

	signer := types.MakeSigner(fx.Chain, block.Number(), block.Time())
	txs := make([]*Transaction, len(block.Transactions()))

	for i, tx := range block.Transactions() {
		r := receipts[i]
		from, _ := types.Sender(signer, tx)

		var contract *common.Address
		if r.ContractAddress != (common.Address{}) {
			contract = &r.ContractAddress
		}

		var fn string
		var args map[string]any
		if tx.To() != nil && len(tx.Data()) >= 4 {
			fn, args = fx.decodeInput(abis, *tx.To(), tx.Data())
		}

		logs := make([]*Event, len(r.Logs))
		for j, l := range r.Logs {
			name, decoded := fx.decodeLog(abis, l)
			logs[j] = &Event{
				Address: l.Address,
				Name:    name,
				Args:    decoded,
			}
		}

		txs[i] = &Transaction{
			TxHash:            tx.Hash(),
			TxIndex:           uint(i),
			From:              from,
			To:                tx.To(),
			Value:             tx.Value(),
			Function:          fn,
			Args:              args,
			Status:            r.Status,
			Gas:               r.GasUsed,
			EffectiveGasPrice: r.EffectiveGasPrice,
			ContractAddress:   contract,
			Logs:              logs,
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
