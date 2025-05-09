package universe

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wealdtech/go-ens/v3"
)

type Peer struct {
	ENS         string `json:"ens,omitempty"`
	LoopringENS string `json:"loopringEns,omitempty"`
	LoopringID  int64  `json:"loopringId,omitempty"`
	Address     string `json:"address,omitempty"`
}

const (
	byAddress = "https://api3.loopring.io/api/v3/account?owner=%s"
	byId      = "https://api3.loopring.io/api/v3/account?accountId=%s"
	dotLoop   = "https://api3.loopring.io/api/wallet/v3/resolveName?owner=%s"
)

func (o *One) LoadPeers() error {
	o.Factory.Rw.Lock()
	defer o.Factory.Rw.Unlock()

	source, err := o.Factory.Data.RB.SMembers(o.Factory.Ctx, "peer").Result()
	if err != nil {
		return fmt.Errorf("failed to fetch peers from Redis hash: %v", err)
	}
	peers := make([]*Peer, 0, len(source))
	for _, peerJSON := range source {
		var peer *Peer
		if err := json.Unmarshal([]byte(peerJSON), &peer); err != nil {
			log.Printf("Skipping invalid peer: %v (data: %s)", err, peerJSON)
			continue
		}
		peers = append(peers, peer)
		o.Map[peer.Address] = peer
		o.Maps.LoopringId[peer.LoopringID] = peer.Address
	}
	o.Peers = peers
	fmt.Printf("%d peers loaded\n", len(o.Peers))
	return nil
}

func (o *One) Save(peer *Peer) error {
	peerJSON, err := json.Marshal(peer)
	if err != nil {
		return fmt.Errorf("failed to serialize peer: %v", err)
	}

	if err := o.Factory.Data.RB.HSet(o.Factory.Ctx, "peer", strings.ToLower(peer.Address), peerJSON).Err(); err != nil {
		return fmt.Errorf("failed to store peer in Redis hash: %v", err)
	}
	return nil
}

// hex -> .eth
func (o *One) GetENS(peer *Peer) {
	if peer.ENS != "" && peer.ENS != "." {
		return
	}
	ensName, err := ens.ReverseResolve(o.Factory.Eth, common.HexToAddress(peer.Address))
	o.Factory.Rw.Lock()
	defer o.Factory.Rw.Unlock()
	if err != nil || ensName == "" {
		peer.ENS = "."
	} else {
		peer.ENS = o.FormatPeer(ensName)
	}
}

// hex -> LoopringENS [.loopring.eth] or "."
func (o *One) GetLoopringENS(peer *Peer) {
	if peer.LoopringENS == "." || peer.LoopringENS != "" && peer.LoopringENS != "!" {
		url := fmt.Sprintf(dotLoop, peer.Address)
		var resp struct {
			Loopring string `json:"data"`
		}
		err := o.input(url, &resp)
		o.Factory.Rw.Lock()
		defer o.Factory.Rw.Unlock()
		switch {
		case err != nil:
			peer.LoopringENS = "!"
		case resp.Loopring == "":
			peer.LoopringENS = "."
		default:
			peer.LoopringENS = o.FormatPeer(resp.Loopring)
		}
	}
}

// hex -> LoopringId or "."
func (o *One) GetLoopringID(peer *Peer) {
	url := fmt.Sprintf(byAddress, peer.Address)
	var resp struct {
		ID int64 `json:"accountId"`
	}
	_ = o.input(url, &resp)
	o.Factory.Rw.Lock()
	defer o.Factory.Rw.Unlock()
	if resp.ID == 0 {
		peer.LoopringID = -1
	} else {
		peer.LoopringID = resp.ID
	}
}

func (o *One) FormatPeer(address string) string {
	address = strings.ToLower(address)
	if strings.HasPrefix(address, "0x") || strings.HasSuffix(address, ".eth") {
		return address
	}
	return address
}

func (o *One) input(url string, response any) error {
	data, err := o.Factory.Json.In(url, "")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, response)
}

// GetAddressByLoopringID returns the address for a given LoopringID, or an empty string if not found.
func (o *One) GetPeer(id int64) string {
	peer, ok := o.Maps.LoopringId[id]
	if !ok {
		return ""
	}
	return strings.ToLower(peer)
}

// func (o *One) HelloUniverse(value string) *Peer {
// 	o.Factory.Rw.Lock()
// 	defer o.Factory.Rw.Unlock()

// 	peer, exists := o.Maps[value]
// 	if !exists {
// 		peer = &Peer{}
// 		if common.IsHexAddress(value) {
// 			peer.Address = o.FormatPeer(value)
// 		} else {
// 			loopringID, err := strcono.ParseInt(value, 10, 64)
// 			if err != nil {
// 				fmt.Printf("Error converting value to int64: %v\n", err)
// 				return nil
// 			}
// 			peer.LoopringID = loopringID
// 		}
// 		o.Peers = append(o.Peers, peer)
// 	}
// 	o.GetENS(peer)
// 	o.GetLoopringENS(peer)
// 	o.GetLoopringID(peer)

// 	o.Maps[peer.Address] = peer
// 	o.Save(peer)
// 	fmt.Printf("	%s %s %s %d\n", peer.Address, peer.ENS, peer.LoopringENS, peer.LoopringID)
// 	return peer
// }
