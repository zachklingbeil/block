package loopring

import (
	"encoding/json"
	"math/big"

	"github.com/zachklingbeil/block/universe"
)

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

func (l *Loopring) SwapToTx(transaction any) universe.Tx {
	var s Swap
	mapToStruct(transaction, &s)
	tx := universe.Tx{
		Zero:     l.Zero.LoopringId(s.Zero).Address,
		One:      l.Zero.LoopringId(s.One).Address,
		Value:    l.Zero.Format(s.Value, l.Zero.TokenId(s.Token).Decimals),
		ValueOut: l.Zero.Format(s.ValueOut, l.Zero.TokenId(s.TokenOut).Decimals),
		Token:    l.Zero.TokenId(s.Token).Address,
		TokenOut: l.Zero.TokenId(s.TokenOut).Address,
		Type:     "swap",
		Index:    s.Index,
		FeeToken: l.Zero.TokenId(s.Token).Address,
	}

	switch {
	case s.ZeroFee != 0:
		fee := calcFee(s.Value, s.ZeroFee)
		tx.Fee = l.Zero.Format(fee, l.Zero.TokenId(s.Token).Decimals)
	case s.OneFee != 0:
		fee := calcFee(s.Value, s.OneFee)
		tx.Fee = l.Zero.Format(fee, l.Zero.TokenId(s.Token).Decimals)
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

func (l *Loopring) TransferToTx(transaction any) universe.Tx {
	var t Transfer
	mapToStruct(transaction, &t)

	tx := universe.Tx{
		Zero:     l.Zero.LoopringId(t.ZeroId).Address,
		One:      l.Zero.LoopringId(t.OneId).Address,
		Value:    l.Zero.Format(t.Value, l.Zero.TokenId(t.Token).Decimals),
		Token:    l.Zero.TokenId(t.Token).Address,
		Index:    t.Index,
		Type:     "transfer",
		FeeToken: l.Zero.TokenId(t.FeeToken).Address,
	}

	if t.Fee != "" && t.Fee != "0" {
		tx.Fee = l.Zero.Format(t.Fee, l.Zero.TokenId(t.FeeToken).Decimals)
	}
	return tx
}

func (l *Loopring) DepositToTx(transaction any) universe.Tx {
	var d Deposit
	mapToStruct(transaction, &d)

	return universe.Tx{
		Zero:  l.Zero.LoopringId(d.ZeroId).Address,
		Value: l.Zero.Format(d.Value, l.Zero.TokenId(d.Token).Decimals),
		Token: l.Zero.TokenId(d.Token).Address,
		Type:  "deposit",
		Index: d.Index,
	}
}

func (l *Loopring) WithdrawToTx(transaction any) universe.Tx {
	var w Withdrawal
	mapToStruct(transaction, &w)

	tx := universe.Tx{
		Zero:     l.Zero.LoopringId(w.ZeroId).Address,
		Value:    l.Zero.Format(w.Value, l.Zero.TokenId(w.Token).Decimals),
		Token:    l.Zero.TokenId(w.Token).Address,
		Type:     "withdraw",
		Index:    w.Index,
		FeeToken: l.Zero.TokenId(w.FeeToken).Address,
	}
	switch {
	case w.Fee != "" && w.Fee != "0":
		tx.Fee = l.Zero.Format(w.Fee, l.Zero.TokenId(w.FeeToken).Decimals)
	}
	return tx
}

func (l *Loopring) AccountUpdateToTx(transaction any) universe.Tx {
	var a AccountUpdate
	mapToStruct(transaction, &a)
	return universe.Tx{
		Zero:  l.Zero.LoopringId(a.ZeroId).Address,
		Type:  "accountUpdate",
		Index: a.Index,
		Nonce: a.Nonce,
	}
}

func (l *Loopring) AmmUpdateToTx(transaction any) universe.Tx {
	var a AmmUpdate
	mapToStruct(transaction, &a)
	return universe.Tx{
		Zero:  l.Zero.LoopringId(a.ZeroId).Address,
		Type:  "ammUpdate",
		Index: a.Index,
		Nonce: a.Nonce,
	}
}

func (l *Loopring) MintToTx(transaction any) universe.Tx {
	var m Mint
	mapToStruct(transaction, &m)
	tx := universe.Tx{
		Zero:     l.Zero.LoopringId(m.ZeroId).Address,
		Value:    m.Quantity,
		Token:    m.NftAddress,
		Type:     "mint",
		Index:    m.Index,
		FeeToken: l.Zero.TokenId(m.FeeToken).Address,
	}

	if m.Fee != "" && m.Fee != "0" {
		tx.Fee = l.Zero.Format(m.Fee, l.Zero.TokenId(m.FeeToken).Decimals)

	}
	return tx
}

func (l *Loopring) NftDataToTx(transaction any) universe.Tx {
	var n NftData
	mapToStruct(transaction, &n)
	return universe.Tx{
		Zero:  l.Zero.LoopringId(n.ZeroId).Address,
		Type:  "nft",
		Index: n.Index,
		Raw:   []byte(n.NftData),
	}
}
