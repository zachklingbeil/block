package universe

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wealdtech/go-ens/v3"
)

const (
	byAddress = "https://api3.loopring.io/api/v3/account?owner=%s"
	byId      = "https://api3.loopring.io/api/v3/account?accountId=%s"
	dotLoop   = "https://api3.loopring.io/api/wallet/v3/resolveName?owner=%s"
)

// hex -> .eth
func (z *Zero) GetENS(peer *One) {
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
func (z *Zero) GetLoopringENS(peer *One) {
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
func (z *Zero) GetLoopringID(peer *One) {
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

// LoopringId -> hex
func (z *Zero) GetLoopringAddress(peer *One) {
	if peer.Address == "." || (peer.Address != "" && peer.Address != "!") {
		return
	}
	url := fmt.Sprintf(byId, strconv.FormatInt(peer.LoopringID, 10))
	var response struct {
		Address string `json:"owner"`
	}
	err := z.input(url, &response)
	z.Factory.Rw.Lock()
	defer z.Factory.Rw.Unlock()
	switch {
	case err != nil:
		peer.Address = "!"
	case response.Address == "":
		peer.Address = "."
	default:
		peer.Address = z.FormatPeer(response.Address)
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
