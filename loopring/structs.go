package loopring

type Block struct {
	Number       int64         `json:"blockId"`
	Size         int64         `json:"blockSize"`
	TxHash       string        `json:"txHash"`
	Created      int64         `json:"createdAt"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	TxType     TxType `json:"txType"`
	Token      Token  `json:"token"`
	ToToken    Token  `json:"toToken"`
	Fee        Fee    `json:"fee"`
	FromID     int64  `json:"accountId"`
	ToID       int64  `json:"toAccountId"`
	ToAddress  string `json:"toAccountAddress"`
	ValidUntil int64  `json:"validUntil"`
}

type Token struct {
	ID      int64  `json:"tokenId"`
	NftData string `json:"nftData"`
	Amount  string `json:"amount"`
}

type Fee struct {
	ID     int64  `json:"tokenId"`
	Amount string `json:"amount"`
}

type TxType string

const (
	Transfer TxType = "Transfer"
)
