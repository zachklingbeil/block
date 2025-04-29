package value

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
)

type Token struct {
	Token    string `json:"token,omitempty"`
	Address  string `json:"address,omitempty"`
	Decimals string `json:"decimals,omitempty"`
	TokenId  string `json:"tokenId,omitempty"`
	TokenInt int64  `json:"tokenInt,omitempty"`
}

func (v *Value) LoadTokens() error {
	source, err := v.Factory.Data.RB.SMembers(v.Factory.Ctx, "token").Result()
	if err != nil {
		log.Fatalf("Failed to fetch tokens from Redis: %v", err)
	}

	v.Tokens = make([]Token, 0, len(source))
	for _, tokenJSON := range source {
		var token Token
		if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
			log.Printf("Skipping invalid token: %v (data: %s)", err, tokenJSON)
			continue
		}
		v.Tokens = append(v.Tokens, token)
		v.TokenMap[token.TokenInt] = &token
		v.TokenMap[token.TokenId] = &token
		v.TokenMap[token.Token] = &token
	}
	return nil
}

func (v *Value) TokenIdToString(tokenInt int64) string {
	v.Factory.Rw.RLock()
	defer v.Factory.Rw.RUnlock()

	token, exists := v.TokenMap[tokenInt]
	if exists {
		return token.TokenId
	}
	return "!"
}

func (v *Value) GetTokenById(tokenInt int64) *Token {
	if tokenInt > 1000 {
		return nil
	}
	v.Factory.Rw.RLock()
	defer v.Factory.Rw.RUnlock()

	token, exists := v.TokenMap[tokenInt]
	if !exists {
		log.Printf("Token not found for ID: %d", tokenInt)
		return nil
	}
	return token
}

// FormatValue formats a string input as a decimal string based on the token's decimals.
func (v *Value) FormatValue(input string, key any) (string, error) {
	v.Factory.Rw.RLock()
	token, exists := v.TokenMap[key]
	v.Factory.Rw.RUnlock()
	if !exists {
		return "", fmt.Errorf("token not found for key: %v", key)
	}
	// Parse the input string into a big.Int
	value := new(big.Int)
	_, ok := value.SetString(input, 10)
	if !ok {
		return "", fmt.Errorf("invalid input string: %s", input)
	}

	// Parse the decimals from the token
	decimals, err := strconv.Atoi(token.Decimals)
	if err != nil {
		return "", fmt.Errorf("invalid decimals value in token: %v", err)
	}

	// Convert the big.Int value to a string
	valueStr := value.String()

	// Handle cases where the value length is less than or equal to the decimals
	if len(valueStr) <= decimals {
		// Add leading zeros to match the decimals
		paddedValue := strings.Repeat("0", decimals-len(valueStr)+1) + valueStr
		// Insert the decimal point at the correct position
		result := "0." + paddedValue
		return strings.TrimRight(result, "0"), nil
	}

	// Insert the decimal point at the correct position
	intPart := valueStr[:len(valueStr)-decimals]
	fracPart := valueStr[len(valueStr)-decimals:]

	// Combine the integer and fractional parts
	result := intPart + "." + fracPart

	// Trim trailing zeroes from the fractional part
	result = strings.TrimRight(result, "0")

	// If the result ends with a ".", remove it
	result = strings.TrimSuffix(result, ".")

	return result, nil
}

// // FormatValue formats a string input as a decimal string based on the token's decimals.
// // Uses v.TokenMap to look up the token based on the provided key.
// func (v *Value) FormatValue(input string, key any) (string, error) {
// 	// Look up the token in the TokenMap
// 	v.Factory.Rw.RLock()
// 	token, exists := v.TokenMap[key]
// 	v.Factory.Rw.RUnlock()
// 	if !exists {
// 		return "", fmt.Errorf("token not found for key: %v", key)
// 	}

// 	// Parse the input string into a big.Int
// 	value := new(big.Int)
// 	if _, ok := value.SetString(input, 10); !ok {
// 		return "", fmt.Errorf("invalid input string: %s", input)
// 	}

// 	// Parse the decimals from the token
// 	decimals, err := strconv.Atoi(token.Decimals)
// 	if err != nil {
// 		return "", fmt.Errorf("invalid decimals value in token: %v", err)
// 	}

// 	// Convert the big.Int value to a string
// 	valueStr := value.String()

// 	// Handle cases where the value length is less than or equal to the decimals
// 	if len(valueStr) <= decimals {
// 		// Pad with leading zeros and insert the decimal point
// 		paddedValue := strings.Repeat("0", decimals-len(valueStr)+1) + valueStr
// 		return "0." + paddedValue, nil
// 	}

// 	// Split the value into integer and fractional parts
// 	intPart := valueStr[:len(valueStr)-decimals]
// 	fracPart := valueStr[len(valueStr)-decimals:]

// 	// Process the fractional part to handle ellipses and trailing zeroes
// 	processedFracPart := processFractionalPart(fracPart)

// 	// Combine the integer and processed fractional parts
// 	return intPart + "." + processedFracPart, nil
// }

// // Helper function to process the fractional part
// func processFractionalPart(fracPart string) string {
// 	var result strings.Builder
// 	zeroCount := 0
// 	lastNonZeroIndex := strings.LastIndexFunc(fracPart, func(r rune) bool {
// 		return r != '0'
// 	})

// 	for i, r := range fracPart {
// 		if r == '0' {
// 			zeroCount++
// 		} else {
// 			// If there were 3 or more consecutive zeroes, add an ellipsis
// 			if zeroCount >= 3 {
// 				result.WriteString("...")
// 			} else {
// 				// Otherwise, append the zeroes as-is
// 				result.WriteString(strings.Repeat("0", zeroCount))
// 			}
// 			zeroCount = 0
// 			result.WriteRune(r)
// 		}

// 		// Handle the final digit after a long sequence of zeroes
// 		if i == lastNonZeroIndex && zeroCount >= 3 {
// 			result.WriteString("...")
// 		}
// 	}

// 	// Append any remaining zeroes if they are fewer than 3
// 	if zeroCount > 0 && zeroCount < 3 {
// 		result.WriteString(strings.Repeat("0", zeroCount))
// 	}

// 	return result.String()
// }
