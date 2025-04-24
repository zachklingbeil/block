package peer

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
)

func (p *Peers) Format(address string) string {
	address = strings.ToLower(address)
	if strings.HasPrefix(address, "0x") || strings.HasSuffix(address, ".eth") {
		return address
	}
	return address
}

// ENS -> hex
func (p *Peers) GetAddress(peer *Peer, dotEth string) *Peer {
	address, err := ens.Resolve(p.Factory.Eth, dotEth)
	if err != nil {
		peer.Address = dotEth
		return peer
	}
	peer.Address = p.Format(address.Hex())
	return peer
}

// hex -> ENS [.eth] or "."
func (p *Peers) GetENS(peer *Peer, address string) *Peer {
	addr := common.HexToAddress(address)
	if ensName, err := ens.ReverseResolve(p.Factory.Eth, addr); err != nil || ensName == "" {
		peer.ENS = "."
	} else {
		peer.ENS = p.Format(ensName)
	}
	return peer
}

// hex -> LoopringENS [.loopring.eth] or "."
func (p *Peers) GetLoopringENS(peer *Peer, address string) *Peer {
	url := fmt.Sprintf(byAddress, address)
	var response struct {
		Loopring string `json:"data"`
	}
	data, err := p.Factory.Json.In(url, p.LoopringApiKey)
	if err != nil || json.Unmarshal(data, &response) != nil || response.Loopring == "" {
		peer.LoopringENS = "."
		return peer
	}
	peer.LoopringENS = p.Format(response.Loopring)
	return peer
}

// hex -> LoopringId or -1
func (p *Peers) GetLoopringID(peer *Peer, address string) *Peer {
	url := fmt.Sprintf(byAddress, address)
	var response struct {
		ID int64 `json:"accountId"`
	}

	data, err := p.Factory.Json.In(url, p.LoopringApiKey)
	if err != nil {
		fmt.Printf("Failed to fetch LoopringID for address %s (error: %v)\n", address, err)
		peer.LoopringID = -1
		return peer
	}

	if err := json.Unmarshal(data, &response); err != nil || response.ID == 0 {
		fmt.Printf("Unexpected response for address %s: %s\n", address, string(data))
		peer.LoopringID = -1
		return peer
	}
	peer.LoopringID = response.ID
	return peer
}

// LoopringId -> hex
func (p *Peers) GetLoopringAddress(peer *Peer, id string) *Peer {
	accountID, err := strconv.Atoi(id)
	if err != nil {
		return peer
	}
	url := fmt.Sprintf(byId, accountID)
	var response struct {
		Address string `json:"owner"`
	}
	if data, err := p.Factory.Json.In(url, p.LoopringApiKey); err == nil && json.Unmarshal(data, &response) == nil {
		peer.Address = p.Format(response.Address)
	} else {
		peer.Address = "!"
	}
	return peer
}
