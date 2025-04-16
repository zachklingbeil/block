package process

import (
	"encoding/json"
	"fmt"
)

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

// Helper method to process a single transaction
func (p *Process) processTransaction(txJSON []byte) error {
	var txTypeWrapper struct {
		TxType string `json:"txType"`
	}
	if err := json.Unmarshal(txJSON, &txTypeWrapper); err != nil {
		return fmt.Errorf("failed to unmarshal txType: %w", err)
	}

	processors := map[string]func([]byte) error{
		"Deposit":       p.processDeposit,
		"Withdraw":      p.processWithdrawal,
		"SpotTrade":     p.processSwap,
		"Transfer":      p.processTransfer,
		"NftMint":       p.processMint,
		"AccountUpdate": p.processAccountUpdate,
		"NftData":       p.processNftData,
		"AmmUpdate":     p.processAmmUpdate,
		"Unknown":       p.processUnknown,
	}

	if processor, exists := processors[txTypeWrapper.TxType]; exists {
		return processor(txJSON)
	}
	return nil
}

// Process Deposit transactions
func (p *Process) processDeposit(txJSON []byte) error {
	var dw DW
	if err := json.Unmarshal(txJSON, &dw); err != nil {
		return fmt.Errorf("failed to unmarshal DW transaction: %w", err)
	}
	p.Types.Deposit = append(p.Types.Deposit, dw)
	p.Counts["Deposits"]++
	return nil
}

// Process Withdrawal transactions
func (p *Process) processWithdrawal(txJSON []byte) error {
	var dw DW
	if err := json.Unmarshal(txJSON, &dw); err != nil {
		return fmt.Errorf("failed to unmarshal DW transaction: %w", err)
	}
	p.Types.Withdrawal = append(p.Types.Withdrawal, dw)
	p.Counts["Withdrawals"]++
	return nil
}

// Process AccountUpdate transactions
func (p *Process) processAccountUpdate(txJSON []byte) error {
	var update AccountUpdate
	if err := json.Unmarshal(txJSON, &update); err != nil {
		return fmt.Errorf("failed to unmarshal AccountUpdate transaction: %w", err)
	}
	p.Types.AccountUpdate = append(p.Types.AccountUpdate, update)
	p.Counts["AccountUpdates"]++
	return nil
}

// Process SpotTrade transactions
func (p *Process) processSwap(txJSON []byte) error {
	var swap Swap
	if err := json.Unmarshal(txJSON, &swap); err != nil {
		return fmt.Errorf("failed to unmarshal Swap transaction: %w", err)
	}
	p.Types.Swaps = append(p.Types.Swaps, swap)
	p.Counts["Swaps"]++
	return nil
}

// Process Transfer transactions
func (p *Process) processTransfer(txJSON []byte) error {
	var transfer Transfer
	if err := json.Unmarshal(txJSON, &transfer); err != nil {
		return fmt.Errorf("failed to unmarshal Transfer transaction: %w", err)
	}
	p.Types.Transfers = append(p.Types.Transfers, transfer)
	p.Counts["Transfers"]++
	return nil
}

// Process NftMint transactions
func (p *Process) processMint(txJSON []byte) error {
	var mint Mint
	if err := json.Unmarshal(txJSON, &mint); err != nil {
		return fmt.Errorf("failed to unmarshal Mint transaction: %w", err)
	}
	p.Types.Mints = append(p.Types.Mints, mint)
	p.Counts["Mints"]++
	return nil
}

// Process NftData transactions
func (p *Process) processNftData(txJSON []byte) error {
	var nftData NftData
	if err := json.Unmarshal(txJSON, &nftData); err != nil {
		return fmt.Errorf("failed to unmarshal NftData transaction: %w", err)
	}
	p.Types.NftData = append(p.Types.NftData, nftData)
	p.Counts["NftData"]++
	return nil
}

// Process AmmUpdate transactions
func (p *Process) processAmmUpdate(txJSON []byte) error {
	var ammUpdate AmmUpdate
	if err := json.Unmarshal(txJSON, &ammUpdate); err != nil {
		return fmt.Errorf("failed to unmarshal AmmUpdate transaction: %w", err)
	}
	p.Types.AmmUpdate = append(p.Types.AmmUpdate, ammUpdate)
	p.Counts["AmmUpdates"]++
	return nil
}

// Process Unknown transactions
func (p *Process) processUnknown(txJSON []byte) error {
	var unknown any
	if err := json.Unmarshal(txJSON, &unknown); err != nil {
		return fmt.Errorf("failed to unmarshal unknown transaction: %w", err)
	}
	p.Types.TBD = append(p.Types.TBD, unknown)
	p.Counts["Unknown"]++
	return nil
}
