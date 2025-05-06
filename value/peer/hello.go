package peer

// func (p *Peers) Hello(value string) string {
// 	p.Factory.Rw.RLock()
// 	peer, exists := p.Map[value]
// 	p.Factory.Rw.RUnlock()

// 	if !exists {
// 		peer = &Peer{}
// 		if common.IsHexAddress(value) {
// 			peer.Address = p.Format(value)
// 		} else {
// 			loopringID, err := strconv.ParseInt(value, 10, 64)
// 			if err != nil {
// 				fmt.Printf("Error converting value to int64: %v\n", err)
// 				return ""
// 			}
// 			peer.LoopringID = loopringID
// 			p.GetLoopringAddress(peer)
// 		}
// 		p.Peers = append(p.Peers, peer)
// 	}
// 	return peer.Address
// }

// func (p *Peers) HelloUniverse(peer *Peer) *Peer {
// 	p.GetENS(peer)
// 	p.GetLoopringENS(peer)
// 	p.GetLoopringID(peer)

// 	p.Factory.Rw.Lock()
// 	p.Map[peer.Address] = peer
// 	p.Save(peer)
// 	p.Factory.Rw.Unlock()
// 	fmt.Printf("	%s %s %s %d\n", peer.Address, peer.ENS, peer.LoopringENS, peer.LoopringID)
// 	return peer
// }
