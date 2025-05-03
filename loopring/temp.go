package loopring

import (
	"encoding/json"
	"math/big"
	"strconv"
)

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
	AMM      bool   `json:"orderB.isAmm,omitempty"`
}

func (l *Loopring) SwapToTx(transaction any) Tx {
	var s Swap
	mapToStruct(transaction, &s)
	tokenIn := l.Value.GetTokenById(strconv.FormatInt(s.Token, 10)).Token
	tokenOut := l.Value.GetTokenById(strconv.FormatInt(s.TokenOut, 10)).Token

	tx := Tx{
		Zero:     l.Value.Hello(strconv.FormatInt(s.Zero, 10)),
		One:      l.Value.Hello(strconv.FormatInt(s.One, 10)),
		Value:    l.Value.FormatValue(s.Value, tokenIn),
		ValueOut: l.Value.FormatValue(s.ValueOut, tokenOut),
		Token:    tokenIn,
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
		tx.Fee = l.Value.FormatValue(fee.String(), tokenIn)
		tx.FeeToken = tokenIn
	}

	return tx
}
