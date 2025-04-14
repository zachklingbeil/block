package process

import "encoding/json"

type Block struct {
	Number       int64 `json:"blockId"`
	Size         int64 `json:"blockSize"`
	Timestamp    int64 `json:"createdAt"`
	Transactions []any `json:"transactions"`
}

type Txs struct {
	DW       []DW       `json:"depositWithdraw,omitempty"`
	Swap     []Swap     `json:"swap,omitempty"`
	Transfer []Transfer `json:"transfer,omitempty"`
	Mint     []Mint     `json:"mint,omitempty"`
	TBD      []any      `json:"tbd,omitempty"`
	*json.RawMessage
}

// Depost,  Withdraw (fee)
type DW struct {
	Zero        int64       `json:"accountId"`
	ZeroAddress string      `json:"fromAddress"`
	One         string      `json:"toAddress"`
	Value       string      `json:"token.amount"`
	Token       int64       `json:"token.tokenId"`
	Fee         string      `json:"fee.amount,omitempty"`
	FeeToken    int64       `json:"fee.tokenId,omitempty"`
	Coordinates *Coordinate `json:"coordinates,omitempty"`
}

type Swap struct {
	Zero        int64       `json:"orderA.accountID"`
	ZeroValue   string      `json:"orderA.filledS"`
	ZeroToken   int64       `json:"orderB.tokenB"`
	One         string      `json:"orderB.accountID"`
	OneValue    string      `json:"orderB.filledS"`
	OneToken    int64       `json:"orderA.tokenB"`
	ZeroFee     int64       `json:"orderA.feeBips"`
	OneFee      int64       `json:"fee.orderB.feeBips"`
	Coordinates *Coordinate `json:"coordinates,omitempty"`
}

type Transfer struct {
	Zero        int64       `json:"accountId"`
	One         string      `json:"toAccountId"`
	OneAddress  string      `json:"toAccountAddress"`
	Value       string      `json:"token.amount"`
	Token       int64       `json:"token.tokenId"`
	Fee         string      `json:"fee.amount,omitempty"`
	FeeToken    int64       `json:"fee.tokenId,omitempty"`
	Coordinates *Coordinate `json:"coordinates,omitempty"`
}

type Mint struct {
	Zero        int64       `json:"minterAccountId"`
	ZeroAddress string      `json:"toAccountAddress"`
	Nft         string      `json:"toToken.tokenId"`
	NftId       string      `json:"nftToken.nftId"`
	NftData     string      `json:"nftToken.nftData"`
	NftAddress  string      `json:"nftToken.tokenAddress"`
	Quantity    string      `json:"nftToken.amount"`
	Fee         string      `json:"fee.amount,omitempty"`
	FeeToken    int64       `json:"fee.tokenId,omitempty"`
	Coordinates *Coordinate `json:"coordinates,omitempty"`
}
