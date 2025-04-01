package peer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wealdtech/go-ens/v3"
)

// func (p *Peers) HelloUniverse(address string) {
// 	var peer *Peer

// 	if strings.HasPrefix(address, "0x") {
// 		peer = p.FetchENS(address)
// 	} else {
// 		peer = p.FetchAddress(address)
// 	}

// 	if peer != nil {
// 		p.Factory.Json.In(p.FormatAddress(peer.Address), peer)
// 	}

// 	peer = p.FetchLoopringENS(address)
// 	if peer != nil {
// 		p.Factory.Json.In(peer.LoopringENS, peer)
// 	}

// 	peer = p.FetchLoopringID(address)
// 	if peer != nil {
// 		p.Factory.Json.In(peer.LoopringID, peer)
// 	}
// }

func (p *Peers) FormatAddress(address string) string {
	if strings.HasPrefix(address, "0x") {
		return "0x" + strings.ToUpper(address[2:])
	}
	return address
}

func (p *Peers) FetchAddress(name string) *Peer {
	address, err := ens.Resolve(p.Factory.Eth, name)
	if err != nil {
		return &Peer{Address: name}
	}
	formattedAddress := p.FormatAddress(address.Hex())
	return &Peer{Address: formattedAddress, ENS: name}
}

func (p *Peers) FetchENS(address string) *Peer {
	addr := common.HexToAddress(address)
	name, err := ens.ReverseResolve(p.Factory.Eth, addr)
	if err != nil {
		return &Peer{Address: address}
	}
	return &Peer{Address: address, ENS: name}
}

func (p *Peers) FetchLoopringENS(address string) *Peer {
	url := fmt.Sprintf("https://api3.loopring.io/api/wallet/v3/resolveName?owner=%s", address)
	var resName struct {
		Loopring string `json:"data"`
	}

	// Handle the response and error
	response, err := p.Factory.Json.In(url, "")
	if err != nil {
		return &Peer{Address: address}
	}

	// Unmarshal the response into the struct
	if err := json.Unmarshal(response, &resName); err != nil {
		return &Peer{Address: address}
	}

	return &Peer{Address: address, LoopringENS: resName.Loopring}
}

func (p *Peers) FetchLoopringID(address string) *Peer {
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/account?owner=%s", address)
	var resID struct {
		AccountID int64  `json:"accountId"`
		Owner     string `json:"owner"`
	}

	// Handle the response and error
	response, err := p.Factory.Json.In(url, p.LoopringApiKey)
	if err != nil {
		return &Peer{Address: address}
	}

	// Unmarshal the response into the struct
	if err := json.Unmarshal(response, &resID); err != nil {
		return &Peer{Address: address}
	}

	return &Peer{Address: address, LoopringID: fmt.Sprintf("%d", resID.AccountID)}
}
