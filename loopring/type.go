package loopring

import (
	"encoding/json"
	"fmt"
)

type Block struct {
	Number       int64 `json:"blockId"`
	Size         int64 `json:"blockSize"`
	Timestamp    int64 `json:"createdAt"`
	Transactions []any `json:"transactions"`
}

type Tx struct {
	*json.RawMessage
}

func (l *Loopring) CreateTestTable() error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS test (
		year SMALLINT NOT NULL,
		month SMALLINT NOT NULL,
		day SMALLINT NOT NULL,
		hour SMALLINT NOT NULL,
		minute SMALLINT NOT NULL,
		second SMALLINT NOT NULL,
		millisecond SMALLINT NOT NULL,
		index SMALLINT NOT NULL,
		tx JSONB NOT NULL,
		PRIMARY KEY (year, month, day, hour, minute, second, millisecond, index)
		);
		`
	_, err := l.Factory.Db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}
	return nil
}

// <zero> sent <zero.value> of <one.token> to <one>
// <zero> received <one.value> of <zero.token> from <one>

// <one> sent <one.value> of <zero.token> to <zero>
// <one> received <zero.value> of <one.token> from <zero>
// {
//    "txType": "Swap",

//    "zero": 239514,
//    "zero.value": "127447000",
//    "one.token": 6,

//    "one": 108,
//    "one.value": "1414390000000000000000",
//    "zero.token": 1,

//    "zero.feeBips": 10,
//    "one.feeBips": 0,

//	   "coordinates": {
//	      "year": 10,
//	      "month": 4,
//	      "day": 13,
//	      "hour": 17,
//	      "minute": 0,
//	      "second": 10,
//	      "millisecond": 192,
//	      "index": 1,
//	      "string": "10.4.13.17.0.10.192.1"
//	   }
//	}
