package ethereum

import (
	"database/sql"
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Ethereum struct {
	Factory *factory.Factory
	Db      *sql.DB
}

func NewEthereum(factory *factory.Factory) (*Ethereum, error) {
	db, err := factory.Db.Connect("ethereum")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Loopring database: %w", err)
	}
	return &Ethereum{
		Factory: factory,
		Db:      db,
	}, nil
}
