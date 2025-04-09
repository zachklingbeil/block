package loopring

import "github.com/zachklingbeil/factory"

type Loopring struct {
	Factory *factory.Factory
}

func NewLoopring(factory *factory.Factory) *Loopring {
	return &Loopring{
		Factory: factory,
	}
}

// type Loopring struct {
// 	TotalNum     int          `json:"totalNum,omitempty"`
// 	Transactions []LoopringTx `json:"transactions,omitempty"`
// }

// type LoopringTx struct {
// 	TxId  int `json:"id,omitempty"`
// 	Block struct {
// 		Number int64 `json:"blockId,omitempty"`
// 		Index  int64 `json:"indexInBlock,omitempty"`
// 	} `json:"blockIdInfo,omitempty"`
// 	Timestamp int64  `json:"timestamp,omitempty"`
// 	Zero      string `json:"senderAddress,omitempty"`
// 	One       string `json:"receiverAddress,omitempty"`
// 	Value     string `json:"amount,omitempty"`
// 	Token     string `json:"symbol,omitempty"`
// 	FeeToken  string `json:"feeTokenSymbol,omitempty"`
// 	FeeValue  string `json:"feeAmount,omitempty"`
// }

// type Transaction struct {
// 	TxType           TxType   `json:"txType"`
// 	AccountID        *int64   `json:"accountId,omitempty"`
// 	Token            *Token   `json:"token,omitempty"`
// 	ToToken          *ToToken `json:"toToken,omitempty"`
// 	Fee              *Fee     `json:"fee,omitempty"`
// 	ValidUntil       *int64   `json:"validUntil,omitempty"`
// 	ToAccountID      *int64   `json:"toAccountId,omitempty"`
// 	ToAccountAddress *string  `json:"toAccountAddress,omitempty"`
// 	StorageID        *int64   `json:"storageId,omitempty"`
// 	OrderA           *Order   `json:"orderA,omitempty"`
// 	OrderB           *Order   `json:"orderB,omitempty"`
// 	Valid            *bool    `json:"valid,omitempty"`
// 	Owner            *string  `json:"owner,omitempty"`
// 	FromAddress      *string  `json:"fromAddress,omitempty"`
// 	ToAddress        *string  `json:"toAddress,omitempty"`
// }

// type Fee struct {
// 	TokenID int64  `json:"tokenId"`
// 	Amount  string `json:"amount"`
// }

// type ToToken struct {
// 	TokenID int64 `json:"tokenId"`
// }

// type Token struct {
// 	TokenID int64   `json:"tokenId"`
// 	NftData *string `json:"nftData,omitempty"`
// 	Amount  string  `json:"amount"`
// }

// type TxType string

// const (
// 	Deposit   TxType = "Deposit"
// 	SpotTrade TxType = "SpotTrade"
// 	Transfer  TxType = "Transfer"
// )
// type Order struct {
// 	StorageID  int64  `json:"storageID"`
// 	AccountID  int64  `json:"accountID"`
// 	AmountS    string `json:"amountS"`
// 	AmountB    string `json:"amountB"`
// 	TokenS     int64  `json:"tokenS"`
// 	TokenB     int64  `json:"tokenB"`
// 	ValidUntil int64  `json:"validUntil"`
// 	Taker      string `json:"taker"`
// 	FeeBips    int64  `json:"feeBips"`
// 	IsAmm      bool   `json:"isAmm"`
// 	NftData    string `json:"nftData"`
// 	FillS      int64  `json:"fillS"`
// 	FilledS    string `json:"filledS"`
// }
