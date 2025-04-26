package circuit

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	ens "github.com/wealdtech/go-ens/v3"
)

const (
	byAddress = "https://api3.loopring.io/api/v3/account?owner=%s"
	byId      = "https://api3.loopring.io/api/v3/account?accountId=%d"
	dotLoop   = "https://api3.loopring.io/api/wallet/v3/resolveName?owner=%s"
)

func (c *Circuit) Format(address string) string {
	address = strings.ToLower(address)
	if strings.HasPrefix(address, "0x") || strings.HasSuffix(address, ".eth") {
		return address
	}
	return address
}

// hex -> .eth
func (c *Circuit) GetENS(peer *Peer) *Peer {
	if peer.ENS == "." || (peer.ENS != "" && peer.ENS != "!") {
		// Return immediately if ENS is already set or marked as checked
		return peer
	}

	ensName, err := ens.ReverseResolve(c.Factory.Eth, common.HexToAddress(peer.Address))
	if err != nil || ensName == "" {
		peer.ENS = "." // Mark as checked with no ENS
		return peer
	}
	peer.ENS = c.Format(ensName)
	return peer
}

// ENS -> hex
func (c *Circuit) GetAddress(peer *Peer) *Peer {
	address, err := ens.Resolve(c.Factory.Eth, peer.ENS)
	if err != nil {
		peer.Address = peer.ENS
		return peer
	}
	peer.Address = c.Format(address.Hex())
	return peer
}

// hex -> LoopringENS [.loopring.eth] or "."
func (c *Circuit) GetLoopringENS(peer *Peer) *Peer {
	if peer.LoopringENS == "." || (peer.LoopringENS != "" && peer.LoopringENS != "!") {
		// Return immediately if LoopringENS is already set or marked as checked
		return peer
	}

	url := fmt.Sprintf(dotLoop, peer.Address)
	var response struct {
		Loopring string `json:"data"`
	}

	data, err := c.Factory.Json.In(url, c.ApiKey)
	if err != nil || json.Unmarshal(data, &response) != nil || response.Loopring == "" {
		peer.LoopringENS = "." // Mark as checked with no valid LoopringENS
		return peer
	}

	peer.LoopringENS = c.Format(response.Loopring)
	return peer
}

// hex -> LoopringId or -1
func (c *Circuit) GetLoopringID(peer *Peer) *Peer {
	// Return immediately if LoopringID is already set or marked as checked
	if peer.LoopringID == "." || (peer.LoopringID != "" && peer.LoopringID != "!") {
		return peer
	}

	// Fetch LoopringID if empty or marked as error
	url := fmt.Sprintf(byAddress, peer.Address)
	var response struct {
		ID int64 `json:"accountId"`
	}

	data, err := c.Factory.Json.In(url, c.ApiKey)
	if err != nil || json.Unmarshal(data, &response) != nil || response.ID == 0 {
		peer.LoopringID = "." // Mark as checked with no valid LoopringID
		return peer
	}

	// Assign the resolved LoopringID
	peer.LoopringID = strconv.FormatInt(response.ID, 10)
	return peer
}

// LoopringId -> hex
func (c *Circuit) GetLoopringAddress(peer *Peer) *Peer {
	// Return immediately if Address is already set or marked as checked
	if peer.Address == "." || (peer.Address != "" && peer.Address != "!") {
		return peer
	}

	// Fetch Address if empty or marked as error
	accountID, err := strconv.Atoi(peer.LoopringID)
	if err != nil {
		peer.Address = "." // Mark as checked with no valid Address
		return peer
	}

	url := fmt.Sprintf(byId, accountID)
	var response struct {
		Address string `json:"owner"`
	}

	data, err := c.Factory.Json.In(url, c.ApiKey)
	if err != nil || json.Unmarshal(data, &response) != nil || response.Address == "" {
		peer.Address = "." // Mark as checked with no valid Address
		return peer
	}

	// Assign the resolved Address
	peer.Address = c.Format(response.Address)
	return peer
}
