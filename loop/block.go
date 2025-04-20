package loop

import (
	"bytes"
	"encoding/json"
	"time"
)

func (l *Loopring) PrepareBlock(block *Block) []Tx {
	var result []Tx

	baseCoord := l.coordinates(block.Number, block.Timestamp)

	for idx, tx := range block.Transactions.TBD {
		coord := l.updateIndex(baseCoord, int64(idx+1))

		if txMap, ok := tx.(map[string]any); ok {
			txMap["coordinates"] = coord
			flatTx := l.Factory.Json.Flat(txMap, "")
			cleanTx := l.Factory.Json.Cleanup(flatTx)

			txBytes, err := json.Marshal(cleanTx)
			if err != nil {
				continue
			}

			// Process the transaction and append the finalized Tx
			switch txMap["txType"] {
			case "Deposit":
				var dw DW
				if err := json.Unmarshal(txBytes, &dw); err == nil {
					dw.Coordinates = coord
					result = append(result, l.DepositToTx(dw))
				}
			case "Withdraw":
				var dw DW
				if err := json.Unmarshal(txBytes, &dw); err == nil {
					dw.Coordinates = coord
					result = append(result, l.WithdrawToTx(dw))
				}
			case "SpotTrade":
				var swap Swap
				if err := json.Unmarshal(txBytes, &swap); err == nil {
					swap.Coordinates = coord
					result = append(result, l.SwapToTx(swap))
				}
			case "Transfer":
				var transfer Transfer
				if err := json.Unmarshal(txBytes, &transfer); err == nil {
					transfer.Coordinates = coord
					result = append(result, l.TransferToTx(transfer))
				}
			case "NftMint":
				var mint Mint
				if err := json.Unmarshal(txBytes, &mint); err == nil {
					mint.Coordinates = coord
					result = append(result, l.MintToTx(mint))
				}
			case "AccountUpdate":
				var au AccountUpdate
				if err := json.Unmarshal(txBytes, &au); err == nil {
					au.Coordinates = coord
					result = append(result, l.AccountUpdateToTx(au))
				}
			case "AmmUpdate":
				var amm AmmUpdate
				if err := json.Unmarshal(txBytes, &amm); err == nil {
					amm.Coordinates = coord
					result = append(result, l.AmmUpdateToTx(amm))
				}
			case "NftData":
				var nft NftData
				if err := json.Unmarshal(txBytes, &nft); err == nil {
					nft.Coordinates = coord
					result = append(result, l.NftDataToTx(nft))
				}
			default:
				// Handle unknown transaction types if needed
			}
		}
	}

	return result
}

func (l *Loopring) coordinates(blockNumber, timestamp int64) Coordinate {
	t := time.UnixMilli(timestamp)
	return Coordinate{
		Block:       blockNumber,
		Year:        int64(t.Year() - 2015),
		Month:       int64(t.Month()),
		Day:         int64(t.Day()),
		Hour:        int64(t.Hour()),
		Minute:      int64(t.Minute()),
		Second:      int64(t.Second()),
		Millisecond: int64(t.Nanosecond() / 1e6),
	}
}

func (l *Loopring) updateIndex(coord Coordinate, index int64) Coordinate {
	coord.Index = index
	return coord
}

func (l *Loopring) UpdateMap(transactions []Tx) map[Coordinate]*Tx {
	l.Factory.Mu.Lock()
	defer l.Factory.Mu.Unlock()

	delta := make(map[Coordinate]*Tx, len(transactions))
	for _, t := range transactions {
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
		// Only add to delta if new or changed
		if prev, ok := l.Map[coord]; !ok || !txEqual(prev, tx) {
			delta[coord] = tx
			l.Map[coord] = tx
		}
	}
	if len(delta) > 0 {
		_ = l.SaveMap(delta) // handle error as needed
	}
	return delta
}

// Helper to compare two *Tx for equality (can be improved as needed)
func txEqual(a, b *Tx) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Zero == b.Zero &&
		a.One == b.One &&
		a.Value == b.Value &&
		a.Token == b.Token &&
		a.Fee == b.Fee &&
		a.FeeToken == b.FeeToken &&
		a.OneValue == b.OneValue &&
		a.OneToken == b.OneToken &&
		a.OneFee == b.OneFee &&
		a.OneFeeToken == b.OneFeeToken &&
		a.Type == b.Type &&
		bytes.Equal(a.Raw, b.Raw)
}
