package fx

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/timefactoryio/block/zero"
)

type Fx struct {
	*zero.Zero
	sync.RWMutex
	abis    map[common.Address]*abi.ABI
	events  map[common.Hash]*abi.Event
	methods map[string]*abi.Method
}

func Init(url string) *Fx {
	return &Fx{
		Zero:    zero.Init(url),
		abis:    make(map[common.Address]*abi.ABI),
		events:  make(map[common.Hash]*abi.Event),
		methods: make(map[string]*abi.Method),
	}
}

var (
	TopicTransfer = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
	TopicApproval = crypto.Keccak256Hash([]byte("Approval(address,address,uint256)"))
)

func (fx *Fx) Test() error {
	// Fetch block
	block, err := fx.Block(nil)
	if err != nil {
		return fmt.Errorf("Block: %w", err)
	}

	fmt.Printf("Block #%s (%s) — %d txs\n", block.Header.Number, block.Header.Hash(), len(block.Transactions))

	// Write raw block
	rawOutput, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		return fmt.Errorf("Marshal block: %w", err)
	}

	if err := os.WriteFile("../output/block.json", rawOutput, 0644); err != nil {
		return fmt.Errorf("WriteFile block: %w", err)
	}

	fmt.Println("Block written to output/block.json")

	// Decode block
	decoded, err := fx.Decode(fx.Context, block)
	if err != nil {
		return fmt.Errorf("Decode: %w", err)
	}

	fmt.Printf("Block %s: %d transactions\n", decoded.Hash, len(decoded.Txs))

	for i, dt := range decoded.Txs {
		switch {
		case dt.Deploy:
			fmt.Printf("  tx[%d] %s DEPLOY to %s\n", i, dt.Hash, dt.ContractAddress)
		case dt.Method != nil && dt.Method.Name != "":
			fmt.Printf("  tx[%d] %s %s() — %d events\n", i, dt.Hash, dt.Method.Name, len(dt.Events))
		case dt.Method != nil:
			fmt.Printf("  tx[%d] %s %s — %d events\n", i, dt.Hash, dt.Method.Selector, len(dt.Events))
		default:
			fmt.Printf("  tx[%d] %s transfer — %d events\n", i, dt.Hash, len(dt.Events))
		}

		if len(dt.UserOps) > 0 {
			fmt.Printf("    %d user operations\n", len(dt.UserOps))
		}
	}

	// Write decoded block
	decodedOutput, err := json.MarshalIndent(decoded, "", "  ")
	if err != nil {
		return fmt.Errorf("Marshal decoded: %w", err)
	}

	if err := os.WriteFile("../output/decoded.json", decodedOutput, 0644); err != nil {
		return fmt.Errorf("WriteFile decoded: %w", err)
	}

	fmt.Println("Decoded block written to output/decoded.json")
	return nil
}
