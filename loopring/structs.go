package loopring

import "encoding/json"

type Tx struct {
	Zero        any             `json:"zero,omitempty"`
	One         any             `json:"one,omitempty"`
	Value       string          `json:"value,omitempty"`
	Token       any             `json:"token,omitempty"`
	Fee         any             `json:"fee,omitempty"`
	FeeToken    int64           `json:"feeToken,omitempty"`
	OneValue    string          `json:"oneValue,omitempty"`
	OneToken    int64           `json:"oneToken,omitempty"`
	OneFee      any             `json:"oneFee,omitempty"`
	OneFeeToken int64           `json:"oneFeeToken,omitempty"`
	Type        string          `json:"type,omitempty"`
	Coordinates Coordinate      `json:"coordinates"`
	Index       int64           `json:"index"`
	Raw         json.RawMessage `json:"raw,omitempty"`
}

type Coordinate struct {
	Block       int64 `json:"block"`
	Year        int64 `json:"year"`
	Month       int64 `json:"month"`
	Day         int64 `json:"day"`
	Hour        int64 `json:"hour"`
	Minute      int64 `json:"minute"`
	Second      int64 `json:"second"`
	Millisecond int64 `json:"millisecond"`
	Index       int64 `json:"index"`
}

type Swap struct {
	ZeroId      int64      `json:"orderA.accountID"`
	ZeroValue   string     `json:"orderA.filledS"`
	ZeroToken   int64      `json:"orderB.tokenB"`
	OneId       int64      `json:"orderB.accountID"`
	OneValue    string     `json:"orderB.filledS"`
	OneToken    int64      `json:"orderA.tokenB"`
	ZeroFee     int64      `json:"orderA.feeBips"`
	OneFee      int64      `json:"fee.orderB.feeBips"`
	Type        string     `json:"txType,omitempty"`
	Coordinates Coordinate `json:"coordinates"`
	Index       int64      `json:"index"`
}

type Transfer struct {
	ZeroId      int64      `json:"accountId"`
	OneId       int64      `json:"toAccountId"`
	One         string     `json:"toAccountAddress"`
	Value       string     `json:"token.amount"`
	Token       int64      `json:"token.tokenId"`
	Fee         string     `json:"fee.amount,omitempty"`
	FeeToken    int64      `json:"fee.tokenId,omitempty"`
	Type        string     `json:"txType,omitempty"`
	Coordinates Coordinate `json:"coordinates"`
	Index       int64      `json:"index"`
}

// Depost,  Withdraw (fee)
type DW struct {
	Zero        string     `json:"fromAddress"`
	ZeroId      int64      `json:"accountId"`
	One         string     `json:"toAddress"`
	Value       string     `json:"token.amount"`
	Token       int64      `json:"token.tokenId"`
	Fee         string     `json:"fee.amount,omitempty"`
	FeeToken    int64      `json:"fee.tokenId,omitempty"`
	Type        string     `json:"txType,omitempty"`
	Coordinates Coordinate `json:"coordinates"`
	Index       int64      `json:"index"`
}

type AccountUpdate struct {
	ZeroId      int64      `json:"accountId"`
	Type        string     `json:"txType,omitempty"`
	Coordinates Coordinate `json:"coordinates"`
	Index       int64      `json:"index"`
}

type AmmUpdate struct {
	Zero        string     `json:"owner"`
	ZeroId      int64      `json:"accountId"`
	Nonce       int64      `json:"nonce"`
	Type        string     `json:"txType,omitempty"`
	Coordinates Coordinate `json:"coordinates"`
	Index       int64      `json:"index"`
}

type Mint struct {
	ZeroId      int64      `json:"minterAccountId"`
	Zero        string     `json:"toAccountAddress"`
	Nft         any        `json:"toToken.tokenId"`
	NftId       string     `json:"nftToken.nftId"`
	NftData     string     `json:"nftToken.nftData"`
	NftAddress  string     `json:"nftToken.tokenAddress"`
	Quantity    string     `json:"nftToken.amount"`
	Fee         string     `json:"fee.amount,omitempty"`
	FeeToken    int64      `json:"fee.tokenId,omitempty"`
	Type        string     `json:"txType,omitempty"`
	Coordinates Coordinate `json:"coordinates"`
	Index       int64      `json:"index"`
}

type NftData struct {
	ZeroId      int64      `json:"accountId"`
	One         string     `json:"minter"`
	NftId       string     `json:"nftToken.nftId"`
	NftData     string     `json:"nftToken.nftData,omitempty"`
	NftAddress  string     `json:"nftToken.tokenAddress"`
	Type        string     `json:"txType,omitempty"`
	Coordinates Coordinate `json:"coordinates"`
	Index       int64      `json:"index"`
}
