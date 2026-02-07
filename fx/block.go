package fx

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Block struct {
	Number       *big.Int      `json:"number"`
	Hash         common.Hash   `json:"hash"`
	ParentHash   common.Hash   `json:"parentHash"`
	Timestamp    uint64        `json:"timestamp"`
	GasLimit     uint64        `json:"gasLimit"`
	GasUsed      uint64        `json:"gasUsed"`
	BaseFee      *big.Int      `json:"baseFeePerGas,omitempty"`
	Transactions []common.Hash `json:"transactions"`
}

func (fx *Fx) Block(number *big.Int) (*Block, error) {
	block, err := fx.Eth.BlockByNumber(fx.Context, number)
	if err != nil {
		return nil, fmt.Errorf("block: %w", err)
	}

	txs := make([]common.Hash, len(block.Transactions()))
	for i, tx := range block.Transactions() {
		txs[i] = tx.Hash()
	}

	return &Block{
		Number:       block.Number(),
		Hash:         block.Hash(),
		ParentHash:   block.ParentHash(),
		Timestamp:    block.Time(),
		GasLimit:     block.GasLimit(),
		GasUsed:      block.GasUsed(),
		BaseFee:      block.BaseFee(),
		Transactions: txs,
	}, nil
}
