package circuit

type Constant struct {
	TokenMap map[int64]*Token
}
type Token struct {
	Symbol     string `json:"symbol"`
	Address    string `json:"address"`
	LoopringID string `json:"accountId,omitempty"`
	TokenId    int64  `json:"tokenId"`
	Decimals   int    `json:"decimals"`
}

type Peer struct {
	ENS         string `json:"ens"`
	LoopringENS string `json:"loopringEns"`
	LoopringID  string `json:"loopringId"`
	Address     string `json:"address"`
}
