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

func Init(url string) *Fx {
	return &Fx{
		Zero: zero.Init(url),
	}
}

func (fx *Fx) Test() error {
	// Fetch block
	block, err := fx.Block(nil)
	if err != nil {
		return fmt.Errorf("Block: %w", err)
	}
	// Write raw block
	rawOutput, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		return fmt.Errorf("Marshal block: %w", err)
	}

	if err := os.WriteFile("../output/block.json", rawOutput, 0644); err != nil {
		return fmt.Errorf("WriteFile block: %w", err)
	}
	fmt.Println("Block written to output/block.json")
	return nil
}
