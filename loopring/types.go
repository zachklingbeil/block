package loopring

import (
	"encoding/json"
	"fmt"
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
	For      string `json:"orderB.filledS,omitempty"`
	ForToken int64  `json:"orderA.tokenB,omitempty"`
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

func (l *Loopring) Token(tokenId int64) string {
	one := l.Zero.TokenId(tokenId)
	if one != nil && one.Token != "" {
		return one.Token
	}
	return fmt.Sprintf("%d", tokenId)
}
func (l *Loopring) Who(id int64) string {
	peer := l.Zero.LoopringId(id)
	if peer == nil {
		return ""
	}
	if peer.ENS != "" && peer.ENS != "." {
		return peer.ENS
	}
	if peer.LoopringENS != "" && peer.LoopringENS != "." && peer.LoopringENS != "!" {
		return peer.LoopringENS
	}
	return peer.Address
}

func (l *Loopring) Decimals(tokenId int64) int64 {
	one := l.Zero.TokenId(tokenId)
	if one != nil && one.Decimals != 0 {
		return one.Decimals
	}
	return 18
}
func (l *Loopring) SwapToTx(transaction any) universe.Tx {
	var s Swap
	mapToStruct(transaction, &s)
	tx := universe.Tx{
		Zero:     l.Who(s.Zero),
		One:      l.Who(s.One),
		Value:    l.Zero.Format.Value(s.Value, l.Decimals(s.Token)),
		Token:    l.Token(s.Token),
		For:      l.Zero.Format.Value(s.For, l.Decimals(s.ForToken)),
		ForToken: l.Token(s.ForToken),
		Index:    s.Index,
	}

	switch {
	case s.ZeroFee != 0:
		fee := calcFee(s.Value, s.ZeroFee)
		tx.Fee = l.Zero.Format.Value(fee, l.Decimals(s.Token))
	case s.OneFee != 0:
		fee := calcFee(s.Value, s.OneFee)
		tx.Fee = l.Zero.Format.Value(fee, l.Decimals(s.Token))
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
		Zero:  l.Who(t.ZeroId),
		One:   l.Who(t.OneId),
		Value: l.Zero.Format.Value(t.Value, l.Decimals(t.Token)),
		Token: l.Token(t.Token),
		Index: t.Index,
		// Type:  "transfer",
	}

	if t.Fee != "" && t.Fee != "0" {
		tx.Fee = l.Zero.Format.Value(t.Fee, l.Decimals(t.FeeToken))
		tx.FeeToken = l.Token(t.FeeToken)
	}
	return tx
}

func (l *Loopring) DepositToTx(transaction any) universe.Tx {
	var d Deposit
	mapToStruct(transaction, &d)

	return universe.Tx{
		Zero:  l.Who(d.ZeroId),
		Value: l.Zero.Format.Value(d.Value, l.Decimals(d.Token)),
		Token: l.Token(d.Token),
		Type:  "deposit",
		Index: d.Index,
	}
}

func (l *Loopring) WithdrawToTx(transaction any) universe.Tx {
	var w Withdrawal
	mapToStruct(transaction, &w)

	tx := universe.Tx{
		Zero:     l.Who(w.ZeroId),
		Value:    l.Zero.Format.Value(w.Value, l.Decimals(w.Token)),
		Token:    l.Token(w.Token),
		Type:     "withdraw",
		Index:    w.Index,
		FeeToken: l.Token(w.FeeToken),
	}
	switch {
	case w.Fee != "" && w.Fee != "0":
		tx.Fee = l.Zero.Format.Value(w.Fee, l.Decimals(w.FeeToken))
	}
	return tx
}

func (l *Loopring) AccountUpdateToTx(transaction any) universe.Tx {
	var a AccountUpdate
	mapToStruct(transaction, &a)
	return universe.Tx{
		Zero:  l.Who(a.ZeroId),
		Type:  "accountUpdate",
		Index: a.Index,
		Nonce: a.Nonce,
	}
}

func (l *Loopring) AmmUpdateToTx(transaction any) universe.Tx {
	var a AmmUpdate
	mapToStruct(transaction, &a)
	return universe.Tx{
		Zero:  l.Who(a.ZeroId),
		Type:  "ammUpdate",
		Index: a.Index,
		Nonce: a.Nonce,
	}
}

func (l *Loopring) MintToTx(transaction any) universe.Tx {
	var m Mint
	mapToStruct(transaction, &m)
	tx := universe.Tx{
		Zero:     l.Who(m.ZeroId),
		Value:    m.Quantity,
		Token:    m.NftAddress,
		Type:     "mint",
		Index:    m.Index,
		FeeToken: l.Token(m.FeeToken),
	}

	if m.Fee != "" && m.Fee != "0" {
		tx.Fee = l.Zero.Format.Value(m.Fee, l.Decimals(m.FeeToken))
	}
	return tx
}

func (l *Loopring) NftDataToTx(transaction any) universe.Tx {
	var n NftData
	mapToStruct(transaction, &n)
	return universe.Tx{
		Zero:  l.Who(n.ZeroId),
		Type:  "nft",
		Index: n.Index,
		Raw:   []byte(n.NftData),
	}
}
