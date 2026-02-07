package fx

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (fx *Fx) Block(number *big.Int) (json.RawMessage, error) {
	var tag string
	if number == nil {
		tag = "latest"
	} else {
		tag = hexutil.EncodeBig(number)
	}

	var raw json.RawMessage
	if err := fx.Rpc.CallContext(fx.Context, &raw, "eth_getBlockByNumber", tag, true); err != nil {
		return nil, fmt.Errorf("block: %w", err)
	}
	return raw, nil
}
