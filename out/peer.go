package out

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wealdtech/go-ens/v3"
	"github.com/zachklingbeil/factory"
)

type Peers struct {
	Factory        *factory.Factory
	LoopringApiKey string
}

type Peer struct {
	Address     string
	ENS         string
	LoopringENS string
	LoopringID  string
}

func NewPeers(factory *factory.Factory) (*Peers, error) {
<<<<<<< HEAD
	return &Peers{
		Factory:        factory,
		LoopringApiKey: os.Getenv("LOOPRING_API_KEY"),
	}, nil
=======
	peers := &Peers{
		Factory:        factory,
		LoopringApiKey: os.Getenv("LOOPRING_API_KEY"),
	}
	if err := peers.PeerTable(); err != nil {
		return nil, err
	}
	return peers, nil
}

// PeerTable creates the addresses table if it doesn't already exist.
func (p *Peers) PeerTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS peers (
        address TEXT PRIMARY KEY,       -- Ethereum address
        id BIGINT,                      -- Loopring account ID
        ens TEXT,                       -- [peer].eth
        loopringEns TEXT                -- [peer].loopring.eth
    );`
	if _, err := p.Factory.Db.Exec(query); err != nil {
		return fmt.Errorf("failed to create addresses table: %w", err)
	}
	return nil
>>>>>>> simple
}

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

	response, err := p.Factory.Json.In(url, "")
	if err != nil {
		return &Peer{Address: address}
	}

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

	response, err := p.Factory.Json.In(url, p.LoopringApiKey)
	if err != nil {
		return &Peer{Address: address}
	}

	if err := json.Unmarshal(response, &resID); err != nil {
		return &Peer{Address: address}
	}
	return &Peer{Address: address, LoopringID: fmt.Sprintf("%d", resID.AccountID)}
}
