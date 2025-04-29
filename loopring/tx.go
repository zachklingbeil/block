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

	// Safely resolve tokens
	tokenZero := l.Value.GetTokenById(s.ZeroToken)
	tokenOne := l.Value.GetTokenById(s.OneToken)

	tokenZeroValue := strconv.FormatInt(s.ZeroToken, 10) // Default to original token ID
	if tokenZero != nil {
		tokenZeroValue = tokenZero.Token
	}

	tokenOneValue := strconv.FormatInt(s.OneToken, 10) // Default to original token ID
	if tokenOne != nil {
		tokenOneValue = tokenOne.Token
	}

	zero := Tx{
		Zero:  strconv.FormatInt(s.Zero, 10),
		One:   strconv.FormatInt(s.One, 10),
		Value: s.ZeroValue,
		Token: tokenZeroValue, // Use resolved token value or fallback
		Fee:   s.ZeroFee,
		Type:  "swap",
		Index: s.Index,
	}

	one := Tx{
		Zero:  strconv.FormatInt(s.Zero, 10),
		One:   strconv.FormatInt(s.One, 10),
		Value: s.OneValue,
		Token: tokenOneValue, // Use resolved token value or fallback
		Fee:   s.OneFee,
		Type:  "swap",
		Index: s.Index,
	}

	return []Tx{zero, one}
}

func (l *Loopring) TransferToTx(transaction any) Tx {
	var t Transfer
	mapToStruct(transaction, &t)

	// Safely resolve tokens
	token := l.Value.GetTokenById(t.Token)
	feeToken := l.Value.GetTokenById(t.FeeToken)

	tokenValue := strconv.FormatInt(t.Token, 10) // Default to original token ID
	if token != nil {
		tokenValue = token.Token
	}

	feeTokenValue := strconv.FormatInt(t.FeeToken, 10) // Default to original token ID
	if feeToken != nil {
		feeTokenValue = feeToken.Token
	}

	return Tx{
		Zero:     strconv.FormatInt(t.ZeroId, 10),
		One:      t.One,
		Value:    t.Value,
		Token:    tokenValue, // Use resolved token value or fallback
		Fee:      t.Fee,
		FeeToken: feeTokenValue, // Use resolved fee token value or fallback
		Index:    t.Index,
	}
}

func (l *Loopring) DepositToTx(transaction any) Tx {
	var d Deposit
	mapToStruct(transaction, &d)

	// Safely resolve the token
	token := l.Value.GetTokenById(d.Token)

	tokenValue := strconv.FormatInt(d.Token, 10) // Default to original token ID
	if token != nil {
		tokenValue = token.Token
	}

	return Tx{
		Zero: strconv.FormatInt(d.ZeroId, 10),
		// One:   d.One,
		Value: d.Value,    // Use resolved token value or fallback
		Token: tokenValue, // Use resolved token value or fallback
		Type:  "deposit",
		Index: d.Index,
	}
}

func (l *Loopring) WithdrawToTx(transaction any) Tx {
	var w Withdrawal
	mapToStruct(transaction, &w)

	// Safely resolve tokens
	token := l.Value.GetTokenById(w.Token)
	feeToken := l.Value.GetTokenById(w.FeeToken)

	tokenValue := strconv.FormatInt(w.Token, 10) // Default to original token ID
	if token != nil {
		tokenValue = token.Token
	}

	feeTokenValue := strconv.FormatInt(w.FeeToken, 10) // Default to original token ID
	if feeToken != nil {
		feeTokenValue = feeToken.Token
	}

	return Tx{
		Zero: strconv.FormatInt(w.ZeroId, 10),
		// One:      w.One,
		Value:    w.Value,
		Token:    tokenValue, // Use resolved token value or fallback
		Fee:      w.Fee,
		FeeToken: feeTokenValue, // Use resolved fee token value or fallback
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
	feeToken := l.Value.GetTokenById(m.FeeToken)
	feeTokenValue := strconv.FormatInt(m.FeeToken, 10) // Default to original token ID
	if feeToken != nil {
		feeTokenValue = feeToken.Token
	}

	return Tx{
		Zero:     m.Zero,
		Value:    m.Quantity,
		Token:    m.NftAddress,
		Fee:      m.Fee,
		FeeToken: feeTokenValue, // Use resolved fee token value or fallback
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
