package loopring

import (
	"fmt"
	"time"

	"github.com/zachklingbeil/factory"
)

type Loopring struct {
	Factory *factory.Factory
}

type Output struct {
	Number    int64
	Size      int64
	Timestamp int64
	Coords    string
}

func NewLoopring(factory *factory.Factory) *Loopring {
	return &Loopring{
		Factory: factory,
	}
}

// ProcessInputs converts a slice of Block into a slice of Output
func (l *Loopring) ProcessInputs(b []Block) []Output {
	blocks := make([]Output, len(b))

	for i, block := range b {
		blocks[i] = fx(block)
	}
	return blocks
}

// fx processes a single Block into a Output
func fx(block Block) Output {
	t := time.UnixMilli(block.Created)

	// Format the timestamp directly into a string representation of Coordinates
	formattedCoords := fmt.Sprintf("%d.%d.%d.%d.%d.%d.%d",
		t.Year()-2015,      // 0-based year
		int(t.Month()),     // Month
		t.Day(),            // Date of the month (1-31)
		t.Hour(),           // Hour
		t.Minute(),         // Minute
		t.Second(),         // Second
		t.Nanosecond()/1e6) // Millisecond as int64, uncapped

	// Return the Output with the formatted coordinates
	return Output{
		Coords:    formattedCoords,
		Number:    block.Number,
		Size:      block.Size,
		Timestamp: block.Created,
	}
}

// LoadBlocks queries the loopring table, processes the data, and inserts Blocks into the coords table
func (l *Loopring) LoadBlocks() error {
	// Query the loopring table to fetch data
	query := `
        SELECT created, block_id, block_size
        FROM loopring
    `
	rows, err := l.Factory.Db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query loopring table: %w", err)
	}
	defer rows.Close()

	// Create a slice of Block
	var b []Block
	for rows.Next() {
		var block Block
		if err := rows.Scan(&block.Created, &block.Number, &block.Size); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}
		b = append(b, block)
	}

	// Process the b into Blocks
	blocks := l.ProcessInputs(b)

	// Insert the Blocks into the coords table
	for _, block := range blocks {
		if err := l.InsertBlockToCoords(&block); err != nil {
			return fmt.Errorf("failed to insert block into coords table: %w", err)
		}
	}

	return nil
}

// InsertBlockToCoords inserts a block into the coords table
func (l *Loopring) InsertBlockToCoords(o *Output) error {
	query := `
        INSERT INTO coords (block_id, block_size, created, coords)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (coords) DO NOTHING
    `
	if _, err := l.Factory.Db.Exec(query, o.Number, o.Size, o.Timestamp, o.Coords); err != nil {
		return fmt.Errorf("failed to insert block into coords table: %w", err)
	}
	return nil
}

func (l *Loopring) CreateCoordsTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS coords (
            block_id BIGINT NOT NULL,
            block_size BIGINT NOT NULL,
            created BIGINT NOT NULL,
            coords TEXT NOT NULL,
            PRIMARY KEY (coords) -- Use coords as the primary key
        )
    `
	if _, err := l.Factory.Db.Exec(query); err != nil {
		return fmt.Errorf("failed to create coords table: %w", err)
	}
	return nil
}

// func (l *Loopring) OutputCoordsAsJSON() error {
// 	query := `
//         SELECT block_id, block_size, created, coords
//         FROM coords
//     `
// 	rows, err := l.Factory.Db.Query(query)
// 	if err != nil {
// 		return fmt.Errorf("failed to query coords table: %w", err)
// 	}
// 	defer rows.Close()

// 	// Create a slice to hold the results
// 	var results []Output
// 	for rows.Next() {
// 		var output Output
// 		if err := rows.Scan(&output.Number, &output.Size, &output.Timestamp, &output.Coords); err != nil {
// 			return fmt.Errorf("failed to scan row: %w", err)
// 		}
// 		results = append(results, output)
// 	}

