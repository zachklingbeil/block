package loopring

import (
	"database/sql"
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Db      *sql.DB
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

// // Helper function to read transactions from the map for a given block number.
// func (l *Loopring) Read(blockNumber int64) (*Block, bool) {
// 	l.Factory.Mu.Lock()
// 	defer l.Factory.Mu.Unlock()
// 	block, exists := l.Map[blockNumber]
// 	if !exists {
// 		return nil, false
// 	}
// 	return block, true
// }

// // Helper function to update the map with transactions for a given block number.
// func (l *Loopring) Write(block *Block) {
// 	l.Factory.Mu.Lock()
// 	defer l.Factory.Mu.Unlock()
// 	l.Map[block.Number] = block
// }
