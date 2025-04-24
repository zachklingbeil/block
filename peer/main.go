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
	Map            map[any]*Peer
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
		Map:            make(map[any]*Peer),
		LoopringApiKey: os.Getenv("LOOPRING_API_KEY"),
		Circuit:        circuit,
	}

	if err := peers.LoadPeer(); err != nil {
		fmt.Printf("Error loading peers: %v\n", err)
	}
	return peers
}
