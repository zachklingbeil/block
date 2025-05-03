package universe

type Key interface {
	GetAddress() string
	GetLoopringID() string
	GetENS() string
	GetLoopringENS() string
	GetToken() string
	GetDecimals() string
	GetTokenId() string
	GetTokenInt() int64
	GetFirstBlock() string
}

type Value struct {
	ENS         string `json:"ens,omitempty"`
	LoopringENS string `json:"loopringEns,omitempty"`
	LoopringID  string `json:"loopringId,omitempty"`
	Address     string `json:"address,omitempty"`
	FirstBlock  string `json:"firstBlock,omitempty"`
	Token       string `json:"token,omitempty"`
	Decimals    string `json:"decimals,omitempty"`
	TokenId     string `json:"tokenId,omitempty"`
	TokenInt    int64  `json:"tokenInt,omitempty"`
}

// Implement the Value interface for Peer
func (v Value) GetAddress() string {
	return v.Address
}

func (v Value) GetLoopringID() string {
	return v.LoopringID
}

func (v Value) GetENS() string {
	return v.ENS
}

func (v Value) GetLoopringENS() string {
	return v.LoopringENS
}

func (v Value) GetToken() string {
	return v.Token
}

func (v Value) GetDecimals() string {
	return v.Decimals
}
func (v Value) GetTokenId() string {
	return v.TokenId
}
func (v Value) GetTokenInt() int64 {
	return v.TokenInt
}
func (v Value) GetFirstBlock() string {
	return v.FirstBlock
}
