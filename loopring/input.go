package loopring

import (
	"encoding/json"
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
}

func NewLoopring(factory *factory.Factory) *Loopring {
	return &Loopring{
		Factory: factory,
	}
}

type Block struct {
	Created      int64         `json:"createdAt"`
	Number       int64         `json:"blockId"`
	Size         int64         `json:"blockSize"`
	TxHash       string        `json:"txHash"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	TxType    string `json:"txType"`
	From      int64  `json:"accountId"`
	To        int64  `json:"toAccountId"`
	ToAddress string `json:"toAccountAddress"`
}

// GetBlock fetches a block and inserts it into the database.
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

	// // Process the block's transactions to extract peers
	// l.getPeers(&block)

	// Insert the block into the database
	if err := l.InsertBlock(&block); err != nil {
		return fmt.Errorf("failed to insert block into database: %w", err)
	}
	return nil
}

// getPeers extracts unique addresses from a block's transactions.
// func (l *Loopring) getPeers(block *Block) {
// 	one := make(map[int64]string)

// 	for _, tx := range block.Transactions {
// 		if tx.To != 0 {
// 			one[tx.To] = tx.ToAddress
// 		}
// 		if tx.From != 0 && one[tx.From] == "" {
// 			one[tx.From] = ""
// 		}
// 	}

// 	// Convert unique accounts to a slice of addresses
// 	addresses := make([]string, 0, len(one))
// 	for _, address := range one {
// 		if address != "" {
// 			addresses = append(addresses, address)
// 		}
// 	}

// 	// Pass the addresses to HelloUniverse in a non-blocking goroutine
// 	go func() {
// 		l.Factory.Peer.HelloUniverse(addresses)
// 	}()
// }

// InsertBlock inserts a block into the database.
func (l *Loopring) InsertBlock(block *Block) error {
	query := `
        INSERT INTO loopring (block_id, block_size, created, tx_hash, transactions)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (created) DO NOTHING
    `
	transactions, err := json.Marshal(block.Transactions)
	if err != nil {
		return fmt.Errorf("failed to marshal transactions: %w", err)
	}

	if _, err := l.Factory.Db.Exec(query, block.Number, block.Size, block.Created, block.TxHash, transactions); err != nil {
		return fmt.Errorf("failed to insert block into database: %w", err)
	}
	l.Factory.Json.Print(block.Number)
	return nil
}
