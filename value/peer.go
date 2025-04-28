package value

import (
	"encoding/json"
	"fmt"
	"log"
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

type Peer struct {
	ENS         string `json:"ens"`
	LoopringENS string `json:"loopringEns"`
	LoopringID  string `json:"loopringId"`
	Address     string `json:"address"`
}

func (v *Value) LoadPeers() error {
	source, err := v.Factory.Data.RB.SMembers(v.Factory.Ctx, "peers").Result()
	if err != nil {
		log.Fatalf("Failed to fetch peers from Redis: %v", err)
	}

	v.Peers = make([]Peer, 0, len(source))
	for _, peerJSON := range source {
		var peer Peer
		if err := json.Unmarshal([]byte(peerJSON), &peer); err != nil {
			log.Printf("Skipping invalid peer: %v", err)
			continue
		}
		v.Peers = append(v.Peers, peer)
	}
	return nil
}

func (v *Value) HelloUniverse(address string) {
	peer := v.GetPeer(address)
	v.GetENS(peer)
	v.GetLoopringID(peer)
	v.GetLoopringENS(peer)
}

func (v *Value) GetPeer(value string) *Peer {
	v.Factory.Rw.RLock()
	peer, exists := v.Map[value]
	v.Factory.Rw.RUnlock()
	if exists {
		return peer
	}
	return v.CreatePeer(value)
}

func (v *Value) CreatePeer(value string) *Peer {
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

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
	v.Map[value] = new
	v.Peers = append(v.Peers, *new)
	return new
}

func (v *Value) Format(address string) string {
	address = strings.ToLower(address)
	if strings.HasPrefix(address, "0x") || strings.HasSuffix(address, ".eth") {
		return address
	}
	return address
}

// hex -> .eth
func (v *Value) GetENS(peer *Peer) *Peer {
	if isValidField(peer.ENS) {
		return peer
	}

	ensName, err := ens.ReverseResolve(v.Factory.Eth, common.HexToAddress(peer.Address))
	if err != nil || ensName == "" {
		peer.ENS = "."
		return peer
	}

	peer.ENS = v.Format(ensName)
	return peer
}

// ENS -> hex
func (v *Value) GetAddress(peer *Peer) *Peer {
	if isValidField(peer.Address) {
		return peer
	}

	address, err := ens.Resolve(v.Factory.Eth, peer.ENS)
	if err != nil {
		peer.Address = "."
		return peer
	}

	peer.Address = v.Format(address.Hex())
	return peer
}

// Helper function to check if a field is valid
func isValidField(field string) bool {
	return field != "" && field != "." && field != "!"
}

func (v *Value) input(url string, response any) error {
	data, err := v.Factory.Json.In(url, "")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, response)
}

// hex -> LoopringENS [.loopring.eth] or "."
func (v *Value) GetLoopringENS(peer *Peer) *Peer {
	if isValidField(peer.LoopringENS) {
		return peer
	}
	url := fmt.Sprintf(dotLoop, peer.Address)
	var response struct {
		Loopring string `json:"data"`
	}
	if err := v.input(url, &response); err != nil {
		peer.LoopringENS = "!"
		return peer
	}
	if response.Loopring == "" {
		peer.LoopringENS = "."
		return peer
	}
	peer.LoopringENS = v.Format(response.Loopring)
	return peer
}

// hex -> LoopringId or "."
func (v *Value) GetLoopringID(peer *Peer) *Peer {
	if isValidField(peer.LoopringID) {
		return peer
	}
	url := fmt.Sprintf(byAddress, peer.Address)
	var response struct {
		ID int64 `json:"accountId"`
	}
	if err := v.input(url, &response); err != nil {
		peer.LoopringID = "!"
		return peer
	}
	if response.ID == 0 {
		peer.LoopringID = "."
		return peer
	}
	peer.LoopringID = strconv.FormatInt(response.ID, 10)
	return peer
}

// LoopringId -> hex
func (v *Value) GetLoopringAddress(peer *Peer) *Peer {
	if isValidField(peer.Address) {
		return peer
	}
	url := fmt.Sprintf(byId, peer.LoopringID)
	var response struct {
		Address string `json:"owner"`
	}
	if err := v.input(url, &response); err != nil {
		peer.Address = "!"
		return peer
	}
	if response.Address == "" {
		peer.Address = "."
		return peer
	}
	peer.Address = v.Format(response.Address)
	return peer
}
