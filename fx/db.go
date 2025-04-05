package fx

import (
	"log"

	"github.com/zachklingbeil/block/fx/peer"
	"github.com/zachklingbeil/factory"
)

type Fx struct {
	Factory *factory.Factory
	Peers   *peer.Peers
}

func Tables(fx *Fx) {
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

        -- Add indexes for frequently queried columns in the peers table
        CREATE INDEX IF NOT EXISTS idx_peers_id ON peers (id);
        CREATE INDEX IF NOT EXISTS idx_peers_ens ON peers (ens);
        CREATE INDEX IF NOT EXISTS idx_peers_loopringEns ON peers (loopringEns);
    `
	if _, err := fx.Factory.Db.Exec(query); err != nil {
		log.Printf("failed to create tables: %v", err)
	}
}
