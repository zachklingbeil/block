package loopring

import "encoding/json"

type In struct {
	Block        int64 `json:"blockId"`
	Size         int64 `json:"blockSize"`
	Timestamp    int64 `json:"createdAt"`
	Transactions []any `json:"transactions"`
}

type Tx struct {
	Zero     *Zero       `json:"zero"`
	ZeroID   *ZeroID     `json:"zeroId"`
	One      *One        `json:"one"`
	OneID    *OneID      `json:"oneId"`
	TokenIn  *Token      `json:"tokenIn"`
	TokenOut *TokenOut   `json:"tokenOut"`
	Fee      *Fee        `json:"fee"`
	FeeToken *FeeToken   `json:"feeToken"`
	Value    int64       `json:"value"`
	Coord    *Coordinate `json:"coordinates,omitempty"`
	*json.RawMessage
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

type Token struct {
	A *int64  `json:"tokenId,omitempty"`
	B *int64  `json:"orderATokenB,omitempty"`
	C *int64  `json:"orderBTokenB,omitempty"`
	D *string `json:"orderAAmountB,omitempty"`
	E *string `json:"orderBAmountB,omitempty"`
}

type TokenOut struct {
	A *int64  `json:"toTokenId,omitempty"`
	B *int64  `json:"orderATokenS,omitempty"`
	C *int64  `json:"orderBTokenS,omitempty"`
	D *string `json:"orderAAmountS,omitempty"`
	E *string `json:"orderAFilledS,omitempty"`
	F *string `json:"orderBAmountS,omitempty"`
	G *string `json:"orderBFilledS,omitempty"`
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

type Misc struct {
	Nonce           *int64  `json:"nonce,omitempty"`
	Timestamp       *int64  `json:"timestamp,omitempty"`
	Memo            *string `json:"memo,omitempty"`
	WithdrawnTo     *string `json:"withdrawalInfoRecipient,omitempty"`
	ZeroNftData     *string `json:"orderANftData,omitempty"`
	OneNftData      *string `json:"orderBNftData,omitempty"`
	NftTokenAddress *string `json:"nftTokenAddress,omitempty"`
	NftTokenID      *string `json:"nftId,omitempty"`
	NftData         *string `json:"nftData,omitempty"`
}
