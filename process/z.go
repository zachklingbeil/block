package process

// type Txs struct {
// 	Deposit       []DW            `json:"deposit,omitempty"`
// 	Withdrawal    []DW            `json:"withdraw,omitempty"`
// 	Swaps         []Swap          `json:"swap,omitempty"`
// 	Transfers     []Transfer      `json:"transfer,omitempty"`
// 	Mints         []Mint          `json:"mint,omitempty"`
// 	AccountUpdate []AccountUpdate `json:"accountUpdate,omitempty"`
// 	AmmUpdate     []AmmUpdate     `json:"ammUpdate,omitempty"`
// 	NftData       []NftData       `json:"nftData,omitempty"`
// 	TBD           []any           `json:"tbd,omitempty"`
// 	*json.RawMessage
// }

// type Txs struct {
// 	Deposit       []DW            `json:"deposit,omitempty"`
// 	Withdrawal    []DW            `json:"withdraw,omitempty"`
// 	Swaps         []Swap          `json:"swap,omitempty"`
// 	Transfers     []Transfer      `json:"transfer,omitempty"`
// 	Mints         []Mint          `json:"mint,omitempty"`
// 	AccountUpdate []AccountUpdate `json:"accountUpdate,omitempty"`
// 	AmmUpdate     []AmmUpdate     `json:"ammUpdate,omitempty"`
// 	NftData       []NftData       `json:"nftData,omitempty"`
// 	TBD           []any           `json:"tbd,omitempty"`
// 	*json.RawMessage
// }

// // Depost,  Withdraw (fee)
// type DW struct {
// 	Zero        int64       `json:"accountId"`
// 	ZeroAddress string      `json:"fromAddress"`
// 	One         string      `json:"toAddress"`
// 	Value       string      `json:"token.amount"`
// 	Token       int64       `json:"token.tokenId"`
// 	Fee         string      `json:"fee.amount,omitempty"`
// 	FeeToken    int64       `json:"fee.tokenId,omitempty"`
// 	Coordinates *Coordinate `json:"coordinates,omitempty"`
// }

// type Swap struct {
// 	Zero        int64       `json:"orderA.accountID"`
// 	ZeroValue   string      `json:"orderA.filledS"`
// 	ZeroToken   int64       `json:"orderB.tokenB"`
// 	One         int64       `json:"orderB.accountID"`
// 	OneValue    string      `json:"orderB.filledS"`
// 	OneToken    int64       `json:"orderA.tokenB"`
// 	ZeroFee     int64       `json:"orderA.feeBips"`
// 	OneFee      int64       `json:"fee.orderB.feeBips"`
// 	Coordinates *Coordinate `json:"coordinates,omitempty"`
// }

// type Transfer struct {
// 	Zero        int64       `json:"accountId"`
// 	One         int64       `json:"toAccountId"`
// 	OneAddress  string      `json:"toAccountAddress"`
// 	Value       string      `json:"token.amount"`
// 	Token       int64       `json:"token.tokenId"`
// 	Fee         string      `json:"fee.amount,omitempty"`
// 	FeeToken    int64       `json:"fee.tokenId,omitempty"`
// 	Coordinates *Coordinate `json:"coordinates,omitempty"`
// }

// type Mint struct {
// 	Zero        int64  `json:"minterAccountId"`
// 	ZeroAddress string `json:"toAccountAddress"`
// 	// Nft         string      `json:"toToken.tokenId"`
// 	NftId       string      `json:"nftToken.nftId"`
// 	NftData     string      `json:"nftToken.nftData"`
// 	NftAddress  string      `json:"nftToken.tokenAddress"`
// 	Quantity    string      `json:"nftToken.amount"`
// 	Fee         string      `json:"fee.amount,omitempty"`
// 	FeeToken    int64       `json:"fee.tokenId,omitempty"`
// 	Coordinates *Coordinate `json:"coordinates,omitempty"`
// }

// type AccountUpdate struct {
// 	Zero        int64       `json:"accountId"`
// 	Coordinates *Coordinate `json:"coordinates,omitempty"`
// }

// type AmmUpdate struct {
// 	Zero        int64       `json:"accountId"`
// 	ZeroAddress string      `json:"owner"`
// 	Nonce       int64       `json:"nonce"`
// 	Coordinates *Coordinate `json:"coordinates,omitempty"`
// }

// type NftData struct {
// 	Zero        int64       `json:"accountId"`
// 	One         string      `json:"minter"`
// 	NftId       string      `json:"nftToken.nftId"`
// 	NftData     string      `json:"nftToken.nftData"`
// 	NftAddress  string      `json:"nftToken.tokenAddress"`
// 	Coordinates *Coordinate `json:"coordinates,omitempty"`
// }
