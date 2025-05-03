package value

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wealdtech/go-ens/v3"
)

func (v *Value) Hello(value string) string {
	v.Factory.Rw.RLock()
	peer, exists := v.Map[value]
	v.Factory.Rw.RUnlock()

	if !exists {
		peer = &Peer{}
		if common.IsHexAddress(value) {
			peer.Address = v.Format(value)
		} else {
			peer.LoopringID = value
		}
		v.Factory.Rw.Lock()
		v.Peers = append(v.Peers, peer)
		v.Map[value] = peer
		v.Factory.Rw.Unlock()
		go v.HelloUniverse(peer)
	}

	// Prefer ENS, then LoopringENS, then Address
	switch {
	case peer.ENS != "" && peer.ENS != "." && peer.ENS != "!":
		return peer.ENS
	case peer.LoopringENS != "" && peer.LoopringENS != "." && peer.LoopringENS != "!":
		return peer.LoopringENS
	default:
		return peer.Address
	}
}

func (v *Value) HelloUniverse(peer *Peer) *Peer {
	v.GetENS(peer)
	v.GetLoopringENS(peer)
	v.GetLoopringID(peer)

	v.Factory.Rw.Lock()
	if common.IsHexAddress(peer.Address) {
		v.Map[peer.Address] = peer
	}
	if peer.ENS != "" && peer.ENS != "." && peer.ENS != "!" {
		v.Map[peer.ENS] = peer
	}
	if peer.LoopringENS != "" && peer.LoopringENS != "." && peer.LoopringENS != "!" {
		v.Map[peer.LoopringENS] = peer
	}
	if peer.LoopringID != "" && peer.LoopringID != "." && peer.LoopringID != "!" {
		v.Map[peer.LoopringID] = peer
	}
	v.Save(peer)
	fmt.Printf("	%s %s %s %s\n", peer.Address, peer.ENS, peer.LoopringENS, peer.LoopringID)
	v.Factory.Rw.Unlock()
	return peer
}

// ENS -> hex
func (v *Value) GetAddress(peer *Peer) {
	if peer.Address == "." && peer.Address != "" {
		return
	}
	address, err := ens.Resolve(v.Factory.Eth, peer.ENS)

	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()
	if err != nil {
		peer.Address = "."
	} else {
		peer.Address = v.Format(address.Hex())
	}
}

// hex -> .eth
func (v *Value) GetENS(peer *Peer) {
	if peer.ENS == "." && peer.ENS != "" {
		return
	}
	ensName, err := ens.ReverseResolve(v.Factory.Eth, common.HexToAddress(peer.Address))

	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()
	if err != nil || ensName == "" {
		peer.ENS = "."
	} else {
		peer.ENS = v.Format(ensName)
	}
}

// hex -> LoopringENS [.loopring.eth] or "."
func (v *Value) GetLoopringENS(peer *Peer) {
	if peer.LoopringENS == "." || peer.LoopringENS != "" {
		return
	}
	url := fmt.Sprintf(dotLoop, peer.Address)
	var response struct {
		Loopring string `json:"data"`
	}
	err := v.input(url, &response)
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()
	if err != nil {
		peer.LoopringENS = "."
	} else if response.Loopring == "" {
		peer.LoopringENS = "."
	} else {
		peer.LoopringENS = v.Format(response.Loopring)
	}
}

// hex -> LoopringId or "."
func (v *Value) GetLoopringID(peer *Peer) {
	if peer.LoopringID == "." || (peer.LoopringID != "" && peer.LoopringID != "!") {
		return
	}
	url := fmt.Sprintf(byAddress, peer.Address)
	var response struct {
		ID int64 `json:"accountId"`
	}
	err := v.input(url, &response)
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()
	switch {
	case err != nil:
		peer.LoopringID = "!"
	case response.ID == 0:
		peer.LoopringID = "."
	default:
		peer.LoopringID = strconv.FormatInt(response.ID, 10)
	}
}

// LoopringId -> hex
func (v *Value) GetLoopringAddress(peer *Peer) {
	if peer.Address == "." || (peer.Address != "" && peer.Address != "!") {
		return
	}
	url := fmt.Sprintf(byId, peer.LoopringID)
	var response struct {
		Address string `json:"owner"`
	}
	err := v.input(url, &response)
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()
	switch {
	case err != nil:
		peer.Address = "!"
	case response.Address == "":
		peer.Address = "."
	default:
		peer.Address = v.Format(response.Address)
	}
}
