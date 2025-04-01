package loopring

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
	Map     map[int64]*Block
	Db      *sql.DB
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

func NewLoopring(factory *factory.Factory) (*Loopring, error) {
	db, err := factory.Db.Connect("loopring")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Loopring database: %w", err)
	}

	loopring := &Loopring{
		Factory: factory,
		Map:     make(map[int64]*Block),
		Db:      db,
	}

	if err := loopring.CreateTable(); err != nil {
		return nil, fmt.Errorf("failed to create blocks table: %w", err)
	}
	return loopring, nil
}

func (l *Loopring) GetBlock(number int) error {
	url := fmt.Sprintf("https://api3.loopring.io/api/v3/block/getBlock?id=%d", number)
	response, err := l.Factory.Json.In(url, "")
	if err != nil {
		return fmt.Errorf("failed to fetch block data for block number %d: %w", number, err)
	}

	var block Block
	if err := json.Unmarshal(response, &block); err != nil {
		return fmt.Errorf("failed to parse block data for block number %d: %w", number, err)
	}

	if err := l.InsertBlock(&block); err != nil {
		return fmt.Errorf("failed to insert block into database: %w", err)
	}

	return nil
}

func (l *Loopring) InsertBlock(block *Block) error {
	query := `
        INSERT INTO blocks (created, block_id, block_size, tx_hash, transactions)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (created) DO NOTHING
    `

	transactionsJSON, err := json.Marshal(block.Transactions)
	if err != nil {
		return fmt.Errorf("failed to marshal transactions: %w", err)
	}

	_, err = l.Db.Exec(query, block.Created, block.Number, block.Size, block.TxHash, transactionsJSON)
	if err != nil {
		return fmt.Errorf("failed to insert block into database: %w", err)
	}

	return nil
}

func (l *Loopring) CreateTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS blocks (
        created BIGINT PRIMARY KEY,
        block_id BIGINT NOT NULL,
        block_size BIGINT NOT NULL,
        tx_hash TEXT NOT NULL,
        transactions JSONB NOT NULL
    );
    `

	_, err := l.Db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create blocks table: %w", err)
	}

	return nil
}
