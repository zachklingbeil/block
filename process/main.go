package process

import (
	"encoding/json"
	"fmt"

	"github.com/zachklingbeil/factory"
)

func InitProcess(factory *factory.Factory) *Process {
	process := &Process{
		Factory: factory,
		Blocks:  []Block{},
		Txs:     &Txs{},
	}

	if err := process.LoadRecentBlocks(500); err != nil {
		fmt.Printf("Warning: failed to load blocks: %v\n", err)
	}

	process.Txs.Deposit = make([]DW, 0, 1000)
	process.Txs.Withdrawal = make([]DW, 0, 1000)
	process.Txs.Swaps = make([]Swap, 0, 10000)
	process.Txs.Transfers = make([]Transfer, 0, 10000)
	process.Txs.Mints = make([]Mint, 0, 1000)
	process.Txs.NftData = make([]NftData, 0, 1000)
	process.Txs.AmmUpdate = make([]AmmUpdate, 0, 1000)
	process.Txs.AccountUpdate = make([]AccountUpdate, 0, 1000)
	process.Txs.TBD = make([]any, 0, 10000)

	return process
}

func (p *Process) ProcessTransactions() error {
	for _, block := range p.Blocks {
		for _, tx := range block.Transactions {
			if err := p.processTransaction(tx); err != nil {
				fmt.Printf("Error processing transaction in block %d: %v\n", block.Number, err)
				continue
			}
		}
	}

	counts := map[string]int{
		"Deposits":       len(p.Txs.Deposit),
		"Withdrawals":    len(p.Txs.Withdrawal),
		"Swaps":          len(p.Txs.Swaps),
		"Transfers":      len(p.Txs.Transfers),
		"Mints":          len(p.Txs.Mints),
		"AccountUpdates": len(p.Txs.AccountUpdate),
		"AmmUpdates":     len(p.Txs.AmmUpdate),
		"NftData":        len(p.Txs.NftData),
		"Unknown":        len(p.Txs.TBD),
	}
	p.PrintTransactionExamples()
	fmt.Printf("Processed transactions: %+v\n", counts)
	return nil
}

func (p *Process) PrintTransactionExamples() {
	fmt.Println("Printing examples of each transaction type:")

	// Print an example Deposit transaction
	if len(p.Txs.Deposit) > 0 {
		fmt.Println("Example Deposit transaction:")
		p.Factory.Json.Print(p.Txs.Deposit[0])
	}

	// Print an example Withdrawal transaction
	if len(p.Txs.Withdrawal) > 0 {
		fmt.Println("Example Withdrawal transaction:")
		p.Factory.Json.Print(p.Txs.Withdrawal[0])
	}

	// Print an example Swap transaction
	if len(p.Txs.Swaps) > 0 {
		fmt.Println("Example Swap transaction:")
		p.Factory.Json.Print(p.Txs.Swaps[0])
	}

	// Print an example Transfer transaction
	if len(p.Txs.Transfers) > 0 {
		fmt.Println("Example Transfer transaction:")
		p.Factory.Json.Print(p.Txs.Transfers[0])
	}

	// Print an example Mint transaction
	if len(p.Txs.Mints) > 0 {
		fmt.Println("Example Mint transaction:")
		p.Factory.Json.Print(p.Txs.Mints[0])
	}

	// Print an example AccountUpdate transaction
	if len(p.Txs.AccountUpdate) > 0 {
		fmt.Println("Example AccountUpdate transaction:")
		p.Factory.Json.Print(p.Txs.AccountUpdate[0])
	}

	// Print an example AmmUpdate transaction
	if len(p.Txs.AmmUpdate) > 0 {
		fmt.Println("Example AmmUpdate transaction:")
		p.Factory.Json.Print(p.Txs.AmmUpdate[0])
	}

	// Print an example NftData transaction
	if len(p.Txs.NftData) > 0 {
		fmt.Println("Example NftData transaction:")
		p.Factory.Json.Print(p.Txs.NftData[0])
	}

	// Print an example Unknown transaction
	if len(p.Txs.TBD) > 0 {
		fmt.Println("Example Unknown transaction:")
		p.Factory.Json.Print(p.Txs.TBD[0])
	}
}

