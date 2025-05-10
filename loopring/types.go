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
type Deposit struct {
	Zero   string `json:"fromAddress"`
	ZeroId int64  `json:"accountId"`
	Value  string `json:"token.amount"`
	Token  int64  `json:"token.tokenId"`
	Type   string `json:"txType,omitempty"`
	Index  uint16 `json:"index"`
}

type Withdrawal struct {
	Zero     string `json:"fromAddress"`
	ZeroId   int64  `json:"accountId"`
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
	Nonce  int64  `json:"nonce"`
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

func (l *Loopring) SwapToTx(transaction any) Tx {
	var s Swap
	mapToStruct(transaction, &s)
	tx := Tx{
		Zero:     l.One.LoopringId(s.Zero).Address,
		One:      l.One.LoopringId(s.One).Address,
		Value:    l.One.Format(s.Value, l.One.TokenId(s.Token).Decimals),
		ValueOut: l.One.Format(s.ValueOut, l.One.TokenId(s.TokenOut).Decimals),
		Token:    l.One.TokenId(s.Token).Address,
		TokenOut: l.One.TokenId(s.TokenOut).Address,
		Type:     "swap",
		Index:    s.Index,
		FeeToken: l.One.TokenId(s.Token).Address,
	}

	switch {
	case s.ZeroFee != 0:
		fee := calcFee(s.Value, s.ZeroFee)
		tx.Fee = l.One.Format(fee, l.One.TokenId(s.Token).Decimals)
	case s.OneFee != 0:
		fee := calcFee(s.Value, s.OneFee)
		tx.Fee = l.One.Format(fee, l.One.TokenId(s.Token).Decimals)
	}
	return tx
}

func calcFee(valueStr string, feeBips int64) string {
	valueIn := new(big.Int)
	valueIn.SetString(valueStr, 10)
	fee := new(big.Int).Mul(valueIn, big.NewInt(feeBips))
	fee.Div(fee, big.NewInt(10000)) // Convert basis points to percentage
	return fee.String()
}

func (l *Loopring) TransferToTx(transaction any) Tx {
	var t Transfer
	mapToStruct(transaction, &t)

	tx := Tx{
		Zero:     l.One.LoopringId(t.ZeroId).Address,
		One:      l.One.LoopringId(t.OneId).Address,
		Value:    l.One.Format(t.Value, l.One.TokenId(t.Token).Decimals),
		Token:    l.One.TokenId(t.Token).Address,
		Index:    t.Index,
		Type:     "transfer",
		FeeToken: l.One.TokenId(t.FeeToken).Address,
	}

	if t.Fee != "" && t.Fee != "0" {
		tx.Fee = l.One.Format(t.Fee, l.One.TokenId(t.FeeToken).Decimals)
	}
	return tx
}

func (l *Loopring) DepositToTx(transaction any) Tx {
	var d Deposit
	mapToStruct(transaction, &d)

	return Tx{
		Zero:  l.One.LoopringId(d.ZeroId).Address,
		Value: l.One.Format(d.Value, l.One.TokenId(d.Token).Decimals),
		Token: l.One.TokenId(d.Token).Address,
		Type:  "deposit",
		Index: d.Index,
	}
}

func (l *Loopring) WithdrawToTx(transaction any) Tx {
	var w Withdrawal
	mapToStruct(transaction, &w)

	tx := Tx{
		Zero:     l.One.LoopringId(w.ZeroId).Address,
		Value:    l.One.Format(w.Value, l.One.TokenId(w.Token).Decimals),
		Token:    l.One.TokenId(w.Token).Address,
		Type:     "withdraw",
		Index:    w.Index,
		FeeToken: l.One.TokenId(w.FeeToken).Address,
	}
	switch {
	case w.Fee != "" && w.Fee != "0":
		tx.Fee = l.One.Format(w.Fee, l.One.TokenId(w.FeeToken).Decimals)
	}
	return tx
}

func (l *Loopring) AccountUpdateToTx(transaction any) Tx {
	var a AccountUpdate
	mapToStruct(transaction, &a)
	return Tx{
		Zero:  l.One.LoopringId(a.ZeroId).Address,
		Type:  "accountUpdate",
		Index: a.Index,
		Nonce: a.Nonce,
	}
}

func (l *Loopring) AmmUpdateToTx(transaction any) Tx {
	var a AmmUpdate
	mapToStruct(transaction, &a)
	return Tx{
		Zero:  l.One.LoopringId(a.ZeroId).Address,
		Type:  "ammUpdate",
		Index: a.Index,
		Nonce: a.Nonce,
	}
}

func (l *Loopring) MintToTx(transaction any) Tx {
	var m Mint
	mapToStruct(transaction, &m)
	tx := Tx{
		Zero:     l.One.LoopringId(m.ZeroId).Address,
		Value:    m.Quantity,
		Token:    m.NftAddress,
		Type:     "mint",
		Index:    m.Index,
		FeeToken: l.One.TokenId(m.FeeToken).Address,
	}

	if m.Fee != "" && m.Fee != "0" {
		tx.Fee = l.One.Format(m.Fee, l.One.TokenId(m.FeeToken).Decimals)

	}
	return tx
}

func (l *Loopring) NftDataToTx(transaction any) Tx {
	var n NftData
	mapToStruct(transaction, &n)
	return Tx{
		Zero:  l.One.LoopringId(n.ZeroId).Address,
		Type:  "nft",
		Index: n.Index,
		Raw:   []byte(n.NftData),
	}
}
