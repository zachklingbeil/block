package fx

import "github.com/ethereum/go-ethereum/common"

// Action classifies what an event represents in terms of value movement.
type Action uint8

const (
	ActionNone     Action = iota // Approval, Sync, etc — no value movement
	ActionTransfer               // value moved between two parties
	ActionMint                   // new supply: from = 0x0
	ActionBurn                   // destroyed: to = 0x0
	ActionDeposit                // native → wrapped (e.g. WETH)
	ActionWithdraw               // wrapped → native (e.g. WETH)
)

var actionNames = [...]string{"none", "transfer", "mint", "burn", "deposit", "withdraw"}

func (a Action) String() string { return actionNames[a] }

var zx0 common.Address

// Well-known event topic hashes.
var (
	// ERC-20 / ERC-721: Transfer(address,address,uint256)
	sigTransfer = common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
	// ERC-1155: TransferSingle(operator,from,to,id,value)
	sigTransferSingle = common.HexToHash("0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62")
	// ERC-1155: TransferBatch(operator,from,to,ids[],values[])
	sigTransferBatch = common.HexToHash("0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb")
	// WETH: Deposit(address,uint256)
	sigDeposit = common.HexToHash("0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c")
	// WETH: Withdrawal(address,uint256)
	sigWithdrawal = common.HexToHash("0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65")
)

// Classify returns the Action for a decoded event based on its topic0 and addresses.
func Classify(topic0 common.Hash, topics []common.Hash, values map[string]any) Action {
	switch topic0 {
	case sigDeposit:
		return ActionDeposit
	case sigWithdrawal:
		return ActionWithdraw
	case sigTransfer:
		from := extractAddress(topics, values, 1, "from")
		to := extractAddress(topics, values, 2, "to")
		return classifyTransfer(from, to)
	case sigTransferSingle, sigTransferBatch:
		// topics: [sig, operator, from, to]
		from := extractAddress(topics, values, 2, "from")
		to := extractAddress(topics, values, 3, "to")
		return classifyTransfer(from, to)
	default:
		return ActionNone
	}
}

func classifyTransfer(from, to common.Address) Action {
	switch {
	case from == zx0 && to == zx0:
		return ActionNone // degenerate
	case from == zx0:
		return ActionMint
	case to == zx0:
		return ActionBurn
	default:
		return ActionTransfer
	}
}

// extractAddress pulls an address from indexed topics first, falling back to decoded values.
func extractAddress(topics []common.Hash, values map[string]any, topicIdx int, name string) common.Address {
	if topicIdx < len(topics) {
		return common.BytesToAddress(topics[topicIdx].Bytes())
	}
	if v, ok := values[name]; ok {
		if addr, ok := v.(common.Address); ok {
			return addr
		}
	}
	return common.Address{}
}
