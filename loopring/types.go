package loopring

import (
	"encoding/json"
	"math/big"
)

type Raw struct {
	Number       int64 `json:"blockId"`
	Timestamp    int64 `json:"createdAt"`
	Size         int64 `json:"blockSize"`
	Transactions []any `json:"transactions"`
}

type Block struct {
	Number int64      `json:"block"`
	Zero   Coordinate `json:"zero"`
	Ones   []Tx       `json:"one"`
}

type Coordinate struct {
	Year        uint8  `json:"year"`
	Month       uint8  `json:"month"`
	Day         uint8  `json:"day"`
	Hour        uint8  `json:"hour"`
	Minute      uint8  `json:"minute"`
	Second      uint8  `json:"second"`
	Millisecond uint16 `json:"millisecond"`
	Index       uint16 `json:"index"`
	Depth       uint16 `json:"depth,omitempty"`
}

type Tx struct {
	Zero     any             `json:"zero,omitempty"`
	One      any             `json:"one,omitempty"`
	Value    any             `json:"value,omitempty"`
	Token    any             `json:"token,omitempty"`
	ValueOut any             `json:"valueOut,omitempty"`
	TokenOut any             `json:"tokenOut,omitempty"`
	Fee      any             `json:"fee,omitempty"`
	FeeToken any             `json:"feeToken,omitempty"`
	Type     any             `json:"type,omitempty"`
	Index    any             `json:"index"`
	Nonce    any             `json:"nonce,omitempty"`
	Raw      json.RawMessage `json:"raw,omitempty"`
}

func mapToStruct(data any, target any) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, target)
}

type Swap struct {
	Zero     int64  `json:"orderA.accountID,omitempty"`
	One      int64  `json:"orderB.accountID,omitempty"`
	Value    string `json:"orderA.filledS,omitempty"`
	Token    int64  `json:"orderB.tokenB,omitempty"`
	ValueOut string `json:"orderB.filledS,omitempty"`
	TokenOut int64  `json:"orderA.tokenB,omitempty"`
	ZeroFee  int64  `json:"orderA.feeBips,omitempty"`
	OneFee   int64  `json:"orderB.feeBips,omitempty"`
	Type     string `json:"txType,omitempty"`
	Index    uint16 `json:"index"`
}

