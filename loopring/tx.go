package loopring

import (
	"encoding/json"
)

func mapToStruct(data any, target any) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, target)
}

func (l *Loopring) TransferToTx(transaction any) Tx {
	var t Transfer
	mapToStruct(transaction, &t)
	token := l.Value.Token.GetAddress(t.Token)
	feeToken := l.Value.Token.GetAddress(t.FeeToken)
	tx := Tx{
		Zero:  l.Value.Peer.GetAddress(t.ZeroId),
		One:   l.Value.Peer.GetAddress(t.OneId),
		Value: l.Value.Token.Format(t.Value, token),
		Token: token,
		Index: t.Index,
		Type:  "transfer",
	}
	if t.Fee != "0" {
		tx.Fee = l.Value.Token.Format(t.Fee, feeToken)
		tx.FeeToken = feeToken
	}
	return tx
}

func (l *Loopring) DepositToTx(transaction any) Tx {
	var d Deposit
	mapToStruct(transaction, &d)
	token := l.Value.Token.GetAddress(d.Token)
	return Tx{
		Zero:  l.Value.Peer.GetAddress(d.ZeroId),
		Value: l.Value.Token.Format(d.Value, token),
		Token: token,
		Type:  "deposit",
		Index: d.Index,
	}
}

func (l *Loopring) WithdrawToTx(transaction any) Tx {
	var w Withdrawal
	mapToStruct(transaction, &w)
	token := l.Value.Token.GetAddress(w.Token)
	feeToken := l.Value.Token.GetAddress(w.FeeToken)
	tx := Tx{
		Zero:  l.Value.Peer.GetAddress(w.ZeroId),
		Value: l.Value.Token.Format(w.Value, token),
		Token: token,
		Type:  "withdraw",
		Index: w.Index,
	}
	if w.Fee != "0" {
		tx.Fee = l.Value.Token.Format(w.Fee, feeToken)
		tx.FeeToken = feeToken
	}
	return tx
}

func (l *Loopring) AccountUpdateToTx(transaction any) Tx {
	var a AccountUpdate
	mapToStruct(transaction, &a)
	return Tx{
		Zero:  l.Value.Peer.GetAddress(a.ZeroId),
		Type:  "accountUpdate",
		Index: a.Index,
		Nonce: a.Nonce,
	}
}

func (l *Loopring) AmmUpdateToTx(transaction any) Tx {
	var a AmmUpdate
	mapToStruct(transaction, &a)
	return Tx{
		Zero:  l.Value.Peer.GetAddress(a.ZeroId),
		Type:  "ammUpdate",
		Index: a.Index,
		Nonce: a.Nonce,
	}
}

func (l *Loopring) MintToTx(transaction any) Tx {
	var m Mint
	mapToStruct(transaction, &m)
	feeToken := l.Value.Token.GetAddress(m.FeeToken)

	tx := Tx{
		Zero:  l.Value.Peer.GetAddress(m.ZeroId),
		Value: m.Quantity,
		Token: m.NftAddress,
		Type:  "mint",
		Index: m.Index,
	}
	if m.Fee != "0" {
		tx.Fee = l.Value.Token.Format(m.Fee, feeToken)
		tx.FeeToken = feeToken
	}
	return tx
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

func (l *Loopring) NftDataToTx(transaction any) Tx {
	var n NftData
	mapToStruct(transaction, &n)
	return Tx{
		Zero:  l.Value.Peer.GetAddress(n.ZeroId),
		Type:  "nft",
		Index: n.Index,
		Raw:   []byte(n.NftData),
	}
}
