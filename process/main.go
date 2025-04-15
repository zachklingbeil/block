package process

import (
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Process struct {
	Factory *factory.Factory
	RawTxs  []RawTx
	Txs     []Tx
	Counts  map[string]int
}

func InitProcess(factory *factory.Factory) *Process {
	qtx := 10000

	process := &Process{
		Factory: factory,
		Txs:     make([]Tx, 0, qtx),
		Counts:  make(map[string]int),
	}

	if err := process.LoadRecentBlocks(500); err != nil {
		fmt.Printf("Warning: failed to load blocks: %v\n", err)
	}

	return process
}
