package fx

import (
	"encoding/json"
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Process struct {
	Factory *factory.Factory
}

type Transaction struct {
	TxType    string `json:"txType"`
	From      int64  `json:"accountId"`
	To        int64  `json:"toAccountId"`
	ToAddress string `json:"toAccountAddress"`
}

func NewProcess(factory *factory.Factory) *Process {
	return &Process{
		Factory: factory,
	}
}

func (p *Process) HelloPeers() error {
	// Retrieve the transactions column from the database
	var transactionsJSON []string
	err := p.Factory.Db.ColumnToSlice("loopring", "transactions", &transactionsJSON)
	if err != nil {
		return fmt.Errorf("failed to retrieve transactions column: %w", err)
	}

	// Unmarshal the transactions and extract unique IDs
	uniqueIDs := make(map[int64]struct{})
	for _, jsonStr := range transactionsJSON {
		var txs []Transaction
		if err := json.Unmarshal([]byte(jsonStr), &txs); err != nil {
			return fmt.Errorf("failed to unmarshal transactions: %w", err)
		}

		// Add unique IDs from the transactions
		for _, tx := range txs {
			if tx.From != 0 {
				uniqueIDs[tx.From] = struct{}{}
			}
			if tx.To != 0 {
				uniqueIDs[tx.To] = struct{}{}
			}
		}
	}

	// Get the total number of unique IDs
	count := len(uniqueIDs)
	fmt.Printf("%d unique IDs...\n", count)

	// Call HelloUniverse for each unique ID and log the countdown
	for id := range uniqueIDs {
		idStr := fmt.Sprintf("%d", id)
		p.Factory.Peer.HelloUniverse(idStr)
		count--
		fmt.Printf("%d\n", count)
	}
	fmt.Println("HelloUniverse.")
	return nil
}
