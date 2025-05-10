package universe

import (
	_ "embed"
	"math/big"
	"strings"
)

// // GetAddress returns the common.Address for a given tokenId.
// func (z *Zero) GetAddress(tokenId int64) string {
// 	if tokenId >= 500 {
// 		return strconv.FormatInt(tokenId, 10)
// 	}

// 	z.Factory.Rw.RLock()
// 	defer z.Factory.Rw.RUnlock()

// 	token, exists := z.Maps.TokenId[tokenId]
// 	if !exists {
// 		log.Printf("Token not found for ID: %d", tokenId)
// 		return strconv.FormatInt(tokenId, 10)
// 	}
// 	return token
// }

// Format formats a string input as a decimal string based on the given decimals.
func (z *Zero) Format(input string, decimals int64) string {
	value := new(big.Int)
	if _, ok := value.SetString(input, 10); !ok {
		return input
	}
	valueStr := value.String()
	if int64(len(valueStr)) <= decimals {
		return "0." + strings.Repeat("0", int(decimals)-len(valueStr)) + valueStr
	}
	intPart := valueStr[:len(valueStr)-int(decimals)]
	fracPart := valueStr[len(valueStr)-int(decimals):]
	result := intPart + "." + fracPart
	result = strings.TrimRight(result, "0")
	return strings.TrimSuffix(result, ".")
}
