package fx

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/timefactoryio/block/zero"
	"golang.org/x/time/rate"
)

type Fx struct {
	*zero.Zero
	Chain     *params.ChainConfig
	Contracts map[common.Address]*Contract
	Limiter   *rate.Limiter
}

func Init(password string) *Fx {
	fx := &Fx{
		Zero:      zero.Init(password),
		Chain:     params.MainnetChainConfig,
		Contracts: make(map[common.Address]*Contract),
		Limiter:   rate.NewLimiter(rate.Limit(3), 1),
	}
	return fx
}

func (fx *Fx) Test() error {
	block, err := fx.Block(nil)
	if err != nil {
		return fmt.Errorf("Block: %w", err)
	}

	output, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		return fmt.Errorf("MarshalIndent: %w", err)
	}

	if err := os.WriteFile("../output/block.json", output, 0644); err != nil {
		return fmt.Errorf("WriteFile: %w", err)
	}

	fmt.Println("Block written to output/block.json")
	return nil
}
