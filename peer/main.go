package peer

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/zachklingbeil/factory"
)

type Peers struct {
	Factory        *factory.Factory
	LoopringApiKey string
	Db             *sql.DB
}

type Peer struct {
	Address     string
	ENS         string
	LoopringENS string
	LoopringID  string
}

func NewPeers(factory *factory.Factory) (*Peers, error) {
	db, err := factory.Db.Connect("peer")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database peer: %w", err)
	}

	peers := &Peers{
		Factory:        factory,
		LoopringApiKey: os.Getenv("LOOPRING_API_KEY"),
		Db:             db,
	}
	peers.CreateTable()
	return peers, nil
}

func (p *Peers) CreateTable() {
	query := `
    CREATE TABLE IF NOT EXISTS peers (
        address TEXT PRIMARY KEY,       -- Ethereum address
        id TEXT,              		 	 -- Loopring id
        ens TEXT,                       -- [peer].eth
        loopringEns TEXT                -- [peer].loopring.eth
    );`
	p.Db.Exec(query)
}
