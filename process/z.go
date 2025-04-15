package process

import (
	"encoding/json"
	"fmt"
)

type Tx struct {
	Zero any `json:"zero,omitempty"`
	One  any `json:"one,omitempty"`

	Value    string `json:"value,omitempty"`
	Token    any    `json:"token,omitempty"`
	Fee      any    `json:"fee,omitempty"`
	FeeToken int64  `json:"feeToken,omitempty"`

	OneValue    string          `json:"oneValue,omitempty"`
	OneToken    int64           `json:"oneToken,omitempty"`
	OneFee      any             `json:"oneFee,omitempty"`
	OneFeeToken int64           `json:"oneFeeToken,omitempty"`
	Type        string          `json:"type,omitempty"`
	Coordinates Coordinate      `json:"coordinates"`
	Raw         json.RawMessage `json:"raw,omitempty"`
}

// Depost,  Withdraw (fee)
type DW struct {
	Zero        string     `json:"fromAddress"`
	ZeroId      int64      `json:"accountId"`
	One         string     `json:"toAddress"`
	Value       string     `json:"token.amount"`
	Token       int64      `json:"token.tokenId"`
	Fee         string     `json:"fee.amount,omitempty"`
	FeeToken    int64      `json:"fee.tokenId,omitempty"`
	Type        string     `json:"txType,omitempty"`
	Coordinates Coordinate `json:"coordinates"`
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

type Swap struct {
	ZeroId      int64      `json:"orderA.accountID"`
	ZeroValue   string     `json:"orderA.filledS"`
	ZeroToken   int64      `json:"orderB.tokenB"`
	OneId       int64      `json:"orderB.accountID"`
	OneValue    string     `json:"orderB.filledS"`
	OneToken    int64      `json:"orderA.tokenB"`
	ZeroFee     int64      `json:"orderA.feeBips"`
	OneFee      int64      `json:"fee.orderB.feeBips"`
	Type        string     `json:"txType,omitempty"`
	Coordinates Coordinate `json:"coordinates"`
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

type Transfer struct {
	ZeroId      int64      `json:"accountId"`
	OneId       int64      `json:"toAccountId"`
	One         string     `json:"toAccountAddress"`
	Value       string     `json:"token.amount"`
	Token       int64      `json:"token.tokenId"`
	Fee         string     `json:"fee.amount,omitempty"`
	FeeToken    int64      `json:"fee.tokenId,omitempty"`
	Type        string     `json:"txType,omitempty"`
	Coordinates Coordinate `json:"coordinates"`
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

type Mint struct {
	ZeroId      int64      `json:"minterAccountId"`
	Zero        string     `json:"toAccountAddress"`
	Nft         any        `json:"toToken.tokenId"`
	NftId       string     `json:"nftToken.nftId"`
	NftData     string     `json:"nftToken.nftData"`
	NftAddress  string     `json:"nftToken.tokenAddress"`
	Quantity    string     `json:"nftToken.amount"`
	Fee         string     `json:"fee.amount,omitempty"`
	FeeToken    int64      `json:"fee.tokenId,omitempty"`
	Type        string     `json:"txType,omitempty"`
	Coordinates Coordinate `json:"coordinates"`
}

// Convert a single AccountUpdate transaction to Tx
func (p *Process) AccountUpdateToTx(accountUpdate AccountUpdate) Tx {
	return Tx{
		Zero:        accountUpdate.ZeroId,
		Type:        "accountUpdate",
		Coordinates: accountUpdate.Coordinates,
	}
}

type AccountUpdate struct {
	ZeroId      int64      `json:"accountId"`
	Type        string     `json:"txType,omitempty"`
	Coordinates Coordinate `json:"coordinates"`
}

// Convert a single AmmUpdate transaction to Tx
func (p *Process) AmmUpdateToTx(ammUpdate AmmUpdate) Tx {
	return Tx{
		Zero:        ammUpdate.ZeroId,
		Type:        "ammUpdate",
		Coordinates: ammUpdate.Coordinates,
	}
}

type AmmUpdate struct {
	Zero        string     `json:"owner"`
	ZeroId      int64      `json:"accountId"`
	Nonce       int64      `json:"nonce"`
	Type        string     `json:"txType,omitempty"`
	Coordinates Coordinate `json:"coordinates"`
}

// Convert a single NftData transaction to Tx
func (p *Process) NftDataToTx(nftData NftData) Tx {
	// Marshal the original NftData struct to JSON
	raw, err := json.Marshal(nftData)
	if err != nil {
		fmt.Printf("Error marshaling NftData to raw JSON: %v\n", err)
	}

	// Unmarshal the JSON into a map to extract additional fields
	var rawMap map[string]interface{}
	if err := json.Unmarshal(raw, &rawMap); err != nil {
		fmt.Printf("Error unmarshaling NftData to map: %v\n", err)
	}

	// Remove fields that are explicitly mapped in the Tx struct
	delete(rawMap, "accountId")
	delete(rawMap, "txType")
	delete(rawMap, "coordinates")

	// Marshal the remaining fields back into JSON for dynamic storage
	filteredRaw, err := json.Marshal(rawMap)
	if err != nil {
		fmt.Printf("Error marshaling filtered raw map to JSON: %v\n", err)
	}

	return Tx{
		Zero:        nftData.ZeroId,
		Type:        "nftData",
		Coordinates: nftData.Coordinates,
		Raw:         filteredRaw, // Store only the additional fields
	}
}

type NftData struct {
	ZeroId      int64      `json:"accountId"`
	One         string     `json:"minter"`
	NftId       string     `json:"nftToken.nftId"`
	NftData     string     `json:"nftToken.nftData,omitempty"`
	NftAddress  string     `json:"nftToken.tokenAddress"`
	Type        string     `json:"txType,omitempty"`
	Coordinates Coordinate `json:"coordinates"`
}
