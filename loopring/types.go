package loopring

import "encoding/json"

type Raw struct {
	Number       int64 `json:"blockId"`
	Timestamp    int64 `json:"createdAt"`
	Size         int64 `json:"blockSize"`
	Transactions []any `json:"transactions"`
}

type Block struct {
	Number int64      `json:"block"`
	Zero   Coordinate `json:"zero"`
	Ones   []Tx       `json:"one"`
}

type Coordinate struct {
	Year        uint8  `json:"year"`
	Month       uint8  `json:"month"`
	Day         uint8  `json:"day"`
	Hour        uint8  `json:"hour"`
	Minute      uint8  `json:"minute"`
	Second      uint8  `json:"second"`
	Millisecond uint16 `json:"millisecond"`
	Index       uint16 `json:"index"`
	Depth       uint16 `json:"depth,omitempty"`
}

type Tx struct {
	Zero     any             `json:"zero,omitempty"`
	One      any             `json:"one,omitempty"`
	Value    any             `json:"value,omitempty"`
	Token    any             `json:"token,omitempty"`
	Fee      any             `json:"fee,omitempty"`
	FeeToken any             `json:"feeToken,omitempty"`
	Type     any             `json:"type,omitempty"`
	Index    any             `json:"index"`
	Raw      json.RawMessage `json:"raw,omitempty"`
}

type Transfer struct {
	ZeroId   int64  `json:"accountId"`
	OneId    int64  `json:"toAccountId"`
	One      string `json:"toAccountAddress"`
	Value    string `json:"token.amount"`
	Token    int64  `json:"token.tokenId"`
	Fee      string `json:"fee.amount,omitempty"`
	FeeToken int64  `json:"fee.tokenId,omitempty"`
	Type     string `json:"txType,omitempty"`
	Index    uint16 `json:"index"`
}

// Depost,  Withdraw (fee)
type Deposit struct {
	Zero   string `json:"fromAddress"`
	ZeroId int64  `json:"accountId"`
	One    string `json:"toAddress"`
	Value  string `json:"token.amount"`
	Token  int64  `json:"token.tokenId"`
	Type   string `json:"txType,omitempty"`
	Index  uint16 `json:"index"`
}

type Withdrawal struct {
	Zero     string `json:"fromAddress"`
	ZeroId   int64  `json:"accountId"`
	One      string `json:"toAddress"`
	Value    string `json:"token.amount"`
	Token    int64  `json:"token.tokenId"`
	Fee      string `json:"fee.amount,omitempty"`
	FeeToken int64  `json:"fee.tokenId,omitempty"`
	Type     string `json:"txType,omitempty"`
	Index    uint16 `json:"index"`
}

type AccountUpdate struct {
	ZeroId int64  `json:"accountId"`
	Type   string `json:"txType,omitempty"`
	Index  uint16 `json:"index"`
}

type AmmUpdate struct {
	Zero   string `json:"owner"`
	ZeroId int64  `json:"accountId"`
	Nonce  int64  `json:"nonce"`
	Type   string `json:"txType,omitempty"`
	Index  uint16 `json:"index"`
}

type Mint struct {
	ZeroId     int64  `json:"minterAccountId"`
	Zero       string `json:"toAccountAddress"`
	Nft        any    `json:"toToken.tokenId"`
	NftId      string `json:"nftToken.nftId"`
	NftData    string `json:"nftToken.nftData"`
	NftAddress string `json:"nftToken.tokenAddress"`
	Quantity   string `json:"nftToken.amount"`
	Fee        string `json:"fee.amount,omitempty"`
	FeeToken   int64  `json:"fee.tokenId,omitempty"`
	Type       string `json:"txType,omitempty"`
	Index      uint16 `json:"index"`
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

type SpotTrade struct {
	Zero      int64  `json:"orderA.accountID"`
	ZeroValue string `json:"orderA.filledS"`
	ZeroToken int64  `json:"orderB.tokenB"`
	One       int64  `json:"orderB.accountID"`
	OneValue  string `json:"orderB.filledS"`
	OneToken  int64  `json:"orderA.tokenB"`
	ZeroFee   int64  `json:"orderA.feeBips"`
	OneFee    int64  `json:"orderB.feeBips"`
	Type      string `json:"txType,omitempty"`
	Index     uint16 `json:"index"`
}
