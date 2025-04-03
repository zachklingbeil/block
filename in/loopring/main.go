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

func NewLoopring(factory *factory.Factory) (*Loopring, error) {
	db, err := factory.Db.Connect("block")
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

type Block struct {
	Created      int64         `json:"createdAt"`
	Number       int64         `json:"blockId"`
	Size         int64         `json:"blockSize"`
	TxHash       string        `json:"txHash"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	TxType    TxType `json:"txType"`
	From      int64  `json:"accountId"`
	To        int64  `json:"toAccountId"`
	ToAddress string `json:"toAccountAddress"`
}

type TxType string

const (
	Transfer TxType = "Transfer, Deposit, Withdraw"
)
