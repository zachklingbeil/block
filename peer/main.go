package peer

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/zachklingbeil/factory"
)

type Peers struct {
	Factory        *factory.Factory
	Map            map[string]*Peer
	LoopringApiKey string
}

type Peer struct {
	Address     string `json:"address"`
	ENS         string `json:"ens"`
	LoopringENS string `json:"loopringEns"`
	LoopringID  int64  `json:"loopringId"`
}

func HelloPeers(factory *factory.Factory) *Peers {
	peers := &Peers{
		Factory:        factory,
		Map:            make(map[string]*Peer),
		LoopringApiKey: os.Getenv("LOOPRING_API_KEY"),
	}

	if err := peers.LoadPeers(); err != nil {
		fmt.Printf("Error loading peers: %v\n", err)
	}

	return peers
}

func (p *Peers) NewBlock(addresses []string) {
	p.Factory.Mu.Lock()
	defer p.Factory.Mu.Unlock()

	for _, address := range addresses {
		if _, exists := p.Map[address]; !exists {
			p.Map[address] = &Peer{Address: address}
			fmt.Printf("Added peer: %s\n", address)
		}
	}
	p.Factory.When.Signal()
}

// func (p *Peers) HelloUniverse() {
// 	const batchSize = 1000
// 	var batch []*Peer

// 	for {
// 		p.Factory.Mu.Lock()

// 		if len(p.Map) == 0 {
// 			p.saveBatch(&batch)
// 			fmt.Println("Hello Universe")
// 			p.Factory.When.Wait()
// 		}

// 		for _, peer := range p.Map {
// 			p.Factory.Mu.Unlock()
// 			p.processPeer(peer)
// 			batch = append(batch, peer)
// 			fmt.Printf("%d %s %s %d\n", len(batch), peer.ENS, peer.LoopringENS, peer.LoopringID)

// 		}
// 	}
// }

// func (p *Peers) processPeer(peer *Peer) {
// 	p.GetENS(peer, peer.Address)
// 	p.GetLoopringENS(peer, peer.Address)
// 	p.GetLoopringID(peer, peer.Address)
// }
