package loopring

type Main struct {
	BlockID      int64  `json:"blockId"`
	BlockSize    int64  `json:"blockSize"`
	Exchange     string `json:"exchange"`
	TxHash       string `json:"txHash"`
	Status       string `json:"status"`
	CreatedAt    int64  `json:"createdAt"`
	Transactions []any  `json:"transactions"`
}

type Txs struct {
	AccountID       *int64  `json:"accountId,omitempty"`
	MinterAccountID *int64  `json:"minterAccountId,omitempty"`
	OrderAAccountID *int64  `json:"orderAAccountID,omitempty"`
	Owner           *string `json:"owner,omitempty"`
	FromAddress     *string `json:"fromAddress,omitempty"`
	SenderAddress   *string `json:"senderAddress,omitempty"`

	ToAccountID      *int64  `json:"toAccountId,omitempty"`
	OrderBAccountID  *int64  `json:"orderBAccountID,omitempty"`
	ToAddress        *string `json:"toAddress,omitempty"`
	ToAccountAddress *string `json:"toAccountAddress,omitempty"`
	ReceiverAddress  *string `json:"receiverAddress,omitempty"`

	TokenID        *int64  `json:"tokenId,omitempty"`
	ToTokenID      *int64  `json:"toTokenId,omitempty"`
	FeeTokenID     *int64  `json:"feeTokenId,omitempty"`
	FeeTokenSymbol *string `json:"feeTokenSymbol,omitempty"`

	Value     *string `json:"amount,omitempty"`
	FeeAmount *string `json:"feeAmount,omitempty"`

	OrderAFeeBips *int64 `json:"orderAFeeBips,omitempty"`
	OrderBFeeBips *int64 `json:"orderBFeeBips,omitempty"`

	OrderATokenB *int64 `json:"orderATokenB,omitempty"`
	OrderATokenS *int64 `json:"orderATokenS,omitempty"`
	OrderBTokenB *int64 `json:"orderBTokenB,omitempty"`
	OrderBTokenS *int64 `json:"orderBTokenS,omitempty"`

	OrderAAmountB *string `json:"orderAAmountB,omitempty"`
	OrderAAmountS *string `json:"orderAAmountS,omitempty"`
	OrderBAmountB *string `json:"orderBAmountB,omitempty"`
	OrderBAmountS *string `json:"orderBAmountS,omitempty"`

	OrderAFillS   *int64  `json:"orderAFillS,omitempty"`
	OrderAFilledS *string `json:"orderAFilledS,omitempty"`
	OrderBFillS   *int64  `json:"orderBFillS,omitempty"`
	OrderBFilledS *string `json:"orderBFilledS,omitempty"`

	OrderATaker *string `json:"orderATaker,omitempty"`
	OrderBTaker *string `json:"orderBTaker,omitempty"`

	OrderANftData *string `json:"orderANftData,omitempty"`
	OrderBNftData *string `json:"orderBNftData,omitempty"`

	NftTokenAddress *string `json:"nftTokenAddress,omitempty"`
	NftTokenID      *string `json:"nftId,omitempty"`
	NftData         *string `json:"nftData,omitempty"`

	Nonce     *int64  `json:"nonce,omitempty"`
	Timestamp *int64  `json:"timestamp,omitempty"`
	Memo      *string `json:"memo,omitempty"`

	WithdrawalInfo *struct {
		Recipient *string `json:"recipient,omitempty"`
	} `json:"withdrawalInfo,omitempty"`
}

type Tx struct {
	Zero     string `json:"zero"`
	One      string `json:"one"`
	Token    string `json:"token"`
	Value    int64  `json:"value"`
	Fee      int64  `json:"fee"`
	FeeToken string `json:"feeToken"`
}
