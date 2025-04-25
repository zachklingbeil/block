package circuit

type Raw struct {
	Number       int64 `json:"blockId"`
	Timestamp    int64 `json:"createdAt"`
	Size         int64 `json:"blockSize"`
	Transactions []any `json:"transactions"`
}
type Token struct {
	Symbol    string `json:"symbol,omitempty"`
	Address   string `json:"address,omitempty"`
	TokenId   int64  `json:"tokenId,omitempty"`
	Decimals  int    `json:"decimals,omitempty"`
	AccountID int64  `json:"accountId"`
}
