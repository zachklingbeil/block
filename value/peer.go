package value

// type Peers struct {
// 	Factory *factory.Factory
// 	Peers   []One
// 	Map     map[string]*One
// }

// func NewPeers(factory *factory.Factory) *Peers {
// 	p := &Peers{
// 		Factory: factory,
// 		Peers:   make([]One, 0),
// 	}

// 	if err := p.LoadPeers(); err != nil {
// 		log.Fatalf("Failed to load peers from Redis: %v", err)
// 	}
// 	return p
// }

// func (p *Peers) GetPeer(value string) *One {
// 	p.Factory.Rw.RLock()
// 	peer, exists := p.Map[value]
// 	p.Factory.Rw.RUnlock()
// 	if exists {
// 		return peer
// 	}
// 	return p.CreatePeer(value)
// }

// func (p *Peers) CreatePeer(value string) *One {
// 	p.Factory.Rw.Lock()
// 	defer p.Factory.Rw.Unlock()

// 	new := &One{}

// 	switch {
// 	case common.IsHexAddress(value):
// 		new.Address = value
// 	case len(value) > 12 && value[len(value)-13:] == ".loopring.eth":
// 		new.LoopringENS = value
// 	case len(value) > 4 && value[len(value)-4:] == ".eth":
// 		new.ENS = value
// 	default:
// 		new.LoopringID = value
// 	}
// 	p.Map[value] = new
// 	p.Peers = append(p.Peers, *new)
// 	return new
// }

// func (p *Peers) Process(address string) {
// 	peer := p.GetPeer(address)
// 	p.GetENS(peer)
// 	p.GetLoopringID(peer)
// 	p.GetLoopringENS(peer)
// }

// func (p *Peers) LoadPeers() error {
// 	source, err := p.Factory.Data.RB.SMembers(p.Factory.Ctx, "peers").Result()
// 	if err != nil {
// 		return err
// 	}

// 	for _, peerJSON := range source {
// 		var peer One
// 		if err := json.Unmarshal([]byte(peerJSON), &peer); err != nil {
// 			log.Printf("Skipping invalid peer: %v", err)
// 			continue
// 		}
// 		p.Peers = append(p.Peers, peer)
// 	}
// 	return nil
// }
// func (p *Peers) HelloUniverse(key string) {
// 	p.GetPeer(key)
// 	p.Process(key)
// }
