package loopring

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/log"
)

type Tx struct {
	Zero        any             `json:"zero,omitempty"`
	One         any             `json:"one,omitempty"`
	Value       any             `json:"value,omitempty"`
	Token       any             `json:"token,omitempty"`
	Fee         any             `json:"fee,omitempty"`
	FeeToken    int64           `json:"feeToken,omitempty"`
	OneValue    any             `json:"oneValue,omitempty"`
	OneToken    int64           `json:"oneToken,omitempty"`
	OneFee      any             `json:"oneFee,omitempty"`
	OneFeeToken int64           `json:"oneFeeToken,omitempty"`
	Type        string          `json:"type,omitempty"`
	Index       uint16          `json:"index"`
	Raw         json.RawMessage `json:"raw,omitempty"`
}

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

func Int(value string) *big.Int {
	bigIntValue := new(big.Int)
	if _, ok := bigIntValue.SetString(value, 10); !ok {
		log.Error("Failed to convert string to big.Int: %s", value)
	}
	return bigIntValue
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

func (l *Loopring) SwapToTx(transaction any) Tx {
	var swap SpotTrade
	if err := mapToStruct(transaction, &swap); err != nil {
		fmt.Printf("Error unmarshaling transaction to SpotTrade: %v\n", err)
		return Tx{}
	}

	return Tx{
		Zero:        swap.Zero,
		One:         swap.One,
		Value:       Int(swap.ZeroValue),
		Token:       swap.ZeroToken,
		OneValue:    Int(swap.OneValue),
		OneToken:    swap.OneToken,
		Fee:         swap.ZeroFee,
		OneFeeToken: swap.OneFee,
		Type:        "swap",
		Index:       swap.Index,
	}
}

func (l *Loopring) TransferToTx(transaction any) Tx {
	var transfer Transfer
	if err := mapToStruct(transaction, &transfer); err != nil {
		fmt.Printf("Error unmarshaling transaction to Transfer: %v\n", err)
		return Tx{}
	}

	return Tx{
		Zero:     transfer.ZeroId,
		One:      transfer.One,
		Value:    Int(transfer.Value),
		Token:    transfer.Token,
		Fee:      Int(transfer.Fee),
		FeeToken: transfer.FeeToken,
		Type:     "transfer",
		Index:    transfer.Index,
	}
}

func (l *Loopring) DepositToTx(transaction any) Tx {
	var deposit Deposit
	if err := mapToStruct(transaction, &deposit); err != nil {
		fmt.Printf("Error unmarshaling transaction to Deposit: %v\n", err)
		return Tx{}
	}

	return Tx{
		Zero:  deposit.ZeroId,
		One:   deposit.One,
		Value: Int(deposit.Value),
		Token: deposit.Token,
		Type:  "deposit",
		Index: deposit.Index,
	}
}

func (l *Loopring) WithdrawToTx(transaction any) Tx {
	var withdrawal Withdrawal
	if err := mapToStruct(transaction, &withdrawal); err != nil {
		fmt.Printf("Error unmarshaling transaction to Withdrawal: %v\n", err)
		return Tx{}
	}

	return Tx{
		Zero:     withdrawal.ZeroId,
		One:      withdrawal.One,
		Value:    Int(withdrawal.Value),
		Token:    withdrawal.Token,
		Fee:      Int(withdrawal.Fee),
		FeeToken: withdrawal.FeeToken,
		Type:     "withdraw",
		Index:    withdrawal.Index,
	}
}

func (l *Loopring) AccountUpdateToTx(transaction any) Tx {
	var accountUpdate AccountUpdate
	if err := mapToStruct(transaction, &accountUpdate); err != nil {
		fmt.Printf("Error unmarshaling transaction to AccountUpdate: %v\n", err)
		return Tx{}
	}

	return Tx{
		Zero:  accountUpdate.ZeroId,
		Type:  "accountUpdate",
		Index: accountUpdate.Index,
	}
}

func (l *Loopring) AmmUpdateToTx(transaction any) Tx {
	var ammUpdate AmmUpdate
	if err := mapToStruct(transaction, &ammUpdate); err != nil {
		fmt.Printf("Error unmarshaling transaction to AmmUpdate: %v\n", err)
		return Tx{}
	}

	return Tx{
		Zero:  ammUpdate.ZeroId,
		Type:  "ammUpdate",
		Index: ammUpdate.Index,
	}
}

func (l *Loopring) MintToTx(transaction any) Tx {
	var mint Mint
	if err := mapToStruct(transaction, &mint); err != nil {
		fmt.Printf("Error unmarshaling transaction to Mint: %v\n", err)
		return Tx{}
	}

	return Tx{
		Zero:     mint.Zero,
		Value:    mint.Quantity,
		Token:    mint.NftAddress,
		Fee:      Int(mint.Fee),
		FeeToken: mint.FeeToken,
		Type:     "mint",
		Index:    mint.Index,
	}
}
func (l *Loopring) NftDataToTx(transaction any) Tx {
	var nftData NftData
	if err := mapToStruct(transaction, &nftData); err != nil {
		fmt.Printf("Error unmarshaling transaction to NftData: %v\n", err)
		return Tx{}
	}

	return Tx{
		Zero:  nftData.ZeroId,
		Type:  "nftData",
		Index: nftData.Index,
	}
}

// func (l *Loopring) NftDataToTx(nftData NftData) Tx {
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
// 	return Tx{
// 		Zero:  nftData.ZeroId,
// 		Type:  "nftData",
// 		Raw:   filteredRaw,
// 		Index: nftData.Index,
// 	}
// }
