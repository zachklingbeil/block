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
)

// func (c *Circuit) HelloUniverse(address any) *Peer {
// 	c.Factory.Rw.RLock()
// 	var peer *Peer
// 	var ok bool

// 	switch addr := address.(type) {
// 	case string:
// 		peer, ok = c.Map[addr].(*Peer)
// 	case int64:
// 		peer, ok = c.Map[addr].(*Peer)
// 	default:
// 		c.Factory.Rw.RUnlock()
// 		fmt.Printf("Unsupported address type: %T\n", address)
// 		return nil
// 	}
// 	c.Factory.Rw.RUnlock()

// 	if !ok || peer == nil {
// 		fmt.Printf("Peer not found for address: %v\n", address)
// 		return nil
// 	}

// 	c.GetENS(peer, peer.Address)
// 	c.GetLoopringENS(peer, peer.Address)
// 	c.GetLoopringID(peer, peer.Address)

// 	fmt.Printf("%s %s %d\n", peer.ENS, peer.LoopringENS, peer.LoopringID)
// 	return peer
// }

func (c *Circuit) Format(address string) string {
	address = strings.ToLower(address)
	if strings.HasPrefix(address, "0x") || strings.HasSuffix(address, ".eth") {
		return address
	}
	return address
}

// ENS -> hex
func (c *Circuit) GetAddress(peer *Peer, dotEth string) *Peer {
	address, err := ens.Resolve(c.Factory.Eth, dotEth)
	if err != nil {
		peer.Address = dotEth
		return peer
	}
	peer.Address = c.Format(address.Hex())
	return peer
}

// hex -> ENS [.eth] or "."
func (c *Circuit) GetENS(peer *Peer, address string) *Peer {
	addr := common.HexToAddress(address)
	if ensName, err := ens.ReverseResolve(c.Factory.Eth, addr); err != nil || ensName == "" {
		peer.ENS = "."
	} else {
		peer.ENS = c.Format(ensName)
	}
	return peer
}

// hex -> LoopringENS [.loopring.eth] or "."
func (c *Circuit) GetLoopringENS(peer *Peer, address string) *Peer {
	url := fmt.Sprintf(byAddress, address)
	var response struct {
		Loopring string `json:"data"`
	}
	data, err := c.Factory.Json.In(url, c.LoopringApiKey)
	if err != nil || json.Unmarshal(data, &response) != nil || response.Loopring == "" {
		peer.LoopringENS = "."
		return peer
	}
	peer.LoopringENS = c.Format(response.Loopring)
	return peer
}

// hex -> LoopringId or -1
func (c *Circuit) GetLoopringID(peer *Peer, address string) *Peer {
	url := fmt.Sprintf(byAddress, address)
	var response struct {
		ID int64 `json:"accountId"`
	}

	data, err := c.Factory.Json.In(url, c.LoopringApiKey)
	if err != nil {
		fmt.Printf("Failed to fetch LoopringID for address %s (error: %v)\n", address, err)
		peer.LoopringID = "."
		return peer
	}

	if err := json.Unmarshal(data, &response); err != nil || response.ID == 0 {
		fmt.Printf("Unexpected response for address %s: %s\n", address, string(data))
		peer.LoopringID = "."
		return peer
	}
	id := strconv.FormatInt(response.ID, 10)
	peer.LoopringID = id
	return peer
}

// LoopringId -> hex
func (c *Circuit) GetLoopringAddress(peer *Peer, id string) *Peer {
	accountID, err := strconv.Atoi(id)
	if err != nil {
		return peer
	}
	url := fmt.Sprintf(byId, accountID)
	var response struct {
		Address string `json:"owner"`
	}
	if data, err := c.Factory.Json.In(url, c.LoopringApiKey); err == nil && json.Unmarshal(data, &response) == nil {
		peer.Address = c.Format(response.Address)
	} else {
		peer.Address = "!"
	}
	return peer
}
