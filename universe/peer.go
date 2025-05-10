package universe

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wealdtech/go-ens/v3"
)

const (
	byAddress = "https://api3.loopring.io/api/v3/account?owner=%s"
	byId      = "https://api3.loopring.io/api/v3/account?accountId=%s"
	dotLoop   = "https://api3.loopring.io/api/wallet/v3/resolveName?owner=%s"
)

func (z *Zero) LoadPeers() error {
	z.Factory.Rw.Lock()
	defer z.Factory.Rw.Unlock()

	source, err := z.Factory.Data.RB.SMembers(z.Factory.Ctx, "peer").Result()
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
		z.Map[peer.Address] = peer
		z.Maps.LoopringId[peer.LoopringID] = peer.Address
	}
	z.Peers = peers
	fmt.Printf("%d peers loaded\n", len(z.Peers))
	return nil
}

func (z *Zero) Save(peer *Peer) error {
	peerJSON, err := json.Marshal(peer)
	if err != nil {
		return fmt.Errorf("failed to serialize peer: %v", err)
	}

	if err := z.Factory.Data.RB.HSet(z.Factory.Ctx, "peer", strings.ToLower(peer.Address), peerJSON).Err(); err != nil {
		return fmt.Errorf("failed to store peer in Redis hash: %v", err)
	}
	return nil
}

// hex -> .eth
func (z *Zero) GetENS(peer *Peer) {
	if peer.ENS != "" && peer.ENS != "." {
		return
	}
	ensName, err := ens.ReverseResolve(z.Factory.Eth, common.HexToAddress(peer.Address))
	z.Factory.Rw.Lock()
	defer z.Factory.Rw.Unlock()
	if err != nil || ensName == "" {
		peer.ENS = "."
	} else {
		peer.ENS = z.FormatPeer(ensName)
	}
}

// hex -> LoopringENS [.loopring.eth] or "."
func (z *Zero) GetLoopringENS(peer *Peer) {
	if peer.LoopringENS == "." || peer.LoopringENS != "" && peer.LoopringENS != "!" {
		url := fmt.Sprintf(dotLoop, peer.Address)
		var resp struct {
			Loopring string `json:"data"`
		}
		err := z.input(url, &resp)
		z.Factory.Rw.Lock()
		defer z.Factory.Rw.Unlock()
		switch {
		case err != nil:
			peer.LoopringENS = "!"
		case resp.Loopring == "":
			peer.LoopringENS = "."
		default:
			peer.LoopringENS = z.FormatPeer(resp.Loopring)
		}
	}
}

// hex -> LoopringId or "."
func (z *Zero) GetLoopringID(peer *Peer) {
	url := fmt.Sprintf(byAddress, peer.Address)
	var resp struct {
		ID int64 `json:"accountId"`
	}
	_ = z.input(url, &resp)
	z.Factory.Rw.Lock()
	defer z.Factory.Rw.Unlock()
	if resp.ID == 0 {
		peer.LoopringID = -1
	} else {
		peer.LoopringID = resp.ID
	}
}

func (z *Zero) FormatPeer(address string) string {
	address = strings.ToLower(address)
	if strings.HasPrefix(address, "0x") || strings.HasSuffix(address, ".eth") {
		return address
	}
	return address
}

func (z *Zero) input(url string, response any) error {
	data, err := z.Factory.Json.In(url, "")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, response)
}

// GetAddressByLoopringID returns the address for a given LoopringID, or an empty string if not found.
func (z *Zero) GetPeer(id int64) string {
	peer, ok := z.Maps.LoopringId[id]
	if !ok {
		return ""
	}
	return strings.ToLower(peer)
}

// func (z *Zero) HelloUniverse(value string) *Peer {
// 	z.Factory.Rw.Lock()
// 	defer z.Factory.Rw.Unlock()

// 	peer, exists := z.Maps[value]
// 	if !exists {
// 		peer = &Peer{}
// 		if common.IsHexAddress(value) {
// 			peer.Address = z.FormatPeer(value)
// 		} else {
// 			loopringID, err := strconz.ParseInt(value, 10, 64)
// 			if err != nil {
// 				fmt.Printf("Error converting value to int64: %v\n", err)
// 				return nil
// 			}
// 			peer.LoopringID = loopringID
// 		}
// 		z.Peers = append(z.Peers, peer)
// 	}
// 	z.GetENS(peer)
// 	z.GetLoopringENS(peer)
// 	z.GetLoopringID(peer)

// 	z.Maps[peer.Address] = peer
// 	z.Save(peer)
// 	fmt.Printf("	%s %s %s %d\n", peer.Address, peer.ENS, peer.LoopringENS, peer.LoopringID)
// 	return peer
// }
