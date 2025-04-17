package loop

import (
	"encoding/json"
	"fmt"
	"time"
)

type NewBlock struct {
	Number       int64 `json:"blockId"`
	Size         int64 `json:"blockSize"`
	Timestamp    int64 `json:"createdAt"`
	Transactions Types `json:"transactions"`
}

type Types struct {
	Deposits       []DW            `json:"Deposit,omitempty"`
	Withdrawals    []DW            `json:"Withdraw,omitempty"`
	Swaps          []Swap          `json:"SpotTrade,omitempty"`
	Transfers      []Transfer      `json:"Transfer,omitempty"`
	Mints          []Mint          `json:"NftMint,omitempty"`
	AccountUpdates []AccountUpdate `json:"AccountUpdate,omitempty"`
	AmmUpdates     []AmmUpdate     `json:"AmmUpdate,omitempty"`
	NftData        []NftData       `json:"NftData,omitempty"`
	TBD            []any           `json:"tbd,omitempty"`
	*json.RawMessage
}

func (l *Loopring) SourceBlocks(lastBlock int64, limit int) ([]NewBlock, error) {
	query := `
        SELECT block, tx
        FROM loopring
        WHERE block > $1
        ORDER BY block ASC
        LIMIT $2;
    `
	rows, err := l.Factory.Db.Query(query, lastBlock, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query blocks: %w", err)
	}
	defer rows.Close()

	var blocks []NewBlock
	for rows.Next() {
		var blockNumber int64
		var txJSON []byte
		if err := rows.Scan(&blockNumber, &txJSON); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		var block NewBlock
		if err := json.Unmarshal(txJSON, &block.Transactions); err != nil {
			return nil, fmt.Errorf("failed to unmarshal transactions for block %d: %w", blockNumber, err)
		}
		block.Number = blockNumber
		blocks = append(blocks, block)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return blocks, nil
}

func (l *Loopring) ProcessBlock() error {

	l.flatten(blockTime, transactions)
	return nil
}

func (l *Loopring) flatten(blockNumber, blockTime int64, newBlock *NewBlock) {
	var flattened []any
	for i, tx := range transactions {
		if txData, ok := tx.(map[string]any); ok {
			coordinates := l.coordinates(blockTime, int64(i+1))
			flatTx := l.Factory.Json.FlattenMap(txData, "")
			flatTx["coordinates"] = coordinates
			cleanedTx := l.Factory.Json.Cleanup(flatTx)
			flattened = append(flattened, cleanedTx)
		} else {
			fmt.Printf("Unexpected transaction format: %+v\n", tx)
		}
	}
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
