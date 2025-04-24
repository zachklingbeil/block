package peer

import "fmt"

func (p *Peers) LoadPeer() error {
	query := `
        SELECT address, ens, loopringEns, loopringId FROM peers
    `
	rows, err := p.Factory.Pg.Query(query)
	if err != nil {
		return fmt.Errorf("failed to load peers from database: %w", err)
	}
	defer rows.Close()

	p.Factory.Mu.Lock()

	defer p.Factory.Mu.Unlock()
	for rows.Next() {
		var peer Peer
		if err := rows.Scan(&peer.Address, &peer.ENS, &peer.LoopringENS, &peer.LoopringID); err != nil {
			return fmt.Errorf("failed to scan peer row: %w", err)
		}
		p.Map[peer.Address] = &peer
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over peer rows: %w", err)
	}

	fmt.Printf("%d peers\n", len(p.Map))
	p.SavePeers()
	return nil
}
