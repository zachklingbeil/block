package process

import (
	"fmt"

	"github.com/zachklingbeil/factory"
	"github.com/zachklingbeil/peer"
)

type Process struct {
	Factory *factory.Factory
	RawTxs  []any
	Types   *Types
	Txs     []Tx
	Map     map[*Coordinate]*Tx
	Counts  map[string]int
	Peer    *peer.Peers
}

func InitProcess(factory *factory.Factory, peer *peer.Peers) *Process {
	qtx := 10000

	process := &Process{
		Factory: factory,
		Peer:    peer,
		Txs:     make([]Tx, 0, qtx),
		Counts:  make(map[string]int),
		Map:     make(map[*Coordinate]*Tx),
		RawTxs:  make([]any, 0, qtx),
		Types: &Types{
			Deposit:       make([]DW, 0, qtx),
			Withdrawal:    make([]DW, 0, qtx),
			Swaps:         make([]Swap, 0, qtx),
			Transfers:     make([]Transfer, 0, qtx),
			Mints:         make([]Mint, 0, qtx),
			NftData:       make([]NftData, 0, qtx),
			AmmUpdate:     make([]AmmUpdate, 0, qtx),
			AccountUpdate: make([]AccountUpdate, 0, qtx),
			TBD:           make([]any, 0, qtx),
		},
	}
	// if err := process.CreateTxTable(); err != nil {
	// 	fmt.Printf("Warning: failed to create transactions table: %v\n", err)
	// }
	if err := process.LoadRecentBlocks(1000); err != nil {
		fmt.Printf("Warning: failed to load blocks: %v\n", err)
	}
	return process
}

func (p *Process) PopulateTxMap() {
	for i := range p.Txs {
		tx := p.Txs[i]

		txWithoutCoordinates := tx
		txWithoutCoordinates.Coordinates = Coordinate{}

		p.Map[&tx.Coordinates] = &txWithoutCoordinates
	}
}

func (p *Process) ExtractUniqueAddresses() error {
	uniqueAddresses := make(map[string]struct{}) // Use a map to track unique addresses

	// Iterate over the transactions
	for _, tx := range p.Txs {
		// Add the Zero field to the map if it's not nil and is a string
		if address, ok := tx.Zero.(string); ok {
			uniqueAddresses[address] = struct{}{}
		}

		// Add the One field to the map if it's not nil and is a string
		if address, ok := tx.One.(string); ok {
			uniqueAddresses[address] = struct{}{}
		}
	}

	// Convert the map keys to a slice of strings
	addresses := make([]string, 0, len(uniqueAddresses))
	for address := range uniqueAddresses {
		addresses = append(addresses, address)
	}

	// Pass the addresses to Peer.NewBlock in a goroutine
	go func() {
		p.Peer.NewBlock(addresses)
	}()

	return nil
}
