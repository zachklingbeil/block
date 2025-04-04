package loopring

import (
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
}

func NewLoopring(factory *factory.Factory) (*Loopring, error) {
	loopring := &Loopring{
		Factory: factory,
	}

	if err := loopring.CreateTable(); err != nil {
		return nil, fmt.Errorf("failed to create blocks table: %w", err)
	}
	return loopring, nil
}

func (l *Loopring) CreateTable() error {
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
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}
