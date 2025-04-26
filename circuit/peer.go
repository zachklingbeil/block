package circuit

import (
	"encoding/json"
	"log"

	"github.com/ethereum/go-ethereum/common"
)

func (c *Circuit) LoadPeers() error {
	source, err := c.Factory.Redis.SMembers(c.Factory.Ctx, "peers").Result()
	if err != nil {
		return err
	}

	for _, peerJSON := range source {
		var peer Peer
		if err := json.Unmarshal([]byte(peerJSON), &peer); err != nil {
			log.Printf("Skipping invalid peer: %v", err)
			continue
		}
		c.Peers = append(c.Peers, peer)
		c.Map[peer.Address] = &c.Peers
	}
	c.Factory.State.Add("peers", len(c.Peers))
	return nil
}

func (c *Circuit) HelloUniverse(key string) {
	c.GetPeer(key)
	c.Process(key)
}

func (c *Circuit) GetPeer(value string) *Peer {
	c.Factory.Rw.RLock()
	peer, exists := c.PeerMap[value]
	c.Factory.Rw.RUnlock()
	if exists {
		return peer
	}
	return c.CreatePeer(value)
}

func (c *Circuit) CreatePeer(value string) *Peer {
	c.Factory.Rw.Lock()
	defer c.Factory.Rw.Unlock()

	new := &Peer{}

	switch {
	case common.IsHexAddress(value):
		new.Address = value
	case len(value) > 12 && value[len(value)-13:] == ".loopring.eth":
		new.LoopringENS = value
	case len(value) > 4 && value[len(value)-4:] == ".eth":
		new.ENS = value
	default:
		new.LoopringID = value
	}
	c.PeerMap[value] = new
	c.Peers = append(c.Peers, *new)
	return new
}

func (c *Circuit) Process(address string) {
	peer := c.GetPeer(address)
	c.GetENS(peer)
	c.GetLoopringID(peer)
	c.GetLoopringENS(peer)
}