func (l *Loopring) SwapToTx(transaction any) Tx {
	var s Swap
	mapToStruct(transaction, &s)
	token := l.Value.GetAddress(s.Token)
	tokenOut := l.Value.GetAddress(s.TokenOut)
	tx := Tx{
		Zero:     l.Value.GetPeer(s.Zero),
		One:      l.Value.GetPeer(s.One),
		Value:    l.Value.Format(s.Value, token),
		ValueOut: l.Value.Format(s.ValueOut, tokenOut),
		Token:    token,
		TokenOut: tokenOut,
		Type:     "swap",
		Index:    s.Index,
	}

	var feeBips int64
	var valueInStr string

	if s.ZeroFee != 0 {
		feeBips = s.ZeroFee
		valueInStr = s.Value
	} else if s.OneFee != 0 {
		feeBips = s.OneFee
		valueInStr = s.Value
	}

	if feeBips != 0 && valueInStr != "" {
		valueIn := new(big.Int)
		valueIn.SetString(valueInStr, 10)
		fee := new(big.Int).Mul(valueIn, big.NewInt(feeBips))
		fee.Div(fee, big.NewInt(10000)) // Convert basis points to percentage
		tx.Fee = l.Value.Format(fee.String(), token)
		tx.FeeToken = token
	}
	return tx
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

func (l *Loopring) TransferToTx(transaction any) Tx {
	var t Transfer
	mapToStruct(transaction, &t)
	token := l.Value.GetAddress(t.Token)
	feeToken := l.Value.GetAddress(t.FeeToken)
	tx := Tx{
		Zero:  l.Value.GetPeer(t.ZeroId),
		One:   l.Value.GetPeer(t.OneId),
		Value: l.Value.Format(t.Value, token),
		Token: token,
		Index: t.Index,
		Type:  "transfer",
	}
	if t.Fee != "0" {
		tx.Fee = l.Value.Format(t.Fee, feeToken)
		tx.FeeToken = feeToken
	}
	return tx
}

// Depost,  Withdraw (fee)
type Deposit struct {
	Zero   string `json:"fromAddress"`
	ZeroId int64  `json:"accountId"`
	// One    string `json:"toAddress,omitempty"`
	Value string `json:"token.amount"`
	Token int64  `json:"token.tokenId"`
	Type  string `json:"txType,omitempty"`
	Index uint16 `json:"index"`
}

func (l *Loopring) DepositToTx(transaction any) Tx {
	var d Deposit
	mapToStruct(transaction, &d)
	token := l.Value.GetAddress(d.Token)
	return Tx{
		Zero:  l.Value.GetPeer(d.ZeroId),
		Value: l.Value.Format(d.Value, token),
		Token: token,
		Type:  "deposit",
		Index: d.Index,
	}
}

type Withdrawal struct {
	Zero   string `json:"fromAddress"`
	ZeroId int64  `json:"accountId"`
	// One      string `json:"toAddress,omitempty"`
	Value    string `json:"token.amount"`
	Token    int64  `json:"token.tokenId"`
	Fee      string `json:"fee.amount,omitempty"`
	FeeToken int64  `json:"fee.tokenId,omitempty"`
	Type     string `json:"txType,omitempty"`
	Index    uint16 `json:"index"`
}

func (l *Loopring) WithdrawToTx(transaction any) Tx {
	var w Withdrawal
	mapToStruct(transaction, &w)
	token := l.Value.GetAddress(w.Token)
	feeToken := l.Value.GetAddress(w.FeeToken)
	tx := Tx{
		Zero:  l.Value.GetPeer(w.ZeroId),
		Value: l.Value.Format(w.Value, token),
		Token: token,
		Type:  "withdraw",
		Index: w.Index,
	}
	if w.Fee != "0" {
		tx.Fee = l.Value.Format(w.Fee, feeToken)
		tx.FeeToken = feeToken
	}
	return tx
}

type AccountUpdate struct {
	ZeroId int64  `json:"accountId"`
	Type   string `json:"txType,omitempty"`
	Index  uint16 `json:"index"`
	Nonce  int64  `json:"nonce"`
}

func (l *Loopring) AccountUpdateToTx(transaction any) Tx {
	var a AccountUpdate
	mapToStruct(transaction, &a)
	return Tx{
		Zero:  l.Value.GetPeer(a.ZeroId),
		Type:  "accountUpdate",
		Index: a.Index,
		Nonce: a.Nonce,
	}
}

type AmmUpdate struct {
	Zero   string `json:"owner"`
	ZeroId int64  `json:"accountId"`
	Nonce  int64  `json:"nonce"`
	Type   string `json:"txType,omitempty"`
	Index  uint16 `json:"index"`
}

func (l *Loopring) AmmUpdateToTx(transaction any) Tx {
	var a AmmUpdate
	mapToStruct(transaction, &a)
	return Tx{
		Zero:  l.Value.GetPeer(a.ZeroId),
		Type:  "ammUpdate",
		Index: a.Index,
		Nonce: a.Nonce,
	}
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

func (l *Loopring) MintToTx(transaction any) Tx {
	var m Mint
	mapToStruct(transaction, &m)
	feeToken := l.Value.GetAddress(m.FeeToken)

	tx := Tx{
		Zero:  l.Value.GetPeer(m.ZeroId),
		Value: m.Quantity,
		Token: m.NftAddress,
		Type:  "mint",
		Index: m.Index,
	}
	if m.Fee != "0" {
		tx.Fee = l.Value.Format(m.Fee, feeToken)
		tx.FeeToken = feeToken
	}
	return tx
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

func (l *Loopring) NftDataToTx(transaction any) Tx {
	var n NftData
	mapToStruct(transaction, &n)
	return Tx{
		Zero:  l.Value.GetPeer(n.ZeroId),
		Type:  "nft",
		Index: n.Index,
		Raw:   []byte(n.NftData),
	}
}
