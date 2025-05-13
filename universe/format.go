package universe

import (
	"math/big"
	"strings"
)

// Format formats a string input as a decimal string based on the given decimals.
func (f *Format) Value(input string, decimals int64) string {
	value := new(big.Int)
	_, ok := value.SetString(input, 10)
	if !ok {
		return input
	}
	valueStr := value.String()
	dec := int(decimals)
	if len(valueStr) <= dec {
		paddedValue := strings.Repeat("0", dec-len(valueStr)+1) + valueStr
		result := "0." + paddedValue
		return strings.TrimRight(result, "0")
	}

	left := valueStr[:len(valueStr)-dec]
	right := valueStr[len(valueStr)-dec:]
	result := left + "." + right
	result = strings.TrimRight(result, "0")
	result = strings.TrimSuffix(result, ".")
	return result
}

func (f *Format) Peer(address string) string {
	address = strings.ToLower(address)
	if strings.HasPrefix(address, "0x") || strings.HasSuffix(address, ".eth") {
		return address
	}
	return address
}
