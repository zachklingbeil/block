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
	zero := l.One.LoopringId(s.Zero)
	one := l.One.LoopringId(s.One)
	token := l.One.TokenId(s.Token)
	tokenOut := l.One.TokenId(s.TokenOut)

	value := l.One.Format(s.Value, token.Decimals)
	valueOut := l.One.Format(s.ValueOut, tokenOut.Decimals)

	var fee string
	var feeTokenAddr string
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
		feeInt := new(big.Int).Mul(valueIn, big.NewInt(feeBips))
		feeInt.Div(feeInt, big.NewInt(10000)) // Convert basis points to percentage
		fee = l.One.Format(feeInt.String(), token.Decimals)
		feeTokenAddr = token.Address
	}

	return Tx{
		Zero:     zero.Address,
		One:      one.Address,
		Value:    value,
		ValueOut: valueOut,
		Token:    token.Address,
		TokenOut: tokenOut.Address,
		Type:     "swap",
		Index:    s.Index,
		Fee:      fee,
		FeeToken: feeTokenAddr,
	}
}

func (l *Loopring) TransferToTx(transaction any) Tx {
	var t Transfer
	mapToStruct(transaction, &t)

	zero := l.One.LoopringId(t.ZeroId)
	one := l.One.LoopringId(t.OneId)
	token := l.One.TokenId(t.Token)
	feeToken := l.One.TokenId(t.FeeToken)

	value := l.One.Format(t.Value, token.Decimals)
	fee := ""
	feeTokenAddr := ""
	if t.Fee != "" && t.Fee != "0" {
		fee = l.One.Format(t.Fee, feeToken.Decimals)
		feeTokenAddr = feeToken.Address
	}

	return Tx{
		Zero:     zero.Address,
		One:      one.Address,
		Value:    value,
		Token:    token.Address,
		Index:    t.Index,
		Type:     "transfer",
		Fee:      fee,
		FeeToken: feeTokenAddr,
	}
}

func (l *Loopring) DepositToTx(transaction any) Tx {
	var d Deposit
	mapToStruct(transaction, &d)
	zero := l.One.LoopringId(d.ZeroId)
	token := l.One.TokenId(d.Token)

	value := l.One.Format(d.Value, token.Decimals)
	return Tx{
		Zero:  zero.Address,
		Value: value,
		Token: token.Address,
		Type:  "deposit",
		Index: d.Index,
	}
}

func (l *Loopring) WithdrawToTx(transaction any) Tx {
	var w Withdrawal
	mapToStruct(transaction, &w)
	zero := l.One.LoopringId(w.ZeroId)
	token := l.One.TokenId(w.Token)
	feeToken := l.One.TokenId(w.FeeToken)

	value := l.One.Format(w.Value, token.Decimals)
	fee := ""
	feeTokenAddr := ""
	if w.Fee != "" && w.Fee != "0" {
		fee = l.One.Format(w.Fee, feeToken.Decimals)
		feeTokenAddr = feeToken.Address
	}

	return Tx{
		Zero:     zero.Address,
		Value:    value,
		Token:    token.Address,
		Type:     "withdraw",
		Index:    w.Index,
		Fee:      fee,
		FeeToken: feeTokenAddr,
	}
}

func (l *Loopring) AccountUpdateToTx(transaction any) Tx {
	var a AccountUpdate
	mapToStruct(transaction, &a)
	zero := l.One.LoopringId(a.ZeroId)
	return Tx{
		Zero:  zero.Address,
		Type:  "accountUpdate",
		Index: a.Index,
		Nonce: a.Nonce,
	}
}

func (l *Loopring) AmmUpdateToTx(transaction any) Tx {
	var a AmmUpdate
	mapToStruct(transaction, &a)
	zero := l.One.LoopringId(a.ZeroId)
	return Tx{
		Zero:  zero.Address,
		Type:  "ammUpdate",
		Index: a.Index,
		Nonce: a.Nonce,
	}
}

func (l *Loopring) MintToTx(transaction any) Tx {
	var m Mint
	mapToStruct(transaction, &m)
	zero := l.One.LoopringId(m.ZeroId)
	feeToken := l.One.TokenId(m.FeeToken)

	value := l.One.Format(m.Quantity, feeToken.Decimals)
	fee := ""
	feeTokenAddr := ""
	if m.Fee != "" && m.Fee != "0" {
		fee = l.One.Format(m.Fee, feeToken.Decimals)
		feeTokenAddr = feeToken.Address
	}

	return Tx{
		Zero:     zero.Address,
		Value:    value,
		Token:    m.NftAddress,
		Type:     "mint",
		Index:    m.Index,
		Fee:      fee,
		FeeToken: feeTokenAddr,
	}
}

func (l *Loopring) NftDataToTx(transaction any) Tx {
	var n NftData
	mapToStruct(transaction, &n)
	zero := l.One.LoopringId(n.ZeroId)
	return Tx{
		Zero:  zero.Address,
		Type:  "nft",
		Index: n.Index,
		Raw:   []byte(n.NftData),
	}
}
