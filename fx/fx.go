package fx

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/timefactoryio/block/zero"
)

type Fx struct {
	*zero.Zero
}

func Init(password string) *Fx {
	fx := &Fx{
		Zero: zero.Init(password),
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
