package loop

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory      *factory.Factory
	Types        *Type
	Transactions []Transaction
	Map          map[Coordinate]*Tx
}

func Connect(factory *factory.Factory) *Loopring {
	loop := &Loopring{
		Factory: factory,
		Map:     make(map[Coordinate]*Tx),
	}
	loop.CreateTxTable()
	go loop.Listen()
	go loop.BlockByBlock()
	return loop
}

func (l *Loopring) BlockByBlock() {
	current := l.currentBlock()
	for blockNum := current; blockNum >= 1; blockNum-- {
		fmt.Printf("%d\n", blockNum)
		err := l.fetchBlock(blockNum)
		if err != nil {
			fmt.Printf("Error fetching block %d: %v\n", blockNum, err)
			continue
		}
		if blockNum%1000 == 0 {
			if err := l.SaveMap(); err != nil {
				fmt.Printf("Error saving map at block %d: %v\n", blockNum, err)
			}
		}
	}
}

func (l *Loopring) fetchBlock(number int64) error {
	l.Loop()
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/block/getBlock?id=%d", number)
	response, err := l.Factory.Json.In(url, "")
	if err != nil {
		return fmt.Errorf("failed to fetch block %d: %w", number, err)
	}
	var block Block
	if err := json.Unmarshal(response, &block); err != nil {
		return fmt.Errorf("failed to unmarshal block %d: %w", number, err)
	}

	for idx, tx := range block.Transactions {
		coord := l.coordinates(block.Timestamp, int64(idx+1))
		if txMap, ok := tx.(map[string]any); ok {
			txMap["coordinates"] = coord
			flatTx := l.Factory.Json.FlattenMap(txMap, "")
			cleanTx := l.Factory.Json.Cleanup(flatTx)
			block.Transactions[idx] = cleanTx
		}
	}
	l.Unload(block.Transactions)
	l.Simplify()
	l.UpdateMap()
	return nil
}

func (l *Loopring) coordinates(timestamp int64, index int64) Coordinate {
	t := time.UnixMilli(timestamp)
	coordinates := Coordinate{
		Year:        int64(t.Year() - 2015),
		Month:       int64(t.Month()),
		Day:         int64(t.Day()),
		Hour:        int64(t.Hour()),
		Minute:      int64(t.Minute()),
		Second:      int64(t.Second()),
		Millisecond: int64(t.Nanosecond() / 1e6),
		Index:       index,
	}
	return coordinates
}

func (l *Loopring) UpdateMap() {
	l.Factory.Mu.Lock()
	defer l.Factory.Mu.Unlock()
	for _, t := range l.Transactions {
		coord := t.Coordinates
		tx := &Tx{
			Zero:        t.Zero,
			One:         t.One,
			Value:       t.Value,
			Token:       t.Token,
			Fee:         t.Fee,
			FeeToken:    t.FeeToken,
			OneValue:    t.OneValue,
			OneToken:    t.OneToken,
			OneFee:      t.OneFee,
			OneFeeToken: t.OneFeeToken,
			Type:        t.Type,
			Raw:         t.Raw,
		}
		l.Map[coord] = tx
	}
}
