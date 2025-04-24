package loopring

import (
	"encoding/json"
	"fmt"

	"github.com/zachklingbeil/block/circuit"
)

type Transfer struct {
	ZeroId   int64  `json:"accountId"`
	OneId    int64  `json:"toAccountId"`
	One      string `json:"toAccountAddress"`
	Value    string `json:"token.amount"`
	Token    int64  `json:"token.tokenId"`
	Fee      string `json:"fee.amount,omitempty"`
	FeeToken int64  `json:"fee.tokenId,omitempty"`
	Type     string `json:"txType,omitempty"`
	Index    uint16 `json:"index"`
}

// Depost,  Withdraw (fee)
type Deposit struct {
	Zero   string `json:"fromAddress"`
	ZeroId int64  `json:"accountId"`
	One    string `json:"toAddress"`
	Value  string `json:"token.amount"`
	Token  int64  `json:"token.tokenId"`
	Type   string `json:"txType,omitempty"`
	Index  uint16 `json:"index"`
}

type Withdrawal struct {
	Zero     string `json:"fromAddress"`
	ZeroId   int64  `json:"accountId"`
	One      string `json:"toAddress"`
	Value    string `json:"token.amount"`
	Token    int64  `json:"token.tokenId"`
	Fee      string `json:"fee.amount,omitempty"`
	FeeToken int64  `json:"fee.tokenId,omitempty"`
	Type     string `json:"txType,omitempty"`
	Index    uint16 `json:"index"`
}

type AccountUpdate struct {
	ZeroId int64  `json:"accountId"`
	Type   string `json:"txType,omitempty"`
	Index  uint16 `json:"index"`
}

type AmmUpdate struct {
	Zero   string `json:"owner"`
	ZeroId int64  `json:"accountId"`
	Nonce  int64  `json:"nonce"`
	Type   string `json:"txType,omitempty"`
	Index  uint16 `json:"index"`
}

type Mint struct {
	ZeroId     int64  `json:"minterAccountId"`
	Zero       string `json:"toAccountAddress"`
	Nft        any    `json:"toToken.tokenId"`
	NftId      string `json:"nftToken.nftId"`
	NftData    string `json:"nftToken.nftData"`
	NftAddress string `json:"nftToken.tokenAddress"`
	Quantity   string `json:"nftToken.amount"`
	Fee        string `json:"fee.amount,omitempty"`
	FeeToken   int64  `json:"fee.tokenId,omitempty"`
	Type       string `json:"txType,omitempty"`
	Index      uint16 `json:"index"`
}

type NftData struct {
	ZeroId     int64  `json:"accountId"`
	One        string `json:"minter"`
	NftId      string `json:"nftToken.nftId"`
	NftData    string `json:"nftToken.nftData,omitempty"`
	NftAddress string `json:"nftToken.tokenAddress"`
	Type       string `json:"txType,omitempty"`
	Index      uint16 `json:"index"`
}

func mapToStruct(data any, target any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	if err := json.Unmarshal(bytes, target); err != nil {
		return fmt.Errorf("failed to unmarshal data into target struct: %w", err)
	}
	return nil
}

type Tx struct {
	Zero     any             `json:"zero,omitempty"`
	One      any             `json:"one,omitempty"`
	Value    any             `json:"value,omitempty"`
	Token    any             `json:"token,omitempty"`
	Fee      any             `json:"fee,omitempty"`
	FeeToken int64           `json:"feeToken,omitempty"`
	Type     string          `json:"type,omitempty"`
	Index    uint16          `json:"index"`
	Raw      json.RawMessage `json:"raw,omitempty"`
}

type SpotTrade struct {
	Zero      int64  `json:"orderA.accountID"`
	ZeroValue string `json:"orderA.filledS"`
	ZeroToken int64  `json:"orderB.tokenB"`
	One       int64  `json:"orderB.accountID"`
	OneValue  string `json:"orderB.filledS"`
	OneToken  int64  `json:"orderA.tokenB"`
	ZeroFee   int64  `json:"orderA.feeBips"`
	OneFee    int64  `json:"orderB.feeBips"`
	Type      string `json:"txType,omitempty"`
	Index     uint16 `json:"index"`
}

func (l *Loopring) SwapToTx(transaction any) []circuit.Tx {
	var s SpotTrade
	if err := mapToStruct(transaction, &s); err != nil {
		fmt.Printf("Error unmarshaling transaction to SpotTrade: %v\n", err)
		return nil
	}

	zero := circuit.Tx{
		Zero:  s.Zero,
		One:   s.One,
		Value: s.ZeroValue,
		Token: s.ZeroToken,
		Fee:   s.ZeroFee,
		// Type:  "swap",
		Index: s.Index,
	}

	one := circuit.Tx{
		Zero:  s.One,
		One:   s.Zero,
		Value: s.OneValue,
		Token: s.OneToken,
		Fee:   s.OneFee,
		// Type:  "swap",
		Index: s.Index,
	}
	return []circuit.Tx{zero, one}
}

