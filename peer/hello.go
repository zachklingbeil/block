package peer

// func (p *Peers) Hello(value string) string {
// 	p.Factory.Rw.RLock()
// 	peer, exists := p.Map[value]
// 	p.Factory.Rw.RUnlock()

// 	if !exists {
// 		return value
// 	}
// 	// Prefer ENS, then LoopringENS, then Address
// 	switch {
// 	case peer.ENS != "" && peer.ENS != "." && peer.ENS != "!":
// 		return peer.ENS
// 	case peer.LoopringENS != "" && peer.LoopringENS != "." && peer.LoopringENS != "!":
// 		return peer.LoopringENS
// 	default:
// 		return peer.Address
// 	}
// }

// // func (p *Peers) Hello(value string) string {
// // 	p.Factory.Rw.RLock()
// // 	peer, exists := p.Map[value]
// // 	p.Factory.Rw.RUnlock()

// // 	if !exists {
// // 		peer = &Peer{}
// // 		if common.IsHexAddress(value) {
// // 			peer.Address = p.Format(value)
// // 		} else {
// // 			peer.LoopringID = value
// // 		}
// // 		p.Factory.Rw.Lock()
// // 		p.Peers = append(p.Peers, peer)
// // 		p.Map[value] = peer
// // 		p.Factory.Rw.Unlock()
// // 		p.HelloUniverse(peer)
// // 	}

// // 	// Prefer ENS, then LoopringENS, then Address
// // 	switch {
// // 	case peer.ENS != "" && peer.ENS != "." && peer.ENS != "!":
// // 		return peer.ENS
// // 	case peer.LoopringENS != "" && peer.LoopringENS != "." && peer.LoopringENS != "!":
// // 		return peer.LoopringENS
// // 	default:
// // 		return peer.Address
// // 	}
// // }

// func (p *Peers) HelloUniverse(peer *Peer) *Peer {
// 	p.GetENS(peer)
// 	// p.GetLoopringENS(peer) // Ensure this is commented out if not needed
// 	p.GetLoopringID(peer)

// 	p.Factory.Rw.Lock()
// 	if common.IsHexAddress(peer.Address) {
// 		p.Map[peer.Address] = peer
// 	}
// 	if peer.ENS != "" && peer.ENS != "." && peer.ENS != "!" {
// 		p.Map[peer.ENS] = peer
// 	}
// 	if peer.LoopringENS != "" && peer.LoopringENS != "." && peer.LoopringENS != "!" {
// 		p.Map[peer.LoopringENS] = peer
// 	}
// 	if peer.LoopringID != "" && peer.LoopringID != "." && peer.LoopringID != "!" {
// 		p.Map[peer.LoopringID] = peer
// 	}
// 	p.Save(peer)
// 	p.Factory.Rw.Unlock()

// 	// Print all peer details, including manually assigned LoopringENS
// 	fmt.Printf("	%s %s %s %s\n", peer.Address, peer.ENS, peer.LoopringENS, peer.LoopringID)
// 	return peer
// }
