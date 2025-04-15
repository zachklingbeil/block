package process

import (
	"encoding/json"
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Process struct {
	Factory *factory.Factory
	RawTxs  []any
	Types   *Types
	Txs     []Tx
	Map     map[*Coordinate]*Tx
	Counts  map[string]int
}

type Types struct {
	Deposit       []DW            `json:"deposit,omitempty"`
	Withdrawal    []DW            `json:"withdraw,omitempty"`
	Swaps         []Swap          `json:"swap,omitempty"`
	Transfers     []Transfer      `json:"transfer,omitempty"`
	Mints         []Mint          `json:"mint,omitempty"`
	AccountUpdate []AccountUpdate `json:"accountUpdate,omitempty"`
	AmmUpdate     []AmmUpdate     `json:"ammUpdate,omitempty"`
	NftData       []NftData       `json:"nftData,omitempty"`
	TBD           []any           `json:"tbd,omitempty"`
	*json.RawMessage
}

func InitProcess(factory *factory.Factory) *Process {
	qtx := 10000

	process := &Process{
		Factory: factory,
		Txs:     make([]Tx, 0, qtx),
		Counts:  make(map[string]int),
		Map:     make(map[*Coordinate]*Tx),
		RawTxs:  make([]any, 0, qtx),
		Types: &Types{
			Deposit:       make([]DW, 0, qtx),
			Withdrawal:    make([]DW, 0, qtx),
			Swaps:         make([]Swap, 0, qtx),
			Transfers:     make([]Transfer, 0, qtx),
			Mints:         make([]Mint, 0, qtx),
			NftData:       make([]NftData, 0, qtx),
			AmmUpdate:     make([]AmmUpdate, 0, qtx),
			AccountUpdate: make([]AccountUpdate, 0, qtx),
			TBD:           make([]any, 0, qtx),
		},
	}
	// if err := process.CreateTxTable(); err != nil {
	// 	fmt.Printf("Warning: failed to create transactions table: %v\n", err)
	// }
	if err := process.LoadRecentBlocks(500); err != nil {
		fmt.Printf("Warning: failed to load blocks: %v\n", err)
	}
	return process
}

func (p *Process) LoadRecentBlocks(limit int) error {
	query := `
        SELECT block, tx
        FROM loopring
        ORDER BY block DESC
        LIMIT $1;
    `
	rows, err := p.Factory.Db.Query(query, limit)
	if err != nil {
		return fmt.Errorf("failed to query loopring table: %w", err)
	}
	defer rows.Close()

	var rawTxs []any
	for rows.Next() {
		var txArray []byte
		if err := rows.Scan(new(int64), &txArray); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		var transactions []json.RawMessage
		if err := json.Unmarshal(txArray, &transactions); err != nil {
			return fmt.Errorf("failed to unmarshal transactions array: %w", err)
		}

		for _, tx := range transactions {
			rawTxs = append(rawTxs, tx)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over rows: %w", err)
	}
	p.RawTxs = rawTxs
	p.Counts["Input"] = len(p.RawTxs)
	return nil
}

func (p *Process) ProcessTransactions() error {
	for _, rawTx := range p.RawTxs {
		txJSON, ok := rawTx.(json.RawMessage)
		if !ok {
			fmt.Printf("Skipping invalid transaction format: %v\n", rawTx)
			continue
		}

		if err := p.processTransaction(txJSON); err != nil {
			fmt.Printf("Error processing transaction: %v\n", err)
			continue
		}
	}

	p.ConvertTypesToTxs()
	p.Counts["Txs"] = len(p.Txs)
	p.Factory.Json.Print(p.Counts)
	return nil
}

func (p *Process) ConvertTypesToTxs() {
	for _, deposit := range p.Types.Deposit {
		p.Txs = append(p.Txs, p.DepositToTx(deposit))
	}

	for _, withdrawal := range p.Types.Withdrawal {
		p.Txs = append(p.Txs, p.WithdrawToTx(withdrawal))
	}

	for _, swap := range p.Types.Swaps {
		p.Txs = append(p.Txs, p.SwapToTx(swap))
	}

	for _, transfer := range p.Types.Transfers {
		p.Txs = append(p.Txs, p.TransferToTx(transfer))
	}

	for _, mint := range p.Types.Mints {
		p.Txs = append(p.Txs, p.MintToTx(mint))
	}

	for _, accountUpdate := range p.Types.AccountUpdate {
		p.Txs = append(p.Txs, p.AccountUpdateToTx(accountUpdate))
	}

	for _, ammUpdate := range p.Types.AmmUpdate {
		p.Txs = append(p.Txs, p.AmmUpdateToTx(ammUpdate))
	}

	for _, nftData := range p.Types.NftData {
		p.Txs = append(p.Txs, p.NftDataToTx(nftData))
	}
}

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

	// Print an example transaction for each type
	fmt.Println("Example transactions by type:")
	for txType, tx := range exampleTxs {
		txJSON, err := json.MarshalIndent(tx, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling transaction to JSON for type %s: %v\n", txType, err)
			continue
		}
		fmt.Printf("Example %s transaction:\n%s\n\n", txType, string(txJSON))
	}
}

func (p *Process) PopulateTxMap() {
	for i := range p.Txs {
		tx := p.Txs[i] // Get the transaction

		txWithoutCoordinates := tx
		txWithoutCoordinates.Coordinates = Coordinate{} // Clear the Coordinates field

		p.Map[&tx.Coordinates] = &txWithoutCoordinates
	}
}
