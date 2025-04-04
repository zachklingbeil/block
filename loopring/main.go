package loopring

import (
	"fmt"
	"log"

	"github.com/zachklingbeil/block/fx/peer"
	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Peers   *peer.Peers
}

func NewLoopring(factory *factory.Factory, peers *peer.Peers) (*Loopring, error) {
	loopring := &Loopring{
		Factory: factory,
		Peers:   peers,
	}
	loopring.Tables()
	loopring.FetchBlocks()
	loopring.ExtractPeerInfo()
	loopring.UpdateMissingLoopringENS()
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

        -- Add indexes for frequently queried columns in the peers table
        CREATE INDEX IF NOT EXISTS idx_peers_id ON peers (id);
        CREATE INDEX IF NOT EXISTS idx_peers_ens ON peers (ens);
        CREATE INDEX IF NOT EXISTS idx_peers_loopringEns ON peers (loopringEns);
    `
	if _, err := l.Factory.Db.Exec(query); err != nil {
		log.Printf("failed to create tables: %v", err)
	}
}

func (l *Loopring) UpdateMissingLoopringENS() error {
	// Collect peers that need their LoopringENS updated
	var addressesToUpdate []string
	for address, peer := range l.Peers.Map {
		if peer.LoopringENS == "" {
			addressesToUpdate = append(addressesToUpdate, address)
		}
	}

	totalToUpdate := len(addressesToUpdate)
	if totalToUpdate == 0 {
		fmt.Println("No peers need their LoopringENS updated.")
		return nil
	}

	// Update only the LoopringENS field for each peer and log progress
	for _, address := range addressesToUpdate {
		// Fetch the Loopring ENS for the address
		loopringENS := l.Peers.FetchLoopringENS(address).LoopringENS

		// Ensure LoopringENS is not nil, set to an empty string if missing
		if loopringENS == "" {
			loopringENS = ""
		}

		// Update only the LoopringENS field in the database
		query := `
            UPDATE peers
            SET loopringEns = $1
            WHERE address = $2;
        `
		if _, err := l.Factory.Db.Exec(query, loopringENS, address); err != nil {
			fmt.Printf("Failed to update Loopring ENS for address %s: %v\n", address, err)
			return err
		}

		// Update only the LoopringENS field in the map
		l.Peers.Map[address].LoopringENS = loopringENS

		// Decrement the total and print progress
		totalToUpdate--
		fmt.Printf("%d\n", totalToUpdate)
	}

	fmt.Println("All peers have been updated.")
	return nil
}
