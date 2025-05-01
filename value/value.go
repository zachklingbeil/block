package value

import (
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Value struct {
	Factory  *factory.Factory
	Peers    []*Peer
	Tokens   []*Token
	Map      map[string]*Peer
	TokenMap map[any]*Token
}

type Peer struct {
	ENS         string `json:"ens"`
	LoopringENS string `json:"loopringEns"`
	LoopringID  string `json:"loopringId"`
	Address     string `json:"address"`
}

type Token struct {
	Token    string `json:"token,omitempty"`
	Address  string `json:"address,omitempty"`
	Decimals string `json:"decimals,omitempty"`
	TokenId  string `json:"tokenId,omitempty"`
	TokenInt int64  `json:"tokenInt,omitempty"`
}

func NewValue(factory *factory.Factory) *Value {
	v := &Value{
		Factory:  factory,
		Map:      make(map[string]*Peer),
		TokenMap: make(map[any]*Token),
	}

	v.LoadTokens()
	v.LoadPeers()
	v.rebuildMap()
	v.DotLoop()
	return v
}

func (v *Value) DotLoop() {
	toProcess := len(v.Peers)
	for _, peer := range v.Peers {
		peer.ENS = ""

		// Fetch ENS for all peers
		v.GetENS(peer)

		// Print peer details, decrementing the count each time
		toProcess--
		fmt.Printf("%d %s %s %s\n", toProcess, peer.ENS, peer.LoopringENS, peer.LoopringID)
	}
}
