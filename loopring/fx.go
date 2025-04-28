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
		Token: s.ZeroToken,
		Fee:   s.ZeroFee,
		Type:  "swap",
		Index: s.Index,
	}

	one := Tx{
		Zero:  strconv.FormatInt(s.Zero, 10),
		One:   strconv.FormatInt(s.One, 10),
		Value: s.OneValue,
		Token: s.OneToken,
		Fee:   s.OneFee,
		Type:  "swap",
		Index: s.Index,
	}
	return []Tx{zero, one}
}

func (l *Loopring) TransferToTx(transaction any) Tx {
	var t Transfer
	mapToStruct(transaction, &t)

	return Tx{
		Zero:     strconv.FormatInt(t.ZeroId, 10),
		One:      t.One,
		Value:    t.Value,
		Token:    t.Token,
		Fee:      t.Fee,
		FeeToken: t.FeeToken,
		// Type:     "transfer",
		Index: t.Index,
	}
}

func (l *Loopring) DepositToTx(transaction any) Tx {
	var d Deposit
	mapToStruct(transaction, &d)
	return Tx{
		Zero:  strconv.FormatInt(d.ZeroId, 10),
		One:   d.One,
		Value: d.Value,
		Token: d.Token,
		Type:  "deposit",
		Index: d.Index,
	}
}

func (l *Loopring) WithdrawToTx(transaction any) Tx {
	var w Withdrawal
	mapToStruct(transaction, &w)
	return Tx{
		Zero:     strconv.FormatInt(w.ZeroId, 10),
		One:      w.One,
		Value:    w.Value,
		Token:    w.Token,
		Fee:      w.Fee,
		FeeToken: w.FeeToken,
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
	return Tx{
		Zero:     m.Zero,
		Value:    m.Quantity,
		Token:    m.NftAddress,
		Fee:      m.Fee,
		FeeToken: m.FeeToken,
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
