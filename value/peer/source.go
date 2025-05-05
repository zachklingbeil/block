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
	p.Factory.Rw.Lock()
	defer p.Factory.Rw.Unlock()
	p.Map[peer.Address] = peer
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
	p.Factory.Rw.Lock()
	defer p.Factory.Rw.Unlock()
	p.Map[peer.ENS] = peer
}

// hex -> LoopringENS [.loopring.eth] or "."
func (p *Peers) GetLoopringENS(peer *Peer) {
	// If LoopringENS is empty or "!", set it to "."
	if peer.LoopringENS == "" || peer.LoopringENS == "!" {
		peer.LoopringENS = "."
		return
	}
	// If LoopringENS already has a valid value, do nothing
	if peer.LoopringENS != "." && peer.LoopringENS != "" && peer.LoopringENS != "!" {
		return
	}

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
	p.Factory.Rw.Lock()
	defer p.Factory.Rw.Unlock()
	p.Map[peer.LoopringENS] = peer
}

// func (p *Peers) GetLoopringENS(peer *Peer) {
// 	if peer.LoopringENS == "" || peer.LoopringENS == "." || peer.LoopringENS == "!" {
// 		// Proceed only if any of these conditions are true
// 		url := fmt.Sprintf(dotLoop, peer.Address)
// 		var response struct {
// 			Loopring string `json:"data"`
// 		}
// 		err := p.input(url, &response)

// 		p.Factory.Rw.Lock()
// 		defer p.Factory.Rw.Unlock()
// 		if err != nil {
// 			peer.LoopringENS = "!"
// 		} else if response.Loopring == "" {
// 			peer.LoopringENS = "."
// 		} else {
// 			peer.LoopringENS = p.Format(response.Loopring)
// 		}
// 	}
// }

// hex -> LoopringId or "."
func (p *Peers) GetLoopringID(peer *Peer) {
	// Proceed only if LoopringID is "", "!" or "."
	if peer.LoopringID != "" && peer.LoopringID != "!" && peer.LoopringID != "." {
		return
	}
	url := fmt.Sprintf(byAddress, peer.Address)
	var response struct {
		ID int64 `json:"accountId"`
	}
	err := p.input(url, &response)
	p.Factory.Rw.Lock()
	defer p.Factory.Rw.Unlock()
	switch {
	case err != nil:
		peer.LoopringID = "!"
	case response.ID == 0:
		peer.LoopringID = "."
	default:
		peer.LoopringID = strconv.FormatInt(response.ID, 10)
	}
	p.Factory.Rw.Lock()
	defer p.Factory.Rw.Unlock()
	p.Map[peer.LoopringID] = peer
}

// LoopringId -> hex
func (p *Peers) GetLoopringAddress(peer *Peer) {
	if peer.Address != "" && peer.Address != "!" {
		return
	}
	url := fmt.Sprintf(byId, peer.LoopringID)
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
	p.Factory.Rw.Lock()
	defer p.Factory.Rw.Unlock()
	p.Map[peer.Address] = peer
}
