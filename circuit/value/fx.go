package value

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wealdtech/go-ens"
)

const (
	byAddress = "https://api3.loopring.io/api/v3/account?owner=%s"
	byId      = "https://api3.loopring.io/api/v3/account?accountId=%d"
	dotLoop   = "https://api3.loopring.io/api/wallet/v3/resolveName?owner=%s"
)

func (v *Value) Format(address string) string {
	address = strings.ToLower(address)
	if strings.HasPrefix(address, "0x") || strings.HasSuffix(address, ".eth") {
		return address
	}
	return address
}

// hex -> .eth
func (v *Value) GetENS(peer *Peer) *Peer {
	if peer.ENS == "." || (peer.ENS != "" && peer.ENS != "!") {
		// Return immediately if ENS is already set or marked as checked
		return peer
	}
	address := common.HexToAddress(peer.Address)
	ensName, err := ens.ReverseResolve(v.Factory.Eth, &address)
	if err != nil || ensName == "" {
		peer.ENS = "." // Mark as checked with no ENS
		return peer
	}
	peer.ENS = v.Format(ensName)
	return peer
}

// ENS -> hex
func (v *Value) GetAddress(peer *Peer) *Peer {
	address, err := ens.Resolve(v.Factory.Eth, peer.ENS)
	if err != nil {
		peer.Address = peer.ENS
		return peer
	}
	peer.Address = v.Format(address.Hex())
	return peer
}

// hex -> LoopringENS [.loopring.eth] or "."
func (v *Value) GetLoopringENS(peer *Peer) *Peer {
	if peer.LoopringENS == "." || (peer.LoopringENS != "" && peer.LoopringENS != "!") {
		// Return immediately if LoopringENS is already set or marked as checked
		return peer
	}

	url := fmt.Sprintf(dotLoop, peer.Address)
	var response struct {
		Loopring string `json:"data"`
	}

	data, err := v.Factory.Json.In(url, os.Getenv("LOOPRING_API_KEY"))
	if err != nil || json.Unmarshal(data, &response) != nil || response.Loopring == "" {
		peer.LoopringENS = "." // Mark as checked with no valid LoopringENS
		return peer
	}

	peer.LoopringENS = v.Format(response.Loopring)
	return peer
}

// hex -> LoopringId or -1
func (v *Value) GetLoopringID(peer *Peer) *Peer {
	if peer.LoopringID == "." || (peer.LoopringID != "" && peer.LoopringID != "!") {
		return peer
	}
	url := fmt.Sprintf(byAddress, peer.Address)
	var response struct {
		ID int64 `json:"accountId"`
	}

	data, err := v.Factory.Json.In(url, os.Getenv("LOOPRING_API_KEY"))
	if err != nil || json.Unmarshal(data, &response) != nil || response.ID == 0 {
		peer.LoopringID = "."
		return peer
	}
	peer.LoopringID = strconv.FormatInt(response.ID, 10)
	return peer
}

// LoopringId -> hex
func (v *Value) GetLoopringAddress(peer *Peer) *Peer {
	if peer.Address == "." || (peer.Address != "" && peer.Address != "!") {
		return peer
	}
	accountID, err := strconv.Atoi(peer.LoopringID)
	if err != nil {
		peer.Address = "." // Mark as checked with no valid Address
		return peer
	}

	url := fmt.Sprintf(byId, accountID)
	var response struct {
		Address string `json:"owner"`
	}

	data, err := v.Factory.Json.In(url, os.Getenv("LOOPRING_API_KEY"))
	if err != nil || json.Unmarshal(data, &response) != nil || response.Address == "" {
		peer.Address = "." // Mark as checked with no valid Address
		return peer
	}
	peer.Address = v.Format(response.Address)
	return peer
}
