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
	raw, err := fx.Block(nil)
	if err != nil {
		return fmt.Errorf("Block: %w", err)
	}

	output, err := json.MarshalIndent(json.RawMessage(raw), "", "  ")
	if err != nil {
		return fmt.Errorf("MarshalIndent: %w", err)
	}

	if err := os.WriteFile("../output/block.json", output, 0644); err != nil {
		return fmt.Errorf("WriteFile: %w", err)
	}

	fmt.Println("Block written to output/block.json")
	return nil
}
