package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/zachklingbeil/factory"
)

type Ethereum struct {
	Factory *factory.Factory
	Chain   *params.ChainConfig

	cancel context.CancelFunc
}

type Block struct {
	*types.Block
	Receipts []*types.Receipt `json:"receipts"`
}

func New(factory *factory.Factory) (*Ethereum, error) {
	ethereum := &Ethereum{
		Factory: factory,
		Chain:   params.MainnetChainConfig,
	}

	// Get latest block number
	latest, err := ethereum.GetLatestBlockNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block number: %w", err)
	}

	opts := DecodeOptions{
		StartBlock: latest,        // Start from latest block
		EndBlock:   big.NewInt(0), // Go back to genesis
		ProgressCallback: func(blockNum *big.Int, total *big.Int) {
			fmt.Printf("%s\n", blockNum.String())
		},
		ErrorCallback: func(blockNum *big.Int, err error) {
			fmt.Printf("Error processing block %s: %v\n", blockNum.String(), err)
		},
	}

	ethereum.DecodeBlockRange(opts)
	return ethereum, nil
}
