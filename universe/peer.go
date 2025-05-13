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
		peer.ENS = z.Format.Peer(ensName)
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
			peer.LoopringENS = z.Format.Peer(resp.Loopring)
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
		peer.Address = z.Format.Peer(response.Address)
	}
}

func (z *Zero) input(url string, response any) error {
	data, err := z.Factory.Json.In(url, "")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, response)
}

func (z *Zero) HelloUniverse(value string) *One {
	z.Factory.Rw.Lock()
	defer z.Factory.Rw.Unlock()

	peer, exists := z.Map[value]
	if !exists {
		peer = &One{}
		if common.IsHexAddress(value) {
			peer.Address = strings.ToLower(value)
		}
		z.One = append(z.One, peer)
		z.Map[peer.Address] = peer
	} else if peer.Address == "" && common.IsHexAddress(value) {
		peer.Address = strings.ToLower(value)
		z.Map[peer.Address] = peer
	}

	z.GetENS(peer)
	fmt.Printf("	%s %s\n", peer.Address, peer.ENS)
	return peer
}

// func (z *Zero) HelloUniverse(value string) *One {
// 	z.Factory.Rw.Lock()
// 	defer z.Factory.Rw.Unlock()

// 	peer, exists := z.Map[value]
// 	if !exists {
// 		peer = &One{}
// 		if common.IsHexAddress(value) {
// 			peer.Address = z.Format.Peer(value)
// 		} else {
// 			loopringID, err := strconv.ParseInt(value, 10, 64)
// 			if err != nil {
// 				fmt.Printf("Error converting value to int64: %v\n", err)
// 				return nil
// 			}
// 			peer.LoopringID = loopringID
// 		}
// 		z.One = append(z.One, peer)
// 	}
// 	z.GetENS(peer)
// 	z.GetLoopringENS(peer)
// 	z.GetLoopringID(peer)

// 	z.Map[peer.Address] = peer
// 	fmt.Printf("	%s %s %s %d\n", peer.Address, peer.ENS, peer.LoopringENS, peer.LoopringID)
// 	return peer
// }
