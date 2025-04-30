package value

import (
	"fmt"
	"strconv"
)

func (v *Value) FormatTxValues() error {
	skippedCount := 0 // Counter for skipped transactions

	v.ProcessTxs(func(tx *Tx) error {
		if valueStr, ok := tx.Value.(string); ok {
			formattedValue, err := v.FormatValue(valueStr, tx.Token)
			if err != nil {
				if tokenStr, ok := tx.Token.(string); ok {
					if len(tokenStr) > 2 && tokenStr[:2] == "0x" {
						return nil
					}
					tokenValue, convErr := strconv.Atoi(tokenStr)
					if convErr == nil && tokenValue > 30000 {
						return nil
					}
				}
				skippedCount++ // Increment the skipped counter
				return nil     // Skip this transaction and continue
			}
			tx.Value = formattedValue // Update the Value field with the formatted value
		}
		skippedCount++
		return nil
	})
	fmt.Println("Skipped transactions:", skippedCount)
	return nil
}

func (v *Value) UpdateTxTokens() error {
	v.ProcessTxs(func(tx *Tx) error {
		switch tokenId := tx.Token.(type) {
		case int64:
			// Handle int64 token ID
			token := v.GetTokenById(tokenId)
			tokenValue := strconv.FormatInt(tokenId, 10) // Default to original token ID as a string
			if token != nil {
				tokenValue = token.Token // Use the token's string representation if it exists
			}
			tx.Token = tokenValue // Update the transaction's Token field
		case string:
			// Handle string token ID
			token := v.GetTokenById(tokenId)
			tokenValue := tokenId // Default to the original string token ID
			if token != nil {
				tokenValue = token.Token // Use the token's string representation if it exists
			}
			tx.Token = tokenValue // Update the transaction's Token field
		default:
			// Log unsupported token types
			// fmt.Printf("Unsupported token type: Token=%v (Type=%T)\n", tx.Token, tx.Token)
		}
		return nil
	})
	return nil
}

func (v *Value) UpdateAndFormatTxFeeTokens() error {
	skippedCount := 0 // Counter for skipped transactions

	v.ProcessTxs(func(tx *Tx) error {
		// Update FeeToken
		switch feeTokenId := tx.FeeToken.(type) {
		case int64:
			feeToken := v.GetTokenById(feeTokenId)
			feeTokenValue := strconv.FormatInt(feeTokenId, 10) // Default to original fee token ID as a string
			if feeToken != nil {
				feeTokenValue = feeToken.Token // Use the fee token's string representation if it exists
			}
			tx.FeeToken = feeTokenValue // Update the transaction's FeeToken field
		case string:
			feeToken := v.GetTokenById(feeTokenId)
			feeTokenValue := feeTokenId // Default to the original string fee token ID
			if feeToken != nil {
				feeTokenValue = feeToken.Token // Use the fee token's string representation if it exists
			}
			tx.FeeToken = feeTokenValue // Update the transaction's FeeToken field
		default:
			// Log unsupported fee token types
			// fmt.Printf("Unsupported fee token type: FeeToken=%v (Type=%T)\n", tx.FeeToken, tx.FeeToken)
		}

		// Format Fee
		if feeStr, ok := tx.Fee.(string); ok {
			if feeStr == "0" {
				tx.Fee = nil // Set Fee to nil if it is "0"
				return nil
			}
			formattedFee, err := v.FormatValue(feeStr, tx.FeeToken)
			if err != nil {
				skippedCount++ // Increment the skipped counter
				return nil     // Skip this transaction and continue
			}
			if formattedFee == "0" {
				tx.Fee = nil // Set Fee to nil if the formatted fee is "0"
			} else {
				tx.Fee = formattedFee // Update the Fee field with the formatted value
			}
		} else if feeInt, ok := tx.Fee.(int64); ok && feeInt == 0 {
			tx.Fee = nil // Set Fee to nil if it is 0
		}

		return nil
	})

	fmt.Println("Skipped transactions (Fee):", skippedCount)
	return nil
}
