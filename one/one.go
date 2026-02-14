package one

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/timefactoryio/block/zero"
)

type One struct {
	*zero.Zero
	Chain *params.ChainConfig
}

func Init(password string) *One {
	return &One{Zero: zero.Init(password)}
}

var (
	Transfer    = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
	Approval    = crypto.Keccak256Hash([]byte("Approval(address,address,uint256)"))
	ApprovalAll = crypto.Keccak256Hash([]byte("ApprovalForAll(address,address,bool)"))
	Single      = crypto.Keccak256Hash([]byte("TransferSingle(address,address,address,uint256,uint256)"))
	Batch       = crypto.Keccak256Hash([]byte("TransferBatch(address,address,address,uint256[],uint256[])"))
)

type Block struct {
	Number       *big.Int    `json:"number"`
	Hash         common.Hash `json:"hash"`
	Timestamp    uint64      `json:"timestamp"`
	GasLimit     uint64      `json:"gasLimit"`
	GasUsed      uint64      `json:"gasUsed"`
	BaseFee      *big.Int    `json:"baseFeePerGas,omitempty"`
	Transactions []*Activity `json:"transactions"`
}

type Activity struct {
	Type        string             `json:"type"`
	Index       uint               `json:"index"`
	Transaction *types.Transaction `json:"transaction"`
	Receipt     *types.Receipt     `json:"receipt,omitempty"`
}

// Classify sets the Type field based on token standards detected in the receipt logs.
func (o *One) Classify(a *Activity) {
	if r := a.Receipt; r != nil && len(r.Logs) > 0 {
		has20, has721, has1155 := false, false, false
		for _, log := range r.Logs {
			if len(log.Topics) == 0 {
				continue
			}
			switch log.Topics[0] {
			case Single, Batch:
				has1155 = true
			case Transfer, Approval:
				if len(log.Topics) == 4 {
					has721 = true
				} else {
					has20 = true
				}
			case ApprovalAll:
				has721 = true
			}
		}
		count := 0
		if has20 {
			count++
		}
		if has721 {
			count++
		}
		if has1155 {
			count++
		}
		if count > 1 {
			a.Type = "custom"
			return
		}
		if has1155 {
			a.Type = "erc1155"
			return
		}
		if has721 {
			a.Type = "erc721"
			return
		}
		if has20 {
			a.Type = "erc20"
			return
		}
	}
	if len(a.Transaction.Data()) == 0 {
		a.Type = "eth"
		return
	}
	a.Type = "custom"
}

// Build fetches the block and receipts from the node, assembling Activity per transaction.
func (o *One) Build(number *big.Int) (*Block, error) {
	raw, err := o.Eth.BlockByNumber(o.Context, number)
	if err != nil {
		return nil, fmt.Errorf("block: %w", err)
	}

	txs := raw.Transactions()
	n := len(txs)

	var receipts []*types.Receipt
	if n > 0 {
		if err := o.Rpc.CallContext(o.Context, &receipts, "eth_getBlockReceipts", fmt.Sprintf("0x%x", raw.Number())); err != nil {
			return nil, fmt.Errorf("block receipts: %w", err)
		}
	}

	activities := make([]*Activity, n)
	for i, tx := range txs {
		a := &Activity{
			Index:       uint(i),
			Transaction: tx,
		}
		if i < len(receipts) {
			a.Receipt = receipts[i]
		}
		o.Classify(a)
		activities[i] = a
	}

	return &Block{
		Number:       raw.Number(),
		Hash:         raw.Hash(),
		Timestamp:    raw.Time(),
		GasLimit:     raw.GasLimit(),
		GasUsed:      raw.GasUsed(),
		BaseFee:      raw.BaseFee(),
		Transactions: activities,
	}, nil
}
