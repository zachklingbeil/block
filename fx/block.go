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

func (fx *Fx) Fetch(number *big.Int) (*types.Header, []*types.Transaction, []*types.Receipt, error) {
	tag := "latest"
	if number != nil {
		tag = hexutil.EncodeBig(number)
	}

	var raw json.RawMessage
	var receipts []*types.Receipt

	batch := []rpc.BatchElem{
		{Method: "eth_getBlockByNumber", Args: []any{tag, true}, Result: &raw},
		{Method: "eth_getBlockReceipts", Args: []any{tag}, Result: &receipts},
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
		txs[i] = Transaction{
			From:    from,
			Tx:      tx,
			Receipt: receipts[i],
		}
	}

	return &Block{
		Header:       header,
		Transactions: txs,
		Signer:       signer,
	}, nil
}
