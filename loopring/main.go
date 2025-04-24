package loopring

import (
	"encoding/json"
	"fmt"

	"github.com/zachklingbeil/block/circuit"
	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Circuit *circuit.Circuit
}

func Connect(factory *factory.Factory, circuit *circuit.Circuit) *Loopring {
	loop := &Loopring{
		Factory: factory,
		Circuit: circuit,
	}
	go loop.Listen()
	return loop
}

func (l *Loopring) Loop() {
	past, distance := l.Distance()
	for blockNumber := past + distance; blockNumber > past; blockNumber-- {
		l.BlockByBlock(blockNumber)
	}
}

func (l *Loopring) BlockByBlock(blockNumber int64) {
	input := l.FetchBlock(blockNumber)
	transactions, block := l.Circuit.Coordinates(input)
	txs := l.ProcessBlock(transactions)
	block.Ones = txs
	blockJSON, _ := json.Marshal(block)
	l.Factory.Redis.SAdd(l.Factory.Ctx, "blocks", blockJSON).Err()
	fmt.Printf("%d\n", blockNumber)
}

func (l *Loopring) Distance() (int64, int64) {
	current := l.currentBlock()
	past, _ := l.getHistory()

	distance := current - past
	if distance > 0 {
		return past, distance
	}
	return past, 0
}
