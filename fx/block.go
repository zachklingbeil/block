package fx

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

func (fx *Fx) Block(number *big.Int) (*types.Block, error) {
	block, err := fx.Eth.BlockByNumber(fx.Context, number)
	if err != nil {
		return nil, fmt.Errorf("block: %w", err)
	}
	return block, nil
}
