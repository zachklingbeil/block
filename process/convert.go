package process

import (
	"encoding/json"
	"fmt"
)

type Tx struct {
	Zero        any             `json:"zero,omitempty"`
	One         any             `json:"one,omitempty"`
	Value       string          `json:"value,omitempty"`
	Token       any             `json:"token,omitempty"`
	Fee         any             `json:"fee,omitempty"`
	FeeToken    int64           `json:"feeToken,omitempty"`
	OneValue    string          `json:"oneValue,omitempty"`
	OneToken    int64           `json:"oneToken,omitempty"`
	OneFee      any             `json:"oneFee,omitempty"`
	OneFeeToken int64           `json:"oneFeeToken,omitempty"`
	Type        string          `json:"type,omitempty"`
	Coordinates Coordinate      `json:"coordinates"`
	Raw         json.RawMessage `json:"raw,omitempty"`
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

// Convert a single DW (Deposit/Withdraw) transaction to Tx
func (p *Process) DepositToTx(dw DW) Tx {
	return Tx{
		Zero:        dw.Zero,
		One:         dw.One,
		Value:       dw.Value,
		Token:       dw.Token,
		Type:        "deposit",
		Coordinates: dw.Coordinates,
	}
}

func (p *Process) WithdrawToTx(dw DW) Tx {
	return Tx{
		Zero:        dw.Zero,
		One:         dw.One,
		Value:       dw.Value,
		Token:       dw.Token,
		Fee:         dw.Fee,
		FeeToken:    dw.FeeToken,
		Type:        "withdraw",
		Coordinates: dw.Coordinates,
	}
}

// Convert a single Swap transaction to Tx
func (p *Process) SwapToTx(swap Swap) Tx {
	return Tx{
		Zero:        swap.ZeroId,
		One:         swap.OneId,
		Value:       swap.ZeroValue,
		Token:       swap.ZeroToken,
		OneValue:    swap.OneValue,
		OneToken:    swap.OneToken,
		Fee:         swap.ZeroFee,
		OneFeeToken: swap.OneFee,
		Type:        "swap",
		Coordinates: swap.Coordinates,
	}
}

// Convert a single Transfer transaction to Tx
func (p *Process) TransferToTx(transfer Transfer) Tx {
	return Tx{
		Zero:        transfer.ZeroId,
		One:         transfer.One,
		Value:       transfer.Value,
		Token:       transfer.Token,
		Fee:         transfer.Fee,
		FeeToken:    transfer.FeeToken,
		Type:        "transfer",
		Coordinates: transfer.Coordinates,
	}
}

// Convert a single Mint transaction to Tx
func (p *Process) MintToTx(mint Mint) Tx {
	return Tx{
		Zero:        mint.Zero,
		Value:       mint.Quantity,
		Token:       mint.NftAddress,
		Fee:         mint.Fee,
		FeeToken:    mint.FeeToken,
		Type:        "mint",
		Coordinates: mint.Coordinates,
	}
}

// Convert a single AccountUpdate transaction to Tx
func (p *Process) AccountUpdateToTx(accountUpdate AccountUpdate) Tx {
	return Tx{
		Zero:        accountUpdate.ZeroId,
		Type:        "accountUpdate",
		Coordinates: accountUpdate.Coordinates,
	}
}

// Convert a single AmmUpdate transaction to Tx
func (p *Process) AmmUpdateToTx(ammUpdate AmmUpdate) Tx {
	return Tx{
		Zero:        ammUpdate.ZeroId,
		Type:        "ammUpdate",
		Coordinates: ammUpdate.Coordinates,
	}
}

// Convert a single NftData transaction to Tx
func (p *Process) NftDataToTx(nftData NftData) Tx {
	raw, err := json.Marshal(nftData)
	if err != nil {
		fmt.Printf("Error marshaling NftData to raw JSON: %v\n", err)
	}

	var rawMap map[string]any
	if err := json.Unmarshal(raw, &rawMap); err != nil {
		fmt.Printf("Error unmarshaling NftData to map: %v\n", err)
	}

	delete(rawMap, "accountId")
	delete(rawMap, "txType")
	delete(rawMap, "coordinates")

	filteredRaw, err := json.Marshal(rawMap)
	if err != nil {
		fmt.Printf("Error marshaling filtered raw map to JSON: %v\n", err)
	}
	return Tx{
		Zero:        nftData.ZeroId,
		Type:        "nftData",
		Coordinates: nftData.Coordinates,
		Raw:         filteredRaw,
	}
}
