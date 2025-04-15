package process

type Block struct {
	Number       int64 `json:"blockId"`
	Size         int64 `json:"blockSize"`
	Timestamp    int64 `json:"createdAt"`
	Transactions []any `json:"transactions"`
}

type RawTx struct {
	// Common fields
	Zero        int64       `json:"accountId,omitempty"`
	ZeroAddress string      `json:"fromAddress,omitempty"`
	One         string      `json:"toAddress,omitempty"`
	OneAddress  string      `json:"toAccountAddress,omitempty"`
	Value       string      `json:"token.amount,omitempty"`
	Token       int64       `json:"token.tokenId,omitempty"`
	Fee         string      `json:"fee.amount,omitempty"`
	FeeToken    int64       `json:"fee.tokenId,omitempty"`
	Coordinates *Coordinate `json:"coordinates,omitempty"`

	// Swap-specific fields
	ZeroSwap  int64  `json:"orderA.accountID,omitempty"`
	ZeroValue string `json:"orderA.filledS,omitempty"`
	ZeroToken int64  `json:"orderB.tokenB,omitempty"`
	ZeroFee   int64  `json:"orderA.feeBips,omitempty"`
	OneSwap   int64  `json:"orderB.accountID,omitempty"`
	OneValue  string `json:"orderB.filledS,omitempty"`
	OneToken  int64  `json:"orderA.tokenB,omitempty"`
	OneFee    int64  `json:"fee.orderB.feeBips,omitempty"`

	// Mint-specific fields
	NftId      string `json:"nftToken.nftId,omitempty"`
	NftData    string `json:"nftToken.nftData,omitempty"`
	NftAddress string `json:"nftToken.tokenAddress,omitempty"`
	Quantity   string `json:"nftToken.amount,omitempty"`

	// AmmUpdate-specific fields
	Nonce int64 `json:"nonce,omitempty"`

	// Type field to identify the transaction type
	Type string `json:"txType,omitempty"`
}

type Tx struct {
	Zero string `json:"zero,omitempty"`
	One  string `json:"one,omitempty"`

	Value    string `json:"value,omitempty"`
	OneValue string `json:"oneValue,omitempty"`

	Token    int64 `json:"token,omitempty"`
	OneToken int64 `json:"oneToken,omitempty"`

	Fee         string `json:"fee,omitempty"`
	FeeToken    int64  `json:"feeToken,omitempty"`
	OneFee      string `json:"oneFee,omitempty"`
	OneFeeToken int64  `json:"oneFeeToken,omitempty"`

	Type string `json:"type,omitempty"`
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

// type Tx struct {
// 	Zero        int64       `json:"zero,omitempty"`
// 	ZeroId      string      `json:"zeroId,omitempty"`
// 	One         string      `json:"one,omitempty"`
// 	OneId       int64       `json:"oneId,omitempty"`
// 	Value       string      `json:"value,omitempty"`
// 	Token       int64       `json:"token,omitempty"`
// 	TokenOne    int64       `json:"tokenOne,omitempty"`
// 	Fee         string      `json:"fee,omitempty"`
// 	FeeToken    int64       `json:"feeToken,omitempty"`
// 	Type        string      `json:"type,omitempty"`
// 	Coordinates *Coordinate `json:"coordinates,omitempty"`
// 	*json.RawMessage
// }
