package loopring

import (
	"fmt"
	"time"

	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
}

type Output struct {
	Number    int64
	Size      int64
	Timestamp int64
	Coords    string
}

func NewLoopring(factory *factory.Factory) *Loopring {
	return &Loopring{
		Factory: factory,
	}
}

// ProcessInputs converts a slice of Block into a slice of Output
func (l *Loopring) ProcessInputs(in []BlockIn) []Output {
	blocks := make([]Output, len(in))

	for i, block := range in {
		blocks[i] = fx(block)
	}
	return blocks
}

type Coordinates struct {
	Year        int64
	Month       int64
	Day         int64
	Hour        int64
	Minute      int64
	Second      int64
	Millisecond int64
}

// fx processes a single Block into a Output
func fx(in BlockIn) Output {
	t := time.UnixMilli(in.Created)

	// Format the timestamp directly into a string representation of Coordinates
	formattedCoords := fmt.Sprintf("%d.%d.%d.%d.%d.%d.%d",
		t.Year()-2015,      // 0-based year
		int(t.Month()),     // Month
		t.Day(),            // Date of the month (1-31)
		t.Hour(),           // Hour
		t.Minute(),         // Minute
		t.Second(),         // Second
		t.Nanosecond()/1e6) // Millisecond as int64, uncapped

	// Return the Output with the formatted coordinates
	return Output{
		Coords:    formattedCoords,
		Number:    in.Number,
		Size:      in.Size,
		Timestamp: in.Created,
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
