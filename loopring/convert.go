package loopring

import (
	"encoding/json"
	"fmt"

	"github.com/zachklingbeil/factory/fx"
)

type Swap struct {
	ZeroId      int64   `json:"orderA.accountID"`
	ZeroValue   string  `json:"orderA.filledS"`
	ZeroToken   int64   `json:"orderB.tokenB"`
	OneId       int64   `json:"orderB.accountID"`
	OneValue    string  `json:"orderB.filledS"`
	OneToken    int64   `json:"orderA.tokenB"`
	ZeroFee     int64   `json:"orderA.feeBips"`
	OneFee      int64   `json:"fee.orderB.feeBips"`
	Type        string  `json:"txType,omitempty"`
	Coordinates fx.Zero `json:"coordinates"`
	Index       uint16  `json:"index"`
}

type Transfer struct {
	ZeroId      int64   `json:"accountId"`
	OneId       int64   `json:"toAccountId"`
	One         string  `json:"toAccountAddress"`
	Value       string  `json:"token.amount"`
	Token       int64   `json:"token.tokenId"`
	Fee         string  `json:"fee.amount,omitempty"`
	FeeToken    int64   `json:"fee.tokenId,omitempty"`
	Type        string  `json:"txType,omitempty"`
	Coordinates fx.Zero `json:"coordinates"`
	Index       uint16  `json:"index"`
}

// Depost,  Withdraw (fee)
type DW struct {
	Zero        string  `json:"fromAddress"`
	ZeroId      int64   `json:"accountId"`
	One         string  `json:"toAddress"`
	Value       string  `json:"token.amount"`
	Token       int64   `json:"token.tokenId"`
	Fee         string  `json:"fee.amount,omitempty"`
	FeeToken    int64   `json:"fee.tokenId,omitempty"`
	Type        string  `json:"txType,omitempty"`
	Coordinates fx.Zero `json:"coordinates"`
	Index       uint16  `json:"index"`
}

type AccountUpdate struct {
	ZeroId      int64   `json:"accountId"`
	Type        string  `json:"txType,omitempty"`
	Coordinates fx.Zero `json:"coordinates"`
	Index       uint16  `json:"index"`
}

type AmmUpdate struct {
	Zero        string  `json:"owner"`
	ZeroId      int64   `json:"accountId"`
	Nonce       int64   `json:"nonce"`
	Type        string  `json:"txType,omitempty"`
	Coordinates fx.Zero `json:"coordinates"`
	Index       uint16  `json:"index"`
}

type Mint struct {
	ZeroId      int64   `json:"minterAccountId"`
	Zero        string  `json:"toAccountAddress"`
	Nft         any     `json:"toToken.tokenId"`
	NftId       string  `json:"nftToken.nftId"`
	NftData     string  `json:"nftToken.nftData"`
	NftAddress  string  `json:"nftToken.tokenAddress"`
	Quantity    string  `json:"nftToken.amount"`
	Fee         string  `json:"fee.amount,omitempty"`
	FeeToken    int64   `json:"fee.tokenId,omitempty"`
	Type        string  `json:"txType,omitempty"`
	Coordinates fx.Zero `json:"coordinates"`
	Index       uint16  `json:"index"`
}

type NftData struct {
	ZeroId      int64   `json:"accountId"`
	One         string  `json:"minter"`
	NftId       string  `json:"nftToken.nftId"`
	NftData     string  `json:"nftToken.nftData,omitempty"`
	NftAddress  string  `json:"nftToken.tokenAddress"`
	Type        string  `json:"txType,omitempty"`
	Coordinates fx.Zero `json:"coordinates"`
	Index       uint16  `json:"index"`
}

func (l *Loopring) SwapToTx(swap Swap) Tx {
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
		Index:       swap.Index,
	}
}

func (l *Loopring) TransferToTx(transfer Transfer) Tx {
	return Tx{
		Zero:        transfer.ZeroId,
		One:         transfer.One,
		Value:       transfer.Value,
		Token:       transfer.Token,
		Fee:         transfer.Fee,
		FeeToken:    transfer.FeeToken,
		Type:        "transfer",
		Coordinates: transfer.Coordinates,
		Index:       transfer.Index,
	}
}

func (l *Loopring) DepositToTx(dw DW) Tx {
	return Tx{
		Zero:        dw.Zero,
		One:         dw.One,
		Value:       dw.Value,
		Token:       dw.Token,
		Type:        "deposit",
		Coordinates: dw.Coordinates,
		Index:       dw.Index,
	}
}

func (l *Loopring) WithdrawToTx(dw DW) Tx {
	return Tx{
		Zero:        dw.Zero,
		One:         dw.One,
		Value:       dw.Value,
		Token:       dw.Token,
		Fee:         dw.Fee,
		FeeToken:    dw.FeeToken,
		Type:        "withdraw",
		Coordinates: dw.Coordinates,
		Index:       dw.Index,
	}
}

func (l *Loopring) AccountUpdateToTx(accountUpdate AccountUpdate) Tx {
	return Tx{
		Zero:        accountUpdate.ZeroId,
		Type:        "accountUpdate",
		Coordinates: accountUpdate.Coordinates,
		Index:       accountUpdate.Index,
	}
}

func (l *Loopring) AmmUpdateToTx(ammUpdate AmmUpdate) Tx {
	return Tx{
		Zero:        ammUpdate.ZeroId,
		Type:        "ammUpdate",
		Coordinates: ammUpdate.Coordinates,
		Index:       ammUpdate.Index,
	}
}

func (l *Loopring) MintToTx(mint Mint) Tx {
	return Tx{
		Zero:        mint.Zero,
		Value:       mint.Quantity,
		Token:       mint.NftAddress,
		Fee:         mint.Fee,
		FeeToken:    mint.FeeToken,
		Type:        "mint",
		Coordinates: mint.Coordinates,
		Index:       mint.Index,
	}
}

func (l *Loopring) NftDataToTx(nftData NftData) Tx {
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
		Index:       nftData.Index,
	}
}
