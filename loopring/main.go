package loopring

import (
	"log"

	"github.com/zachklingbeil/block/fx"
	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Peers   *fx.Peers
}

func NewLoopring(factory *factory.Factory, peers *fx.Peers) (*Loopring, error) {
	loopring := &Loopring{
		Factory: factory,
		Peers:   peers,
	}
	loopring.Tables()
	loopring.FetchBlocks()
	return loopring, nil
}

func (l *Loopring) Tables() {
	query := `
        CREATE TABLE IF NOT EXISTS loopring (
            block_id BIGINT PRIMARY KEY,
            block_size BIGINT NOT NULL,
            created BIGINT UNIQUE NOT NULL, -- Add UNIQUE constraint
            tx_hash TEXT NOT NULL,
            transactions JSONB NOT NULL
        );

        CREATE TABLE IF NOT EXISTS peers (
            address TEXT PRIMARY KEY,       -- Ethereum address
            id BIGINT,                      -- Loopring account ID
            ens TEXT,                       -- [peer].eth
            loopringEns TEXT                -- [peer].loopring.eth
        );
    `
	if _, err := l.Factory.Db.Exec(query); err != nil {
		log.Printf("failed to create tables: %v", err)
	}
}
