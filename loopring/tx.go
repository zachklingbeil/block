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
	zero := Tx{
		Zero:  strconv.FormatInt(s.Zero, 10),
		One:   strconv.FormatInt(s.One, 10),
		Value: s.ZeroValue,
		Token: l.Value.GetTokenById(s.ZeroToken),
		Fee:   s.ZeroFee,
		Type:  "swap",
		Index: s.Index,
	}

	one := Tx{
		Zero:  strconv.FormatInt(s.Zero, 10),
		One:   strconv.FormatInt(s.One, 10),
		Value: s.OneValue,
		Token: l.Value.GetTokenById(s.OneToken),
		Fee:   s.OneFee,
		Type:  "swap",
		Index: s.Index,
	}

	return []Tx{zero, one}
}

func (l *Loopring) TransferToTx(transaction any) Tx {
	var t Transfer
	mapToStruct(transaction, &t)

	value, err := l.Value.FormatValue(t.Value, t.Token)
	if err != nil {
		value = t.Value
	}

	fee, err := l.Value.FormatValue(t.Fee, t.FeeToken)
	if err != nil {
		fee = t.Fee
	}

	return Tx{
		Zero:     strconv.FormatInt(t.ZeroId, 10),
		One:      t.One,
		Value:    value,
		Token:    l.Value.GetTokenById(t.Token),
		Fee:      fee,
		FeeToken: l.Value.GetTokenById(t.FeeToken),
		Index:    t.Index,
		Type:     "transfer",
	}
}

func (l *Loopring) DepositToTx(transaction any) Tx {
	var d Deposit
	mapToStruct(transaction, &d)

	value, err := l.Value.FormatValue(d.Value, d.Token)
	if err != nil {
		value = d.Value
	}

	return Tx{
		Zero:  strconv.FormatInt(d.ZeroId, 10),
		Value: value,
		Token: l.Value.GetTokenById(d.Token),
		Type:  "deposit",
		Index: d.Index,
	}
}

func (l *Loopring) WithdrawToTx(transaction any) Tx {
	var w Withdrawal
	mapToStruct(transaction, &w)

	value, err := l.Value.FormatValue(w.Value, w.Token)
	if err != nil {
		value = w.Value
	}

	fee, err := l.Value.FormatValue(w.Fee, w.FeeToken)
	if err != nil {
		fee = w.Fee
	}
	return Tx{
		Zero:     strconv.FormatInt(w.ZeroId, 10),
		Value:    value,
		Token:    l.Value.GetTokenById(w.Token),
		Fee:      fee,
		FeeToken: l.Value.GetTokenById(w.FeeToken),
		Type:     "withdraw",
		Index:    w.Index,
	}
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
	fee, err := l.Value.FormatValue(m.Fee, m.FeeToken)
	if err != nil {
		fee = m.Fee
	}

	return Tx{
		Zero:     m.Zero,
		Value:    m.Quantity,
		Token:    m.NftAddress,
		Fee:      fee,
		FeeToken: l.Value.GetTokenById(m.FeeToken),
		Type:     "mint",
		Index:    m.Index,
	}
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
