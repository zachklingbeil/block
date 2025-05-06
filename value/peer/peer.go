package peer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zachklingbeil/factory"
)

type Peers struct {
	Factory       *factory.Factory
	Peers         []*Peer `json:"peers,omitempty"`
	LoopringIdMap map[int64]*Peer
	Map           map[string]*Peer
}

type Peer struct {
	ENS         string `json:"ens,omitempty"`
	LoopringENS string `json:"loopringEns,omitempty"`
	LoopringID  int64  `json:"loopringId,omitempty"`
	Address     string `json:"address,omitempty"`
}

func NewPeers(factory *factory.Factory) *Peers {
	peers := &Peers{
		Factory: factory,
	}
	peers.LoadPeers()
	return peers
}

const (
	byAddress = "https://api3.loopring.io/api/v3/account?owner=%s"
	byId      = "https://api3.loopring.io/api/v3/account?accountId=%s"
	dotLoop   = "https://api3.loopring.io/api/wallet/v3/resolveName?owner=%s"
)

func (p *Peers) Format(address string) string {
	address = strings.ToLower(address)
	if strings.HasPrefix(address, "0x") || strings.HasSuffix(address, ".eth") {
		return address
	}
	return address
}

func (p *Peers) input(url string, response any) error {
	data, err := p.Factory.Json.In(url, "")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, response)
}

func (p *Peers) Refresh() {
	for i := range p.Peers {
		fmt.Printf("%d\n", i)
		peer := p.Peers[i]
		p.Format(peer.Address)
		p.Save(peer)
	}
}

// GetAddressByLoopringID returns the address for a given LoopringID, or an empty string if not found.
func (p *Peers) GetAddress(id int64) string {
	peer, ok := p.LoopringIdMap[id]
	if !ok || peer == nil {
		return ""
	}
	return strings.ToLower(peer.Address)
}
