package loopring

import (
	"encoding/json"
	"strconv"
)

func mapToStruct(data any, target any) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, target)
}

func (l *Loopring) TransferToTx(transaction any) Tx {
	var t Transfer
	mapToStruct(transaction, &t)
	token := l.Value.GetTokenById(t.Token).Token
	feeToken := l.Value.GetTokenById(t.FeeToken).Token
	tx := Tx{
		Zero:  l.Value.Hello(strconv.FormatInt(t.ZeroId, 10)),
		One:   l.Value.Hello(t.One),
		Value: l.Value.FormatValue(t.Value, token),
		Token: token,
		Index: t.Index,
		Type:  "transfer",
	}
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
		Zero:  l.Value.Hello(strconv.FormatInt(d.ZeroId, 10)),
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
		Zero:  l.Value.Hello(strconv.FormatInt(w.ZeroId, 10)),
		Value: l.Value.FormatValue(w.Value, token),
		Token: token,
		Type:  "withdraw",
		Index: w.Index,
	}
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
		Zero:  l.Value.Hello(strconv.FormatInt(a.ZeroId, 10)),
		Type:  "accountUpdate",
		Index: a.Index,
		Nonce: a.Nonce,
	}
}

func (l *Loopring) AmmUpdateToTx(transaction any) Tx {
	var a AmmUpdate
	mapToStruct(transaction, &a)
	return Tx{
		Zero:  l.Value.Hello(strconv.FormatInt(a.ZeroId, 10)),
		Type:  "ammUpdate",
		Index: a.Index,
		Nonce: a.Nonce,
	}
}

func (l *Loopring) MintToTx(transaction any) Tx {
	var m Mint
	mapToStruct(transaction, &m)
	feeToken := l.Value.GetTokenById(m.FeeToken).Token

	tx := Tx{
		Zero:  l.Value.Hello(m.Zero),
		Value: m.Quantity,
		Token: m.NftAddress,
		Type:  "mint",
		Index: m.Index,
	}
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
		Zero:  l.Value.Hello(strconv.FormatInt(n.ZeroId, 10)),
		Type:  "nft",
		Index: n.Index,
	}
}
