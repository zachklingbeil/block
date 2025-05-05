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
}

func (l *Loopring) SwapToTx(transaction any) Tx {
	var s Swap
	mapToStruct(transaction, &s)
	tokenIn := l.Value.Token.Get(s.Token, 1)
	tokenOut := l.Value.Token.Get(s.TokenOut, 1)
	tx := Tx{
		Zero:     l.Value.Peer.Hello(strconv.FormatInt(s.Zero, 10)),
		One:      l.Value.Peer.Hello(strconv.FormatInt(s.One, 10)),
		Value:    l.Value.Token.Format(s.Value, s.Token),
		ValueOut: l.Value.Token.Format(s.ValueOut, s.TokenOut),
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
		tx.Fee = l.Value.Token.Format(fee.String(), s.Token)
		tx.FeeToken = tokenIn
	}
	return tx
}
