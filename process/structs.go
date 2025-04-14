package process

import (
	"encoding/json"

	"github.com/zachklingbeil/factory"
)

type Process struct {
	Factory *factory.Factory
	Blocks  []Block
	Raw     []RawTx
	Txs     *Txs
}

type Block struct {
	Number       int64 `json:"blockId"`
	Size         int64 `json:"blockSize"`
	Timestamp    int64 `json:"createdAt"`
	Transactions []any `json:"transactions"`
}

type Txs struct {
	Deposit       []DW            `json:"deposit,omitempty"`
	Withdrawal    []DW            `json:"withdraw,omitempty"`
	Swaps         []Swap          `json:"swap,omitempty"`
	Transfers     []Transfer      `json:"transfer,omitempty"`
	Mints         []Mint          `json:"mint,omitempty"`
	AccountUpdate []AccountUpdate `json:"accountUpdate,omitempty"`
	AmmUpdate     []AmmUpdate     `json:"ammUpdate,omitempty"`
	NftData       []NftData       `json:"nftData,omitempty"`
	TBD           []any           `json:"tbd,omitempty"`
	*json.RawMessage
}

type Tx struct {
	Zero        int64       `json:"zero,omitempty"`
	ZeroId      string      `json:"zeroId,omitempty"`
	One         string      `json:"one,omitempty"`
	Value       string      `json:"value,omitempty"`
	Token       int64       `json:"token,omitempty"`
	Fee         string      `json:"fee,omitempty"`
	FeeToken    int64       `json:"feeToken,omitempty"`
	Coordinates *Coordinate `json:"coordinates,omitempty"`
	*json.RawMessage
}

// Depost,  Withdraw (fee)
type DW struct {
	Zero        int64       `json:"accountId"`
	ZeroAddress string      `json:"fromAddress"`
	One         string      `json:"toAddress"`
	Value       string      `json:"token.amount"`
	Token       int64       `json:"token.tokenId"`
	Fee         string      `json:"fee.amount,omitempty"`
	FeeToken    int64       `json:"fee.tokenId,omitempty"`
	Coordinates *Coordinate `json:"coordinates,omitempty"`
}

type Swap struct {
	Zero        int64       `json:"orderA.accountID"`
	ZeroValue   string      `json:"orderA.filledS"`
	ZeroToken   int64       `json:"orderB.tokenB"`
	One         int64       `json:"orderB.accountID"`
	OneValue    string      `json:"orderB.filledS"`
	OneToken    int64       `json:"orderA.tokenB"`
	ZeroFee     int64       `json:"orderA.feeBips"`
	OneFee      int64       `json:"fee.orderB.feeBips"`
	Coordinates *Coordinate `json:"coordinates,omitempty"`
}

type RawTx struct {
	Zero      string `json:"fromAddress,omitempty"`
	ZeroId    int64  `json:"accountId,omitempty"`
	ZeroSwap  int64  `json:"orderA.accountID,omitempty"`
	ZeroValue string `json:"orderA.filledS,omitempty"`
	ZeroToken int64  `json:"orderB.tokenB,omitempty"`
	ZeroFee   int64  `json:"fee.orderA.feeBips,omitempty"`

	One        string `json:"toAddress,omitempty"`
	OneId      int64  `json:"toAccountId,omitempty"`
	OneSwap    int64  `json:"orderB.accountID,omitempty"`
	OneValue   string `json:"orderB.filledS,omitempty"`
	OneToken   int64  `json:"orderA.tokenB,omitempty"`
	OneFee     int64  `json:"fee.orderB.feeBips,omitempty"`
	OneAddress string `json:"toAccountAddress,omitempty"`

	Type        string      `json:"txType,omitempty"`
	Value       string      `json:"token.amount,omitempty"`
	Token       int64       `json:"token.tokenId,omitempty"`
	Fee         string      `json:"fee.amount,omitempty"`
	FeeToken    int64       `json:"fee.tokenId,omitempty"`
	Coordinates *Coordinate `json:"coordinates,omitempty"`
	*json.RawMessage
}

type Transfer struct {
	Zero        int64       `json:"accountId"`
	One         int64       `json:"toAccountId"`
	OneAddress  string      `json:"toAccountAddress"`
	Value       string      `json:"token.amount"`
	Token       int64       `json:"token.tokenId"`
	Fee         string      `json:"fee.amount,omitempty"`
	FeeToken    int64       `json:"fee.tokenId,omitempty"`
	Coordinates *Coordinate `json:"coordinates,omitempty"`
}

type Mint struct {
	Zero        int64  `json:"minterAccountId"`
	ZeroAddress string `json:"toAccountAddress"`
	// Nft         string      `json:"toToken.tokenId"`
	NftId       string      `json:"nftToken.nftId"`
	NftData     string      `json:"nftToken.nftData"`
	NftAddress  string      `json:"nftToken.tokenAddress"`
	Quantity    string      `json:"nftToken.amount"`
	Fee         string      `json:"fee.amount,omitempty"`
	FeeToken    int64       `json:"fee.tokenId,omitempty"`
	Coordinates *Coordinate `json:"coordinates,omitempty"`
}

type AccountUpdate struct {
	Zero        int64       `json:"accountId"`
	Coordinates *Coordinate `json:"coordinates,omitempty"`
}

type AmmUpdate struct {
	Zero        int64       `json:"accountId"`
	ZeroAddress string      `json:"owner"`
	Nonce       int64       `json:"nonce"`
	Coordinates *Coordinate `json:"coordinates,omitempty"`
}

type NftData struct {
	Zero        int64       `json:"accountId"`
	One         string      `json:"minter"`
	NftId       string      `json:"nftToken.nftId"`
	NftData     string      `json:"nftToken.nftData"`
	NftAddress  string      `json:"nftToken.tokenAddress"`
	Coordinates *Coordinate `json:"coordinates,omitempty"`
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