func (l *Loopring) TransferToTx(transaction any) circuit.Tx {
	var t Transfer
	if err := mapToStruct(transaction, &t); err != nil {
		fmt.Printf("Error unmarshaling transaction to Transfer: %v\n", err)
		return circuit.Tx{}
	}

	return circuit.Tx{
		Zero:     t.ZeroId,
		One:      t.One,
		Value:    t.Value,
		Token:    t.Token,
		Fee:      t.Fee,
		FeeToken: t.FeeToken,
		// Type:     "transfer",
		Index: t.Index,
	}
}

func (l *Loopring) DepositToTx(transaction any) circuit.Tx {
	var d Deposit
	if err := mapToStruct(transaction, &d); err != nil {
		fmt.Printf("Error unmarshaling transaction to Deposit: %v\n", err)
		return circuit.Tx{}
	}

	return circuit.Tx{
		Zero:  d.ZeroId,
		One:   d.One,
		Value: d.Value,
		Token: d.Token,
		// Type:  "deposit",
		Index: d.Index,
	}
}

func (l *Loopring) WithdrawToTx(transaction any) circuit.Tx {
	var w Withdrawal
	if err := mapToStruct(transaction, &w); err != nil {
		fmt.Printf("Error unmarshaling transaction to Withdrawal: %v\n", err)
		return circuit.Tx{}
	}

	return circuit.Tx{
		Zero:     w.ZeroId,
		One:      w.One,
		Value:    w.Value,
		Token:    w.Token,
		Fee:      w.Fee,
		FeeToken: w.FeeToken,
		// Type:     "withdraw",
		Index: w.Index,
	}
}

func (l *Loopring) AccountUpdateToTx(transaction any) circuit.Tx {
	var a AccountUpdate
	if err := mapToStruct(transaction, &a); err != nil {
		fmt.Printf("Error unmarshaling transaction to AccountUpdate: %v\n", err)
		return circuit.Tx{}
	}
	return circuit.Tx{
		Zero: a.ZeroId,
		// Type:  "accountUpdate",
		Index: a.Index,
	}
}

func (l *Loopring) AmmUpdateToTx(transaction any) circuit.Tx {
	var a AmmUpdate
	if err := mapToStruct(transaction, &a); err != nil {
		fmt.Printf("Error unmarshaling transaction to AmmUpdate: %v\n", err)
		return circuit.Tx{}
	}
	return circuit.Tx{
		Zero: a.ZeroId,
		// Type:  "ammUpdate",
		Index: a.Index,
	}
}

func (l *Loopring) MintToTx(transaction any) circuit.Tx {
	var m Mint
	if err := mapToStruct(transaction, &m); err != nil {
		fmt.Printf("Error unmarshaling transaction to Mint: %v\n", err)
		return circuit.Tx{}
	}
	return circuit.Tx{
		Zero:     m.Zero,
		Value:    m.Quantity,
		Token:    m.NftAddress,
		Fee:      m.Fee,
		FeeToken: m.FeeToken,
		// Type:     "mint",
		Index: m.Index,
	}
}
func (l *Loopring) NftDataToTx(transaction any) circuit.Tx {
	var n NftData
	if err := mapToStruct(transaction, &n); err != nil {
		fmt.Printf("Error unmarshaling transaction to NftData: %v\n", err)
		return circuit.Tx{}
	}
	return circuit.Tx{
		Zero: n.ZeroId,
		// Type:  "nft",
		Index: n.Index,
	}
}

// func (l *Loopring) NftDataToTx(nftData NftData) circuit.Tx {
// 	raw, err := json.Marshal(nftData)
// 	if err != nil {
// 		fmt.Printf("Error marshaling NftData to raw JSON: %v\n", err)
// 	}

// 	var rawMap map[string]any
// 	if err := json.Unmarshal(raw, &rawMap); err != nil {
// 		fmt.Printf("Error unmarshaling NftData to map: %v\n", err)
// 	}

// 	delete(rawMap, "accountId")
// 	delete(rawMap, "txType")
// 	delete(rawMap, "coordinates")

// 	filteredRaw, err := json.Marshal(rawMap)
// 	if err != nil {
// 		fmt.Printf("Error marshaling filtered raw map to JSON: %v\n", err)
// 	}
// 	return circuit.Tx{
// 		Zero:  nftData.ZeroId,
// 		Type:  "nftData",
// 		Raw:   filteredRaw,
// 		Index: nftData.Index,
// 	}
// }
