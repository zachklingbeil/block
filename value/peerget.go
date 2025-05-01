package value

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wealdtech/go-ens/v3"
)

// Helper function to rebuild the map from v.Peers
func (v *Value) rebuildMap() {
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	v.Map = make(map[string]*Peer)
	for _, p := range v.Peers {
		if p.Address != "" && p.Address != "." && p.Address != "!" {
			v.Map[p.Address] = p
		}
		if p.ENS != "" && p.ENS != "." && p.ENS != "!" {
			v.Map[p.ENS] = p
		}
		if p.LoopringENS != "" && p.LoopringENS != "." && p.LoopringENS != "!" {
			v.Map[p.LoopringENS] = p
		}
		if p.LoopringID != "" && p.LoopringID != "." && p.LoopringID != "!" {
			v.Map[p.LoopringID] = p
		}
	}
}

func (v *Value) HelloUniverse(value string) {
	// First look for existing peer in map, if they don't exist, create new peer and set initial identifying value
	peer, exists := v.Map[value]
	if !exists {
		peer = &Peer{}
		if common.IsHexAddress(value) {
			peer.Address = v.Format(value)
		} else {
			peer.LoopringID = value
		}
		// Add to v.Peers - the sole source of truth
		v.Peers = append(v.Peers, peer)
	}

	v.GetENS(peer)
	v.GetLoopringENS(peer)
	v.GetLoopringID(peer)
	if peer.Address == "" || peer.Address == "!" {
		v.GetLoopringAddress(peer)
	}

	v.rebuildMap()

	v.Factory.Rw.Lock()
	fmt.Printf("%s %s %s\n", peer.ENS, peer.LoopringENS, peer.LoopringID)
	v.Factory.Rw.Unlock()
}

// ENS -> hex
func (v *Value) GetAddress(peer *Peer) {
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	if peer.Address == "." && peer.Address != "" {
		return
	}
	address, err := ens.Resolve(v.Factory.Eth, peer.ENS)
	if err != nil {
		peer.Address = "."
	} else {
		peer.Address = v.Format(address.Hex())
	}
	v.Save(peer)
}

// hex -> .eth
func (v *Value) GetENS(peer *Peer) {
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	if peer.ENS == "." && peer.ENS != "" {
		return
	}
	ensName, err := ens.ReverseResolve(v.Factory.Eth, common.HexToAddress(peer.Address))
	if err != nil || ensName == "" {
		peer.ENS = "."
	} else {
		peer.ENS = v.Format(ensName)
	}
	v.Save(peer)
}

// hex -> LoopringENS [.loopring.eth] or "."
func (v *Value) GetLoopringENS(peer *Peer) {
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	if peer.LoopringENS == "." && peer.LoopringENS != "" {
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
}

// hex -> LoopringId or "."
func (v *Value) GetLoopringID(peer *Peer) {
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	if peer.LoopringID == "." && peer.LoopringID != "" {
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
}

// LoopringId -> hex
func (v *Value) GetLoopringAddress(peer *Peer) {
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	if peer.Address == "." && peer.Address != "" {
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
}
