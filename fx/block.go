package fx

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Block struct {
	Number       *big.Int    `json:"number"`
	Hash         common.Hash `json:"hash"`
	Timestamp    uint64      `json:"timestamp"`
	GasLimit     uint64      `json:"gasLimit"`
	GasUsed      uint64      `json:"gasUsed"`
	BaseFee      *big.Int    `json:"baseFeePerGas,omitempty"`
	Transactions []*Receipt  `json:"transactions"`
}

type Receipt struct {
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
	Logs              []*Log          `json:"logs,omitempty"`
}

type Log struct {
	Address common.Address    `json:"address"`
	Topics  *Topics           `json:"topics"`
	Data    string            `json:"data,omitempty"`
	Indexed map[string]string `json:"indexed,omitempty"`
}

type Topics struct {
	Zero  string `json:"0"`
	One   string `json:"1,omitempty"`
	Two   string `json:"2,omitempty"`
	Three string `json:"3,omitempty"`
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

	ethTxs := block.Transactions()
	signer := types.MakeSigner(fx.Chain, block.Number(), block.Time())
	txs := make([]*Receipt, len(ethTxs))
	for i, tx := range ethTxs {
		r := receipts[i]

		from, _ := types.Sender(signer, tx)

		var contractAddr *common.Address
		if r.ContractAddress != (common.Address{}) {
			contractAddr = &r.ContractAddress
		}

		txs[i] = &Receipt{
			TxHash:            tx.Hash(),
			TxIndex:           uint(i),
			From:              from,
			To:                tx.To(),
			Value:             tx.Value(),
			Input:             "0x" + hex.EncodeToString(tx.Data()),
			Status:            r.Status,
			Gas:               r.GasUsed,
			EffectiveGasPrice: r.EffectiveGasPrice,
			ContractAddress:   contractAddr,
			Logs:              fx.Logs(r.Logs),
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

func (fx *Fx) Logs(raw []*types.Log) []*Log {
	// Collect unique contract addresses
	addrSeen := make(map[common.Address]struct{})
	var addrs []common.Address
	for _, l := range raw {
		if _, ok := addrSeen[l.Address]; !ok {
			addrSeen[l.Address] = struct{}{}
			addrs = append(addrs, l.Address)
		}
	}

	// Build event ABI map per contract address
	abiMap := make(map[common.Address]map[common.Hash]*EventABI)
	for _, addr := range addrs {
		abi, err := fx.ContractABI(addr)
		if err != nil {
			continue
		}
		abiMap[addr] = ParseEvents(abi)
	}

	logs := make([]*Log, len(raw))
	for i, l := range raw {
		t := &Topics{}
		data := "0x" + hex.EncodeToString(l.Data)
		var indexed map[string]string

		if len(l.Topics) > 0 {
			resolved := false

			if events, ok := abiMap[l.Address]; ok {
				if event, ok := events[l.Topics[0]]; ok {
					resolved = true
					t.Zero = event.Sig
					indexed = DecodeIndexed(event, l.Topics)
				}
			}

			if !resolved {
				t.Zero = l.Topics[0].Hex()
				if len(l.Topics) > 1 {
					t.One = l.Topics[1].Hex()
				}
				if len(l.Topics) > 2 {
					t.Two = l.Topics[2].Hex()
				}
				if len(l.Topics) > 3 {
					t.Three = l.Topics[3].Hex()
				}
			}
		}

		logs[i] = &Log{
			Address: l.Address,
			Topics:  t,
			Indexed: indexed,
			Data:    data,
		}
	}
	return logs
}