// 	// Convert the results to JSON
// 	jsonData, err := json.MarshalIndent(results, "", "  ")
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal results to JSON: %w", err)
// 	}

// 	// Write JSON to a file or print to console
// 	file, err := os.Create("coords.json")
// 	if err != nil {
// 		return fmt.Errorf("failed to create JSON file: %w", err)
// 	}
// 	defer file.Close()

// 	if _, err := file.Write(jsonData); err != nil {
// 		return fmt.Errorf("failed to write JSON to file: %w", err)
// 	}

// 	fmt.Println("Coords table exported to coords.json")
// 	return nil
// }

// type Loopring struct {
// 	TotalNum     int          `json:"totalNum,omitempty"`
// 	Transactions []LoopringTx `json:"transactions,omitempty"`
// }

// type LoopringTx struct {
// 	TxId  int `json:"id,omitempty"`
// 	Block struct {
// 		Number int64 `json:"blockId,omitempty"`
// 		Index  int64 `json:"indexInBlock,omitempty"`
// 	} `json:"blockIdInfo,omitempty"`
// 	Timestamp int64  `json:"timestamp,omitempty"`
// 	Zero      string `json:"senderAddress,omitempty"`
// 	One       string `json:"receiverAddress,omitempty"`
// 	Value     string `json:"amount,omitempty"`
// 	Token     string `json:"symbol,omitempty"`
// 	FeeToken  string `json:"feeTokenSymbol,omitempty"`
// 	FeeValue  string `json:"feeAmount,omitempty"`
// }

// type Transaction struct {
// 	TxType           TxType   `json:"txType"`
// 	AccountID        *int64   `json:"accountId,omitempty"`
// 	Token            *Token   `json:"token,omitempty"`
// 	ToToken          *ToToken `json:"toToken,omitempty"`
// 	Fee              *Fee     `json:"fee,omitempty"`
// 	ValidUntil       *int64   `json:"validUntil,omitempty"`
// 	ToAccountID      *int64   `json:"toAccountId,omitempty"`
// 	ToAccountAddress *string  `json:"toAccountAddress,omitempty"`
// 	StorageID        *int64   `json:"storageId,omitempty"`
// 	OrderA           *Order   `json:"orderA,omitempty"`
// 	OrderB           *Order   `json:"orderB,omitempty"`
// 	Valid            *bool    `json:"valid,omitempty"`
// 	Owner            *string  `json:"owner,omitempty"`
// 	FromAddress      *string  `json:"fromAddress,omitempty"`
// 	ToAddress        *string  `json:"toAddress,omitempty"`
// }

// type Fee struct {
// 	TokenID int64  `json:"tokenId"`
// 	Amount  string `json:"amount"`
// }

// type ToToken struct {
// 	TokenID int64 `json:"tokenId"`
// }

// type Token struct {
// 	TokenID int64   `json:"tokenId"`
// 	NftData *string `json:"nftData,omitempty"`
// 	Amount  string  `json:"amount"`
// }

// type TxType string

// const (
// 	Deposit   TxType = "Deposit"
// 	SpotTrade TxType = "SpotTrade"
// 	Transfer  TxType = "Transfer"
// )
// type Order struct {
// 	StorageID  int64  `json:"storageID"`
// 	AccountID  int64  `json:"accountID"`
// 	AmountS    string `json:"amountS"`
// 	AmountB    string `json:"amountB"`
// 	TokenS     int64  `json:"tokenS"`
// 	TokenB     int64  `json:"tokenB"`
// 	ValidUntil int64  `json:"validUntil"`
// 	Taker      string `json:"taker"`
// 	FeeBips    int64  `json:"feeBips"`
// 	IsAmm      bool   `json:"isAmm"`
// 	NftData    string `json:"nftData"`
// 	FillS      int64  `json:"fillS"`
// 	FilledS    string `json:"filledS"`
// }
