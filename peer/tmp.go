package peer

// func (p *Peers) HelloUniverse(address any) *Peer {
// 	p.Factory.Rw.RLock()
// 	var peer *Peer
// 	var ok bool

// 	switch addr := address.(type) {
// 	case string:
// 		peer, ok = p.Map[addr].(*Peer)
// 	case int64:
// 		peer, ok = p.Map[addr].(*Peer)
// 	default:
// 		p.Factory.Rw.RUnlock()
// 		fmt.Printf("Unsupported address type: %T\n", address)
// 		return nil
// 	}
// 	p.Factory.Rw.RUnlock()

// 	if !ok || peer == nil {
// 		fmt.Printf("Peer not found for address: %v\n", address)
// 		return nil
// 	}

// 	p.GetENS(peer, peer.Address)
// 	p.GetLoopringENS(peer, peer.Address)
// 	p.GetLoopringID(peer, peer.Address)

// 	fmt.Printf("%s %s %d\n", peer.ENS, peer.LoopringENS, peer.LoopringID)
// 	return peer
// }
