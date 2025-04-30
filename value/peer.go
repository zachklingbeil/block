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
	source, err := v.Factory.Data.RB.SMembers(v.Factory.Ctx, "peer").Result()
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
		v.Map[peer.Address] = &peer
		v.Map[peer.ENS] = &peer
		v.Map[peer.LoopringENS] = &peer
		v.Map[peer.LoopringID] = &peer
	}
	v.Factory.State.Add("value", "peers", len(v.Peers))
	return nil
}

func (v *Value) Save() {
	for _, peer := range v.Peers {
		peerJSON, err := json.Marshal(peer)
		if err != nil {
			log.Printf("Failed to serialize peer: %v", err)
			continue
		}

		// Add each peer as its own row in the Redis set
		err = v.Factory.Data.RB.SAdd(v.Factory.Ctx, "peer", peerJSON).Err()
		if err != nil {
			log.Printf("Failed to store peer in Redis: %v", err)
		}
	}
}

func (v *Value) ReprocessPeers() {
	addresses := []string{}
	for i := range v.Peers {
		peer := &v.Peers[i]
		if peer.ENS == "!" || peer.LoopringID == "!" || peer.LoopringENS == "!" || peer.Address == "!" {
			addresses = append(addresses, peer.Address)
		}
	}
	fmt.Printf("Reprocessing %d peers...\n", len(addresses))

	v.HelloUniverse(addresses)
	v.Save()
}

func (v *Value) HelloUniverse(addresses []string) {
	remaining := len(addresses) // Track the number of addresses to process

	for _, address := range addresses {
		peer := v.GetPeer(address)

		v.GetENS(peer)
		v.GetLoopringID(peer)
		v.GetLoopringENS(peer)
		v.GetLoopringAddress(peer)

		remaining--
		fmt.Printf("%d %s %s %s\n", remaining, peer.ENS, peer.LoopringENS, peer.LoopringID)
		v.Save()
	}
}

func (v *Value) GetPeer(value string) *Peer {
	v.Factory.Rw.RLock()
	peer, exists := v.Map[value]
	v.Factory.Rw.RUnlock()
	if exists {
		return peer
	}
	return nil // Do not create a new peer
}

func (v *Value) Hello(value string) string {
	v.Factory.Rw.RLock()
	peer, exists := v.Map[value]
	v.Factory.Rw.RUnlock()
	if !exists {
		return ""
	}

	if isValidField(peer.ENS) {
		return peer.ENS
	}
	if isValidField(peer.LoopringENS) {
		return peer.LoopringENS
	}
	return peer.Address
}

func (v *Value) CreatePeer(value string) *Peer {
	v.Factory.Rw.Lock()
	defer v.Factory.Rw.Unlock()

	new := &Peer{}

	switch {
	case common.IsHexAddress(value):
		new.Address = v.Format(value)
	case len(value) > 12 && value[len(value)-13:] == ".loopring.eth":
		new.LoopringENS = v.Format(value)
	case len(value) > 4 && value[len(value)-4:] == ".eth":
		new.ENS = v.Format(value)
	default:
		new.LoopringID = value
		v.GetLoopringAddress(new)
	}

	// Add the new peer to the map and slice
	v.Map[new.Address] = new
	if isValidField(new.ENS) {
		v.Map[new.ENS] = new
	}
	if isValidField(new.LoopringENS) {
		v.Map[new.LoopringENS] = new
	}
	if isValidField(new.LoopringID) {
		v.Map[new.LoopringID] = new
	}
	v.Peers = append(v.Peers, *new)

	// Save the updated peers
	v.Save()

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
	} else {
		peer.ENS = v.Format(ensName)
	}

	v.Factory.Rw.Lock()
	v.Map[peer.ENS] = peer
	v.Factory.Rw.Unlock()
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
	} else {
		peer.Address = v.Format(address.Hex())
	}

	v.Factory.Rw.Lock()
	v.Map[peer.Address] = peer
	v.Factory.Rw.Unlock()
	return peer
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
	} else if response.Loopring == "" {
		peer.LoopringENS = "."
	} else {
		peer.LoopringENS = v.Format(response.Loopring)
	}

	v.Factory.Rw.Lock()
	v.Map[peer.LoopringENS] = peer
	v.Factory.Rw.Unlock()
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
	} else if response.ID == 0 {
		peer.LoopringID = "."
	} else {
		peer.LoopringID = strconv.FormatInt(response.ID, 10)
	}

	v.Factory.Rw.Lock()
	v.Map[peer.LoopringID] = peer
	v.Factory.Rw.Unlock()
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
	} else if response.Address == "" {
		peer.Address = "."
	} else {
		peer.Address = v.Format(response.Address)
	}

	v.Factory.Rw.Lock()
	v.Map[peer.Address] = peer
	v.Factory.Rw.Unlock()
	return peer
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
	v.Save()
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
