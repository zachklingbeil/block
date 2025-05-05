package peer

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

func (p *Peers) Hello(value string) string {
	p.Factory.Rw.RLock()
	peer, exists := p.Map[value]
	p.Factory.Rw.RUnlock()

	if !exists {
		peer = &Peer{}
		if common.IsHexAddress(value) {
			peer.Address = p.Format(value)
		} else {
			peer.LoopringID = value
			p.GetLoopringAddress(peer)
		}
		p.Factory.Rw.RLock()
		p.Peers = append(p.Peers, peer)
		p.Factory.Rw.RUnlock()
		p.HelloUniverse(peer)
	}

	// // Prefer ENS, then LoopringENS, then Address
	// switch {
	// case peer.ENS != "" && peer.ENS != "." && peer.ENS != "!":
	// 	return peer.ENS
	// case peer.LoopringENS != "" && peer.LoopringENS != "." && peer.LoopringENS != "!":
	// 	return peer.LoopringENS
	// default:
	// 	return peer.Address
	// }
	return peer.Address
}

func (p *Peers) HelloUniverse(peer *Peer) *Peer {
	p.GetENS(peer)
	p.GetLoopringENS(peer)
	p.GetLoopringID(peer)

	p.Factory.Rw.Lock()
	p.Map[peer.Address] = peer
	p.Map[peer.ENS] = peer
	p.Map[peer.LoopringENS] = peer
	p.Map[peer.LoopringID] = peer
	p.Save(peer)
	p.Factory.Rw.Unlock()
	fmt.Printf("	%s %s %s %s\n", peer.Address, peer.ENS, peer.LoopringENS, peer.LoopringID)
	return peer
}
