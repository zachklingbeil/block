package loopring

type Main struct {
	BlockID      int64 `json:"blockId"`
	BlockSize    int64 `json:"blockSize"`
	CreatedAt    int64 `json:"createdAt"`
	Transactions []any `json:"transactions"`
	// Exchange     string `json:"exchange"`
	// TxHash       string `json:"txHash"`
	// Status       string `json:"status"`
}

type Txs struct {
	// AccountID       *int64 `json:"accountId,omitempty"`
	// MinterAccountID *int64 `json:"minterAccountId,omitempty"`
	// OrderAAccountID *int64 `json:"orderAAccountID,omitempty"`

	// Owner         *string `json:"owner,omitempty"`
	// FromAddress   *string `json:"fromAddress,omitempty"`
	// SenderAddress *string `json:"senderAddress,omitempty"`

	// ToAccountID      *int64  `json:"toAccountId,omitempty"`
	// OrderBAccountID  *int64  `json:"orderBAccountID,omitempty"`
	// ToAddress        *string `json:"toAddress,omitempty"`
	// ToAccountAddress *string `json:"toAccountAddress,omitempty"`
	// ReceiverAddress  *string `json:"receiverAddress,omitempty"`

	// TokenID   *int64 `json:"tokenId,omitempty"`
	// ToTokenID *int64 `json:"toTokenId,omitempty"`

	// FeeTokenID     *int64  `json:"feeTokenId,omitempty"`
	// FeeTokenSymbol *string `json:"feeTokenSymbol,omitempty"`

	// Value     *string `json:"amount,omitempty"`
	// FeeAmount *string `json:"feeAmount,omitempty"`

	// OrderAFeeBips *int64 `json:"orderAFeeBips,omitempty"`
	// OrderBFeeBips *int64 `json:"orderBFeeBips,omitempty"`

	// OrderATokenB *int64 `json:"orderATokenB,omitempty"`
	// OrderATokenS *int64 `json:"orderATokenS,omitempty"`
	// OrderBTokenB *int64 `json:"orderBTokenB,omitempty"`
	// OrderBTokenS *int64 `json:"orderBTokenS,omitempty"`

	// OrderAAmountB *string `json:"orderAAmountB,omitempty"`
	// OrderAAmountS *string `json:"orderAAmountS,omitempty"`
	// OrderBAmountB *string `json:"orderBAmountB,omitempty"`
	// OrderBAmountS *string `json:"orderBAmountS,omitempty"`

	// OrderAFillS *int64 `json:"orderAFillS,omitempty"` /////////
	// OrderBFillS *int64 `json:"orderBFillS,omitempty"` ////////

	// OrderATaker *string `json:"orderATaker,omitempty"`
	// OrderBTaker *string `json:"orderBTaker,omitempty"`

	// OrderAFilledS *string `json:"orderAFilledS,omitempty"`
	// OrderBFilledS *string `json:"orderBFilledS,omitempty"`

	// OrderANftData *string `json:"orderANftData,omitempty"`
	// OrderBNftData *string `json:"orderBNftData,omitempty"`

	// NftTokenAddress *string `json:"nftTokenAddress,omitempty"`
	// NftTokenID      *string `json:"nftId,omitempty"`
	// NftData         *string `json:"nftData,omitempty"`

	// Nonce     *int64  `json:"nonce,omitempty"`
	// Timestamp *int64  `json:"timestamp,omitempty"`
	// Memo      *string `json:"memo,omitempty"`

	WithdrawnTo *string `json:"withdrawalInfoRecipient,omitempty"`
}

type Tx struct {
	Zero      *Zero     `json:"zero"`
	ZeroID    *ZeroID   `json:"zeroId"`
	One       *One      `json:"one"`
	OneID     *OneID    `json:"oneId"`
	TokenIn   *TokenIn  `json:"tokenIn"`
	TokenOut  *TokenOut `json:"tokenOut"`
	Fee       *Fee      `json:"fee"`
	FeeToken  *FeeToken `json:"feeToken"`
	Value     int64     `json:"value"`
	Timestamp *int64    `json:"timestamp,omitempty"`
	Memo      *string   `json:"memo,omitempty"`
	Nonce     *int64    `json:"nonce,omitempty"`
}

type Zero struct {
	A *string `json:"owner,omitempty"`
	B *string `json:"fromAddress,omitempty"`
	C *string `json:"senderAddress,omitempty"`
}

type ZeroID struct {
	A *int64 `json:"accountId"`
	B *int64 `json:"minterAccountId"`
	C *int64 `json:"orderAAccountID"`
}

type One struct {
	A *string `json:"toAddress,omitempty"`
	B *string `json:"toAccountAddress,omitempty"`
	C *string `json:"receiverAddress,omitempty"`
}

type OneID struct {
	A *int64 `json:"toAccountId,omitempty"`
	B *int64 `json:"orderBAccountID,omitempty"`
}

type TokenIn struct {
	TokenID       *int64  `json:"tokenId,omitempty"`
	OrderATokenB  *int64  `json:"orderATokenB,omitempty"`
	OrderBTokenB  *int64  `json:"orderBTokenB,omitempty"`
	OrderAAmountB *string `json:"orderAAmountB,omitempty"`
	OrderBAmountB *string `json:"orderBAmountB,omitempty"`
}

type TokenOut struct {
	ToTokenID    *int64 `json:"toTokenId,omitempty"`
	OrderATokenS *int64 `json:"orderATokenS,omitempty"`
	OrderBTokenS *int64 `json:"orderBTokenS,omitempty"`

	OrderAAmountS *string `json:"orderAAmountS,omitempty"`
	OrderAFilledS *string `json:"orderAFilledS,omitempty"`

	OrderBAmountS *string `json:"orderBAmountS,omitempty"`
	OrderBFilledS *string `json:"orderBFilledS,omitempty"`
}

type Fee struct {
	Value *string `json:"feeAmount,omitempty"`
	BipsA *int64  `json:"orderAFeeBips,omitempty"`
	BipsB *int64  `json:"orderBFeeBips,omitempty"`
}

type FeeToken struct {
	Token  *int64  `json:"feeTokenId,omitempty"`
	Symbol *string `json:"feeTokenSymbol,omitempty"`
}

type Nft struct {
	OrderANftData *string `json:"orderANftData,omitempty"`
	OrderBNftData *string `json:"orderBNftData,omitempty"`

	NftTokenAddress *string `json:"nftTokenAddress,omitempty"`
	NftTokenID      *string `json:"nftId,omitempty"`
	NftData         *string `json:"nftData,omitempty"`
}

type Misc struct {
	Nonce     *int64  `json:"nonce,omitempty"`
	Timestamp *int64  `json:"timestamp,omitempty"`
	Memo      *string `json:"memo,omitempty"`
}
