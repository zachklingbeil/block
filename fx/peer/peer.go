package peer

import (
	"fmt"

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
