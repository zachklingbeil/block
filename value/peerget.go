package value

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wealdtech/go-ens/v3"
)

func (v *Value) Hello(value string) string {
	v.Factory.Rw.RLock()
	peer, exists := v.Map[value]
	v.Factory.Rw.RUnlock()
	if !exists {
		return ""
	}
	if peer.ENS != "" && peer.ENS != "." && peer.ENS != "!" {
		return peer.ENS
	}
	if peer.LoopringENS != "" && peer.LoopringENS != "." && peer.LoopringENS != "!" {
		return peer.LoopringENS
	}
	return peer.Address
}

func (v *Value) Format(address string) string {
	address = strings.ToLower(address)
	if strings.HasPrefix(address, "0x") || strings.HasSuffix(address, ".eth") {
		return address
	}
	return address
}

// ENS -> hex
func (v *Value) GetAddress(peer *Peer) {
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	if peer.Address == "." {
		return
	}
	address, err := ens.Resolve(v.Factory.Eth, peer.ENS)
	if err != nil {
		peer.Address = "."
	} else {
		peer.Address = v.Format(address.Hex())
	}
	v.Save(peer)
	v.Map[peer.Address] = peer
}

// hex -> .eth
func (v *Value) GetENS(peer *Peer) {
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	if peer.ENS == "." {
		return
	}
	ensName, err := ens.ReverseResolve(v.Factory.Eth, common.HexToAddress(peer.Address))
	if err != nil || ensName == "" {
		peer.ENS = "."
	} else {
		peer.ENS = v.Format(ensName)
	}
	v.Save(peer)
	v.Map[peer.ENS] = peer
}

// hex -> LoopringENS [.loopring.eth] or "."
func (v *Value) GetLoopringENS(peer *Peer) {
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	if peer.LoopringENS == "." {
		return
	}
	url := fmt.Sprintf(dotLoop, peer.Address)
	var response struct {
		Loopring string `json:"data"`
	}
	if err := v.input(url, &response); err != nil {
		peer.LoopringENS = "!"
	} else if response.Loopring == "" {
		peer.LoopringENS = "."
	} else {
		peer.LoopringENS = v.Format(response.Loopring)
	}
	v.Save(peer)
	v.Map[peer.LoopringENS] = peer
}

// hex -> LoopringId or "."
func (v *Value) GetLoopringID(peer *Peer) {
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	if peer.LoopringID == "." {
		return
	}
	url := fmt.Sprintf(byAddress, peer.Address)
	var response struct {
		ID int64 `json:"accountId"`
	}
	if err := v.input(url, &response); err != nil {
		peer.LoopringID = "!"
	} else if response.ID == 0 {
		peer.LoopringID = "."
	} else {
		peer.LoopringID = strconv.FormatInt(response.ID, 10)
	}
	v.Save(peer)
	v.Map[peer.LoopringID] = peer
}

// LoopringId -> hex
func (v *Value) GetLoopringAddress(peer *Peer) {
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	if peer.Address == "." {
		return
	}
	url := fmt.Sprintf(byId, peer.LoopringID)
	var response struct {
		Address string `json:"owner"`
	}
	if err := v.input(url, &response); err != nil {
		peer.Address = "!"
	} else if response.Address == "" {
		peer.Address = "."
	} else {
		peer.Address = v.Format(response.Address)
	}
	v.Save(peer)
	v.Map[peer.Address] = peer
}

func (v *Value) ClearInvalidFields() {
	for i := range v.Peers {
		peer := &v.Peers[i]
		if peer.ENS == "." {
			peer.ENS = ""
		}
		if peer.LoopringID == "." {
			peer.LoopringID = ""
		}
		if peer.LoopringENS == "." {
			peer.LoopringENS = ""
		}
		if peer.Address == "." {
			peer.Address = ""
		}
	}
}

func (v *Value) input(url string, response any) error {
	data, err := v.Factory.Json.In(url, "")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, response)
}
