package fx

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
	Factory *factory.Factory
	Map     map[string]*Peer
}

type Peer struct {
	Address     string
	ENS         string
	LoopringENS string
	LoopringID  int64
}

func HelloUniverse(factory *factory.Factory) (*Peers, error) {
	peers := &Peers{
		Factory: factory,
		Map:     make(map[string]*Peer),
	}

	var value []Peer
	if err := factory.DiskToMem("peers", &value); err != nil {
		return nil, fmt.Errorf("failed to load peers table: %w", err)
	}

	for _, record := range value {
		peers.Map[record.Address] = &Peer{
			Address:     record.Address,
			ENS:         record.ENS,
			LoopringENS: record.LoopringENS,
			LoopringID:  record.LoopringID,
		}
	}
	return peers, nil
}

func (p *Peers) Update(peer *Peer) error {
	query := `
    INSERT INTO peers (address, id, ens, loopringEns)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (address) DO UPDATE
    SET id = EXCLUDED.id,
        ens = EXCLUDED.ens,
        loopringEns = EXCLUDED.loopringEns;
    `
	_, err := p.Factory.Db.Exec(query, peer.Address, peer.LoopringID, peer.ENS, peer.LoopringENS)
	if err != nil {
		return fmt.Errorf("failed to upsert peer: %w", err)
	}

	p.Map[peer.Address] = peer
	return nil
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

	response, err := p.Factory.Json.In(url, os.Getenv("LOOPRING_API_KEY"))
	if err != nil {
		return &Peer{Address: address}
	}

	if err := json.Unmarshal(response, &resID); err != nil {
		return &Peer{Address: address}
	}
	return &Peer{Address: address, LoopringID: resID.AccountID}
}

func (p *Peers) FetchLoopringAddress(id int64) (*Peer, error) {
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/account?accountId=%d", id)
	var resID struct {
		AccountID int64  `json:"accountId"`
		Owner     string `json:"owner"`
	}

	response, err := p.Factory.Json.In(url, os.Getenv("LOOPRING_API_KEY"))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch address for account ID %d: %w", id, err)
	}

	if err := json.Unmarshal(response, &resID); err != nil {
		return nil, fmt.Errorf("failed to parse address for account ID %d: %w", id, err)
	}
	return &Peer{Address: resID.Owner, LoopringID: id}, nil
}
