package loopring

type Main struct {
	BlockID      int64         `json:"blockId"`
	BlockSize    int64         `json:"blockSize"`
	Exchange     string        `json:"exchange"`
	TxHash       string        `json:"txHash"`
	Status       string        `json:"status"`
	CreatedAt    int64         `json:"createdAt"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	TxType           TxType  `json:"txType"`
	AccountID        int64   `json:"accountId"`
	Token            Token   `json:"token"`
	ToToken          ToToken `json:"toToken"`
	Fee              Fee     `json:"fee"`
	ValidUntil       int64   `json:"validUntil"`
	ToAccountID      int64   `json:"toAccountId"`
	ToAccountAddress string  `json:"toAccountAddress"`
	StorageID        int64   `json:"storageId"`
}

type Fee struct {
	TokenID int64  `json:"tokenId"`
	Amount  string `json:"amount"`
}

type ToToken struct {
}

type Token struct {
	TokenID int64  `json:"tokenId"`
	NftData string `json:"nftData"`
	Amount  string `json:"amount"`
}

type TxType string

const (
	Transfer TxType = "Transfer"
)
