package process

import (
	"encoding/json"
	"fmt"
)

func (p *Process) PrintExampleTxForEachType() {
	if len(p.Txs) == 0 {
		fmt.Println("No transactions processed.")
		return
	}

	// Map to store the first transaction for each type
	exampleTxs := make(map[string]Tx)

	// Iterate over p.Txs and store the first transaction for each type
	for _, tx := range p.Txs {
		if _, exists := exampleTxs[tx.Type]; !exists {
			exampleTxs[tx.Type] = tx
		}
	}

	for txType, tx := range exampleTxs {
		txJSON, err := json.MarshalIndent(tx, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling transaction to JSON for type %s: %v\n", txType, err)
			continue
		}
		fmt.Printf("Example %s transaction:\n%s\n\n", txType, string(txJSON))
	}
}
