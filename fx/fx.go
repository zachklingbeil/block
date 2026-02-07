package fx

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/timefactoryio/block/zero"
)

type Fx struct {
	*zero.Zero
	abiCache map[common.Address]*abi.ABI
}

func Init() *Fx {
	return &Fx{
		Zero:     zero.Init(),
		abiCache: make(map[common.Address]*abi.ABI),
	}
}
