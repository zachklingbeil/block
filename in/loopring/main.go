package loopring

import (
	"database/sql"
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Db      *sql.DB
	Blocks  []Block
}

// NewLoopring initializes a new Loopring instance and ensures the database table exists.
func NewLoopring(factory *factory.Factory) (*Loopring, error) {
	db, err := factory.Db.Connect("loopring")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Loopring database: %w", err)
	}

	loopring := &Loopring{
		Factory: factory,
		Db:      db,
	}

	if err := loopring.CreateTable(); err != nil {
		return nil, fmt.Errorf("failed to create blocks table: %w", err)
	}
	return loopring, nil
}
