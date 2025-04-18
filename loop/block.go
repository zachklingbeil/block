package loop

import (
	"encoding/json"
	"fmt"
)

type Block struct {
	Number       int64 `json:"blockId"`
	Size         int64 `json:"blockSize"`
	Timestamp    int64 `json:"createdAt"`
	Transactions []any `json:"transactions"`
}

type Coordinate struct {
	Block       int64
	Year        int64 `json:"year"`
	Month       int64 `json:"month"`
	Day         int64 `json:"day"`
	Hour        int64 `json:"hour"`
	Minute      int64 `json:"minute"`
	Second      int64 `json:"second"`
	Millisecond int64 `json:"millisecond"`
	Index       int64 `json:"index"`
}

type Transaction struct {
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
	Raw         json.RawMessage `json:"raw,omitempty"`
}

type Type struct {
	Deposits       []DW            `json:"Deposit,omitempty"`
	Withdrawals    []DW            `json:"Withdraw,omitempty"`
	Swaps          []Swap          `json:"SpotTrade,omitempty"`
	Transfers      []Transfer      `json:"Transfer,omitempty"`
	Mints          []Mint          `json:"NftMint,omitempty"`
	AccountUpdates []AccountUpdate `json:"AccountUpdate,omitempty"`
	AmmUpdates     []AmmUpdate     `json:"AmmUpdate,omitempty"`
	NftData        []NftData       `json:"NftData,omitempty"`
	TBD            []any           `json:"tbd,omitempty"`
	*json.RawMessage
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

func (l *Loopring) SwapToTx(swap Swap) Transaction {
	return Transaction{
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

func (l *Loopring) TransferToTx(transfer Transfer) Transaction {
	return Transaction{
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

func (l *Loopring) DepositToTx(dw DW) Transaction {
	return Transaction{
		Zero:        dw.Zero,
		One:         dw.One,
		Value:       dw.Value,
		Token:       dw.Token,
		Type:        "deposit",
		Coordinates: dw.Coordinates,
	}
}

func (l *Loopring) WithdrawToTx(dw DW) Transaction {
	return Transaction{
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

type AccountUpdate struct {
	ZeroId      int64      `json:"accountId"`
	Type        string     `json:"txType,omitempty"`
	Coordinates Coordinate `json:"coordinates"`
}

func (l *Loopring) AccountUpdateToTx(accountUpdate AccountUpdate) Transaction {
	return Transaction{
		Zero:        accountUpdate.ZeroId,
		Type:        "accountUpdate",
		Coordinates: accountUpdate.Coordinates,
	}
}

type AmmUpdate struct {
	Zero        string     `json:"owner"`
	ZeroId      int64      `json:"accountId"`
	Nonce       int64      `json:"nonce"`
	Type        string     `json:"txType,omitempty"`
	Coordinates Coordinate `json:"coordinates"`
}

func (l *Loopring) AmmUpdateToTx(ammUpdate AmmUpdate) Transaction {
	return Transaction{
		Zero:        ammUpdate.ZeroId,
		Type:        "ammUpdate",
		Coordinates: ammUpdate.Coordinates,
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

func (l *Loopring) MintToTx(mint Mint) Transaction {
	return Transaction{
		Zero:        mint.Zero,
		Value:       mint.Quantity,
		Token:       mint.NftAddress,
		Fee:         mint.Fee,
		FeeToken:    mint.FeeToken,
		Type:        "mint",
		Coordinates: mint.Coordinates,
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

func (l *Loopring) NftDataToTx(nftData NftData) Transaction {
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
	return Transaction{
		Zero:        nftData.ZeroId,
		Type:        "nftData",
		Coordinates: nftData.Coordinates,
		Raw:         filteredRaw,
	}
}

func (l *Loopring) Unload(transactions []any) {
	for _, tx := range transactions {
		txMap, ok := tx.(map[string]any)
		if !ok {
			continue
		}
		txType, ok := txMap["txType"].(string)
		if !ok {
			continue
		}

		txBytes, err := json.Marshal(txMap)
		if err != nil {
			continue
		}

		switch txType {
		case "Deposit":
			var dw DW
			if err := json.Unmarshal(txBytes, &dw); err == nil {
				l.Types.Deposits = append(l.Types.Deposits, dw)
			}
		case "Withdraw":
			var dw DW
			if err := json.Unmarshal(txBytes, &dw); err == nil {
				l.Types.Withdrawals = append(l.Types.Withdrawals, dw)
			}
		case "SpotTrade":
			var swap Swap
			if err := json.Unmarshal(txBytes, &swap); err == nil {
				l.Types.Swaps = append(l.Types.Swaps, swap)
			}
		case "Transfer":
			var transfer Transfer
			if err := json.Unmarshal(txBytes, &transfer); err == nil {
				l.Types.Transfers = append(l.Types.Transfers, transfer)
			}
		case "NftMint":
			var mint Mint
			if err := json.Unmarshal(txBytes, &mint); err == nil {
				l.Types.Mints = append(l.Types.Mints, mint)
			}
		case "AccountUpdate":
			var au AccountUpdate
			if err := json.Unmarshal(txBytes, &au); err == nil {
				l.Types.AccountUpdates = append(l.Types.AccountUpdates, au)
			}
		case "AmmUpdate":
			var amm AmmUpdate
			if err := json.Unmarshal(txBytes, &amm); err == nil {
				l.Types.AmmUpdates = append(l.Types.AmmUpdates, amm)
			}
		case "NftData":
			var nft NftData
			if err := json.Unmarshal(txBytes, &nft); err == nil {
				l.Types.NftData = append(l.Types.NftData, nft)
			}
		default:
			l.Types.TBD = append(l.Types.TBD, tx)
		}
	}
}

func (l *Loopring) Simplify() {
	for _, deposit := range l.Types.Deposits {
		l.Transactions = append(l.Transactions, l.DepositToTx(deposit))
	}

	for _, withdrawal := range l.Types.Withdrawals {
		l.Transactions = append(l.Transactions, l.WithdrawToTx(withdrawal))
	}

	for _, swap := range l.Types.Swaps {
		l.Transactions = append(l.Transactions, l.SwapToTx(swap))
	}

	for _, transfer := range l.Types.Transfers {
		l.Transactions = append(l.Transactions, l.TransferToTx(transfer))
	}

	for _, mint := range l.Types.Mints {
		l.Transactions = append(l.Transactions, l.MintToTx(mint))
	}

	for _, accountUpdate := range l.Types.AccountUpdates {
		l.Transactions = append(l.Transactions, l.AccountUpdateToTx(accountUpdate))
	}

	for _, ammUpdate := range l.Types.AmmUpdates {
		l.Transactions = append(l.Transactions, l.AmmUpdateToTx(ammUpdate))
	}

	for _, nftData := range l.Types.NftData {
		l.Transactions = append(l.Transactions, l.NftDataToTx(nftData))
	}
}
