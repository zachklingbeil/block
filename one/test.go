package one

import (
	"encoding/json"
	"fmt"
	"os"
)

func (fx *One) Test() error {
	block, err := fx.Build(nil)
	if err != nil {
		return fmt.Errorf("Block: %w", err)
	}

	output, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		return fmt.Errorf("MarshalIndent: %w", err)
	}

	if err := os.WriteFile("../output/build.json", output, 0644); err != nil {
		return fmt.Errorf("WriteFile: %w", err)
	}

	fmt.Println("Block written to output/block.json")
	return nil
}
