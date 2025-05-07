package value

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

func (v *Value) LoadPeers() error {
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	source, err := v.Factory.Data.RB.SMembers(v.Factory.Ctx, "peer").Result()
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
		v.Universe[peer.Address] = peer
		v.Maps.LoopringId[peer.LoopringID] = peer.Address
	}
	v.Peers = peers
	fmt.Printf("%d peers loaded\n", len(v.Peers))
	return nil
}

func (v *Value) Save(peer *Peer) error {
	peerJSON, err := json.Marshal(peer)
	if err != nil {
		return fmt.Errorf("failed to serialize peer: %v", err)
	}

	if err := v.Factory.Data.RB.HSet(v.Factory.Ctx, "peer", strings.ToLower(peer.Address), peerJSON).Err(); err != nil {
		return fmt.Errorf("failed to store peer in Redis hash: %v", err)
	}
	return nil
}

// hex -> .eth
func (v *Value) GetENS(peer *Peer) {
	if peer.ENS != "" && peer.ENS != "." {
		return
	}
	ensName, err := ens.ReverseResolve(v.Factory.Eth, common.HexToAddress(peer.Address))
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()
	if err != nil || ensName == "" {
		peer.ENS = "."
	} else {
		peer.ENS = v.FormatPeer(ensName)
	}
}

// hex -> LoopringENS [.loopring.eth] or "."
func (v *Value) GetLoopringENS(peer *Peer) {
	if peer.LoopringENS == "." || peer.LoopringENS != "" && peer.LoopringENS != "!" {
		url := fmt.Sprintf(dotLoop, peer.Address)
		var resp struct {
			Loopring string `json:"data"`
		}
		err := v.input(url, &resp)
		v.Factory.Rw.Lock()
		defer v.Factory.Rw.Unlock()
		switch {
		case err != nil:
			peer.LoopringENS = "!"
		case resp.Loopring == "":
			peer.LoopringENS = "."
		default:
			peer.LoopringENS = v.FormatPeer(resp.Loopring)
		}
	}
}

// hex -> LoopringId or "."
func (v *Value) GetLoopringID(peer *Peer) {
	url := fmt.Sprintf(byAddress, peer.Address)
	var resp struct {
		ID int64 `json:"accountId"`
	}
	_ = v.input(url, &resp)
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()
	if resp.ID == 0 {
		peer.LoopringID = -1
	} else {
		peer.LoopringID = resp.ID
	}
}

func (v *Value) FormatPeer(address string) string {
	address = strings.ToLower(address)
	if strings.HasPrefix(address, "0x") || strings.HasSuffix(address, ".eth") {
		return address
	}
	return address
}

func (v *Value) input(url string, response any) error {
	data, err := v.Factory.Json.In(url, "")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, response)
}

// GetAddressByLoopringID returns the address for a given LoopringID, or an empty string if not found.
func (v *Value) GetPeer(id int64) string {
	peer, ok := v.Maps.LoopringId[id]
	if !ok {
		return ""
	}
	return strings.ToLower(peer)
}

// func (v *Value) HelloUniverse(value string) *Peer {
// 	v.Factory.Rw.Lock()
// 	defer v.Factory.Rw.Unlock()

// 	peer, exists := v.Maps[value]
// 	if !exists {
// 		peer = &Peer{}
// 		if common.IsHexAddress(value) {
// 			peer.Address = v.FormatPeer(value)
// 		} else {
// 			loopringID, err := strconv.ParseInt(value, 10, 64)
// 			if err != nil {
// 				fmt.Printf("Error converting value to int64: %v\n", err)
// 				return nil
// 			}
// 			peer.LoopringID = loopringID
// 		}
// 		v.Peers = append(v.Peers, peer)
// 	}
// 	v.GetENS(peer)
// 	v.GetLoopringENS(peer)
// 	v.GetLoopringID(peer)

// 	v.Maps[peer.Address] = peer
// 	v.Save(peer)
// 	fmt.Printf("	%s %s %s %d\n", peer.Address, peer.ENS, peer.LoopringENS, peer.LoopringID)
// 	return peer
// }
