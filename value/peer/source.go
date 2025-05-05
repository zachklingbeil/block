package peer

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wealdtech/go-ens/v3"
)

// ENS -> hex
func (p *Peers) GetAddress(peer *Peer) {
	if peer.Address == "." || peer.Address != "" {
		return
	}
	address, err := ens.Resolve(p.Factory.Eth, peer.ENS)

	p.Factory.Rw.Lock()
	defer p.Factory.Rw.Unlock()
	if err != nil {
		peer.Address = "."
	} else {
		peer.Address = p.Format(address.Hex())
	}
}

// hex -> .eth
func (p *Peers) GetENS(peer *Peer) {
	if peer.ENS == "." || peer.ENS != "" {
		return
	}
	ensName, err := ens.ReverseResolve(p.Factory.Eth, common.HexToAddress(peer.Address))

	p.Factory.Rw.Lock()
	defer p.Factory.Rw.Unlock()
	if err != nil || ensName == "" {
		peer.ENS = "."
	} else {
		peer.ENS = p.Format(ensName)
	}
}

// hex -> LoopringENS [.loopring.eth] or "."
func (p *Peers) GetLoopringENS(peer *Peer) {
	if peer.LoopringENS == "." || (peer.LoopringENS != "" && peer.LoopringENS != "!") {
		// Proceed only if any of these conditions are true
		url := fmt.Sprintf(dotLoop, peer.Address)
		var response struct {
			Loopring string `json:"data"`
		}
		err := p.input(url, &response)

		p.Factory.Rw.Lock()
		defer p.Factory.Rw.Unlock()
		if err != nil {
			peer.LoopringENS = "!"
		} else if response.Loopring == "" {
			peer.LoopringENS = "."
		} else {
			peer.LoopringENS = p.Format(response.Loopring)
		}
	}
}

// hex -> LoopringId or "."
func (p *Peers) GetLoopringID(peer *Peer) {
	url := fmt.Sprintf(byAddress, peer.Address)
	var response struct {
		ID int64 `json:"accountId"`
	}
	p.input(url, &response)
	p.Factory.Rw.Lock()
	defer p.Factory.Rw.Unlock()
	switch {
	case response.ID == 0:
		peer.LoopringID = -1
	default:
		peer.LoopringID = response.ID
	}
}

// LoopringId -> hex
func (p *Peers) GetLoopringAddress(peer *Peer) {
	if peer.Address == "." || (peer.Address != "" && peer.Address != "!") {
		return
	}
	url := fmt.Sprintf(byId, strconv.FormatInt(peer.LoopringID, 10))
	var response struct {
		Address string `json:"owner"`
	}
	err := p.input(url, &response)
	p.Factory.Rw.Lock()
	defer p.Factory.Rw.Unlock()
	switch {
	case err != nil:
		peer.Address = "!"
	case response.Address == "":
		peer.Address = "."
	default:
		peer.Address = p.Format(response.Address)
	}
}
