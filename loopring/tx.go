package loopring

import (
	"encoding/json"
	"strconv"
)

func mapToStruct(data any, target any) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, target)
}

func (l *Loopring) SwapToTx(transaction any) []Tx {
	var s SpotTrade
	mapToStruct(transaction, &s)
	tokenZero := l.Value.GetTokenById(s.ZeroToken).Token
	zero := Tx{
		Zero:  strconv.FormatInt(s.Zero, 10),
		One:   strconv.FormatInt(s.One, 10),
		Value: l.Value.FormatValue(s.ZeroValue, tokenZero),
		Token: tokenZero,
		Type:  "swap",
		Index: s.Index,
	}

	if s.ZeroFee != 0 {
		zero.Fee = s.ZeroFee
	}

	tokenOne := l.Value.GetTokenById(s.OneToken).Token
	one := Tx{
		Zero:  strconv.FormatInt(s.Zero, 10),
		One:   strconv.FormatInt(s.One, 10),
		Value: l.Value.FormatValue(s.OneValue, tokenOne),
		Token: tokenOne,
		Type:  "swap",
		Index: s.Index,
	}

	if s.OneFee != 0 {
		one.Fee = s.OneFee
	}

	return []Tx{zero, one}
}

func (l *Loopring) TransferToTx(transaction any) Tx {
	var t Transfer
	mapToStruct(transaction, &t)
	token := l.Value.GetTokenById(t.Token).Token
	feeToken := l.Value.GetTokenById(t.FeeToken).Token

	tx := Tx{
		Zero:  strconv.FormatInt(t.ZeroId, 10),
		One:   t.One,
		Value: l.Value.FormatValue(t.Value, token),
		Token: token,
		Index: t.Index,
		Type:  "transfer",
	}

	// Conditionally include Fee and FeeToken
	if t.Fee != "0" {
		tx.Fee = l.Value.FormatValue(t.Fee, feeToken)
		tx.FeeToken = feeToken
	}

	return tx
}

func (l *Loopring) DepositToTx(transaction any) Tx {
	var d Deposit
	mapToStruct(transaction, &d)
	token := l.Value.GetTokenById(d.Token).Token

	return Tx{
		Zero:  strconv.FormatInt(d.ZeroId, 10),
		Value: l.Value.FormatValue(d.Value, token),
		Token: token,
		Type:  "deposit",
		Index: d.Index,
	}
}

func (l *Loopring) WithdrawToTx(transaction any) Tx {
	var w Withdrawal
	mapToStruct(transaction, &w)
	token := l.Value.GetTokenById(w.Token).Token
	feeToken := l.Value.GetTokenById(w.FeeToken).Token

	tx := Tx{
		Zero:  strconv.FormatInt(w.ZeroId, 10),
		Value: l.Value.FormatValue(w.Value, token),
		Token: token,
		Type:  "withdraw",
		Index: w.Index,
	}

	// Conditionally include Fee and FeeToken
	if w.Fee != "0" {
		tx.Fee = l.Value.FormatValue(w.Fee, feeToken)
		tx.FeeToken = feeToken
	}

	return tx
}

func (l *Loopring) AccountUpdateToTx(transaction any) Tx {
	var a AccountUpdate
	mapToStruct(transaction, &a)
	return Tx{
		Zero:  strconv.FormatInt(a.ZeroId, 10),
		Type:  "accountUpdate",
		Index: a.Index,
	}
}

func (l *Loopring) AmmUpdateToTx(transaction any) Tx {
	var a AmmUpdate
	mapToStruct(transaction, &a)
	return Tx{
		Zero:  strconv.FormatInt(a.ZeroId, 10),
		Type:  "ammUpdate",
		Index: a.Index,
	}
}

func (l *Loopring) MintToTx(transaction any) Tx {
	var m Mint
	mapToStruct(transaction, &m)
	feeToken := l.Value.GetTokenById(m.FeeToken).Token

	tx := Tx{
		Zero:  m.Zero,
		Value: m.Quantity,
		Token: m.NftAddress,
		Type:  "mint",
		Index: m.Index,
	}

	// Conditionally include Fee and FeeToken
	if m.Fee != "0" {
		tx.Fee = l.Value.FormatValue(m.Fee, feeToken)
		tx.FeeToken = feeToken
	}

	return tx
}

func (l *Loopring) NftDataToTx(transaction any) Tx {
	var n NftData
	mapToStruct(transaction, &n)
	return Tx{
		Zero:  strconv.FormatInt(n.ZeroId, 10),
		Type:  "nft",
		Index: n.Index,
	}
}