// Helper method to process a single transaction
func (p *Process) processTransaction(tx any) error {
	// Convert the transaction to JSON for unmarshaling
	txJSON, err := json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}

	var txTypeWrapper struct {
		TxType string `json:"txType"`
	}
	if err := json.Unmarshal(txJSON, &txTypeWrapper); err != nil {
		return fmt.Errorf("failed to unmarshal txType: %w", err)
	}

	// Process the transaction based on its txType
	switch txTypeWrapper.TxType {
	case "Deposit":
		return p.processDeposit(txJSON)
	case "Withdraw":
		return p.processWithdrawal(txJSON)
	case "SpotTrade":
		return p.processSwap(txJSON)
	case "Transfer":
		return p.processTransfer(txJSON)
	case "NftMint":
		return p.processMint(txJSON)
	case "AccountUpdate":
		return p.processAccountUpdate(txJSON)
	case "NftData":
		return p.processNftData(txJSON)
	case "AmmUpdate":
		return p.processAmmUpdate(txJSON)
	default:
		return p.processUnknown(txJSON)
	}
}

// Process Deposit/Withdraw transactions
func (p *Process) processDeposit(txJSON []byte) error {
	var dw DW
	if err := json.Unmarshal(txJSON, &dw); err != nil {
		return fmt.Errorf("failed to unmarshal DW transaction: %w", err)
	}
	p.Txs.Deposit = append(p.Txs.Deposit, dw)
	return nil
}

// Process Deposit/Withdraw transactions
func (p *Process) processWithdrawal(txJSON []byte) error {
	var dw DW
	if err := json.Unmarshal(txJSON, &dw); err != nil {
		return fmt.Errorf("failed to unmarshal DW transaction: %w", err)
	}
	p.Txs.Withdrawal = append(p.Txs.Withdrawal, dw)
	return nil
}

func (p *Process) processAccountUpdate(txJSON []byte) error {
	var update AccountUpdate
	if err := json.Unmarshal(txJSON, &update); err != nil {
		return fmt.Errorf("failed to unmarshal DW transaction: %w", err)
	}

	p.Txs.AccountUpdate = append(p.Txs.AccountUpdate, update)
	return nil
}

// Process SpotTrade transactions
func (p *Process) processSwap(txJSON []byte) error {
	var swap Swap
	if err := json.Unmarshal(txJSON, &swap); err != nil {
		return fmt.Errorf("failed to unmarshal Swap transaction: %w", err)
	}
	p.Txs.Swaps = append(p.Txs.Swaps, swap)
	return nil
}

// Process Transfer transactions
func (p *Process) processTransfer(txJSON []byte) error {
	var transfer Transfer
	if err := json.Unmarshal(txJSON, &transfer); err != nil {
		return fmt.Errorf("failed to unmarshal Transfer transaction: %w", err)
	}
	p.Txs.Transfers = append(p.Txs.Transfers, transfer)
	return nil
}

// Process NftMint transactions
func (p *Process) processMint(txJSON []byte) error {
	var mint Mint
	if err := json.Unmarshal(txJSON, &mint); err != nil {
		return fmt.Errorf("failed to unmarshal Mint transaction: %w", err)
	}
	p.Txs.Mints = append(p.Txs.Mints, mint)
	return nil
}

// Process NftData transactions
func (p *Process) processNftData(txJSON []byte) error {
	var nftData NftData
	if err := json.Unmarshal(txJSON, &nftData); err != nil {
		return fmt.Errorf("failed to unmarshal NftData transaction: %w", err)
	}
	p.Txs.NftData = append(p.Txs.NftData, nftData)
	return nil
}

// Process AmmUpdate transactions
func (p *Process) processAmmUpdate(txJSON []byte) error {
	var ammUpdate AmmUpdate
	if err := json.Unmarshal(txJSON, &ammUpdate); err != nil {
		return fmt.Errorf("failed to unmarshal AmmUpdate transaction: %w", err)
	}
	p.Txs.AmmUpdate = append(p.Txs.AmmUpdate, ammUpdate)
	return nil
}

func (p *Process) processUnknown(txJSON []byte) error {
	var unknown any
	if err := json.Unmarshal(txJSON, &unknown); err != nil {
		return fmt.Errorf("failed to unmarshal unknown transaction: %w", err)
	}
	p.Txs.TBD = append(p.Txs.TBD, unknown)
	return nil
}
