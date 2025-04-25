package peer

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/zachklingbeil/block/circuit"
	"github.com/zachklingbeil/factory"
)

type Peers struct {
	Factory        *factory.Factory
	Circuit        *circuit.Circuit
	Slice          []Peer
	LoopringApiKey string
}

type Peer struct {
	Address     string `json:"address"`
	ENS         string `json:"ens"`
	LoopringENS string `json:"loopringEns"`
	LoopringID  int64  `json:"loopringId"`
}

func HelloPeers(factory *factory.Factory, circuit *circuit.Circuit) *Peers {
	peers := &Peers{
		Factory:        factory,
		LoopringApiKey: os.Getenv("LOOPRING_API_KEY"),
		Circuit:        circuit,
		Slice:          make([]Peer, 240000),
	}

	if err := peers.LoadPeers(); err != nil {
		fmt.Printf("Error loading peers: %v\n", err)
	}

	total := len(peers.Slice)
	for i, peer := range peers.Slice {
		if peer.Address != "" {
			circuit.AddString(peer.ENS, peer)
			circuit.AddString(peer.LoopringENS, peer)
			circuit.AddInt(peer.LoopringID, peer)
			fmt.Printf("%d\n", total-i)
		}
	}
	return peers
}
