package universe

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func (z *Zero) FetchERC20(address common.Address) *One {
	const erc20ABIJSON = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"type":"function"},
		 {"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"},
		 {"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"}]`
	erc20ABI, abiErr := abi.JSON(strings.NewReader(erc20ABIJSON))
	if abiErr != nil {
		return nil
	}
	contract := bind.NewBoundContract(address, erc20ABI, z.Factory.Eth, z.Factory.Eth, z.Factory.Eth)
	callOpts := &bind.CallOpts{Context: z.Factory.Ctx}

	var symbolResult []any
	if err := contract.Call(callOpts, &symbolResult, "symbol"); err != nil {
		return nil
	}
	var decimalsResult []any
	if err := contract.Call(callOpts, &decimalsResult, "decimals"); err != nil {
		return nil
	}

	symbol := ""
	if len(symbolResult) > 0 {
		symbol = fmt.Sprintf("%v", symbolResult[0])
	}
	var decimals int64
	if len(decimalsResult) > 0 {
		switch v := decimalsResult[0].(type) {
		case uint8:
			decimals = int64(v)
		case int64:
			decimals = v
		default:
			decimals = 0
		}
	}

	one := &One{
		Address:  strings.ToLower(address.Hex()),
		Token:    strings.ToLower(symbol),
		Decimals: decimals,
	}
	z.One = append(z.One, one)
	z.Map[one.Address] = one
	z.Maps.Token[one.Token] = one
	return one
}
