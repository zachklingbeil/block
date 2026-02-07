package fx

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/timefactoryio/block/zero/proto/bytecodedb"
	"github.com/timefactoryio/block/zero/proto/sigprovider"
	"github.com/timefactoryio/block/zero/proto/userops"
)

// output types

type DecodedMethod struct {
	Selector  string            `json:"selector"`
	Name      string            `json:"name,omitempty"`
	Signature string            `json:"signature,omitempty"`
	Params    map[string]string `json:"params,omitempty"`
}

type DecodedEvent struct {
	Index     uint              `json:"index"`
	Address   common.Address    `json:"address"`
	Topic     common.Hash       `json:"topic"`
	Name      string            `json:"name,omitempty"`
	Signature string            `json:"signature,omitempty"`
	Params    map[string]string `json:"params,omitempty"`
}

type DecodedUserOp struct {
	Sender   common.Address `json:"sender"`
	Nonce    string         `json:"nonce"`
	CallData hexutil.Bytes  `json:"callData"`
	Method   *DecodedMethod `json:"method,omitempty"`
}

type DecodedTx struct {
	Hash    common.Hash     `json:"hash"`
	From    common.Address  `json:"from"`
	Deploy  bool            `json:"deploy,omitempty"`
	Method  *DecodedMethod  `json:"method,omitempty"`
	Events  []DecodedEvent  `json:"events"`
	UserOps []DecodedUserOp `json:"userOps,omitempty"`
}

type DecodedBlock struct {
	*Block
	Decoded []DecodedTx `json:"decoded"`
}

var entryPoints = map[common.Address]bool{
	common.HexToAddress("0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789"): true,
	common.HexToAddress("0x0000000071727De22E5E9d8BAf0edAc6f37da032"): true,
}

// gRPC lookups

func (f *Fx) lookupFunctions(ctx context.Context, selectors []string) map[string]*sigprovider.Abi {
	if len(selectors) == 0 || f.Sig == nil {
		return nil
	}
	out := make(map[string]*sigprovider.Abi)
	for _, sel := range selectors {
		resp, err := f.Sig.GetFunctionAbi(ctx, &sigprovider.GetFunctionAbiRequest{
			TxInput: sel,
		})
		if err != nil || len(resp.GetAbi()) == 0 {
			continue
		}
		out[sel] = resp.GetAbi()[0]
	}
	return out
}

func (f *Fx) lookupEvents(ctx context.Context, block *Block) map[string]*sigprovider.Abi {
	if f.Sig == nil {
		return nil
	}

	var reqs []*sigprovider.GetEventAbiRequest
	keys := make([]string, 0)
	seen := make(map[string]bool)

	for _, t := range block.Transactions {
		for _, log := range t.Receipt.Logs {
			if len(log.Topics) == 0 {
				continue
			}
			key := log.Topics[0].Hex()
			if seen[key] {
				continue
			}
			seen[key] = true
			keys = append(keys, key)

			topics := make([]string, len(log.Topics))
			for i, t := range log.Topics {
				topics[i] = t.Hex()
			}

			reqs = append(reqs, &sigprovider.GetEventAbiRequest{
				Data:   hexutil.Encode(log.Data),
				Topics: strings.Join(topics, ","),
			})
		}
	}

	if len(reqs) == 0 {
		return nil
	}

	resp, err := f.Sig.BatchGetEventAbis(ctx, &sigprovider.BatchGetEventAbisRequest{
		Requests: reqs,
	})
	if err != nil {
		return nil
	}

	out := make(map[string]*sigprovider.Abi)
	responses := resp.GetResponses()
	for i, key := range keys {
		if i < len(responses) {
			abis := responses[i].GetAbi()
			if len(abis) > 0 {
				out[key] = abis[0]
			}
		}
	}
	return out
}

func (f *Fx) lookupABI(ctx context.Context, addr common.Address) *abi.ABI {
	if f.DB == nil {
		return nil
	}

	f.RLock()
	if cached, ok := f.abiCache[addr]; ok {
		f.RUnlock()
		return cached
	}
	f.RUnlock()

	resp, err := f.DB.SearchSources(ctx, &bytecodedb.SearchSourcesRequest{
		Bytecode:     addr.Hex(),
		BytecodeType: bytecodedb.BytecodeType_DEPLOYED_BYTECODE,
	})
	if err != nil {
		return nil
	}

	sources := resp.GetSources()
	if len(sources) == 0 {
		return nil
	}

	raw := sources[0].GetAbi()
	if raw == "" {
		return nil
	}

	parsed, err := abi.JSON(strings.NewReader(raw))
	if err != nil {
		return nil
	}

	f.Lock()
	f.abiCache[addr] = &parsed
	f.Unlock()

	return &parsed
}

func (f *Fx) lookupUserOps(ctx context.Context, txHash common.Hash) []DecodedUserOp {
	if f.Ops == nil {
		return nil
	}

	hex := txHash.Hex()
	resp, err := f.Ops.ListUserOps(ctx, &userops.ListUserOpsRequest{
		TransactionHash: &hex,
	})
	if err != nil {
		return nil
	}

	items := resp.GetItems()
	if len(items) == 0 {
		return nil
	}

	ops := make([]DecodedUserOp, 0, len(items))
	for _, item := range items {
		full, err := f.Ops.GetUserOp(ctx, &userops.GetUserOpRequest{
			Hash: item.GetHash(),
		})
		if err != nil {
			continue
		}
		ops = append(ops, DecodedUserOp{
			Sender:   common.HexToAddress(full.GetSender()),
			Nonce:    full.GetNonce(),
			CallData: common.FromHex(full.GetCallData()),
		})
	}
	return ops
}

// decode logic

func (f *Fx) decodeMethod(contractABI *abi.ABI, data []byte) *DecodedMethod {
	if len(data) < 4 {
		return nil
	}

	dm := &DecodedMethod{Selector: hexutil.Encode(data[:4])}

	if contractABI != nil {
		method, err := contractABI.MethodById(data[:4])
		if err == nil {
			dm.Name = method.Name
			dm.Signature = method.Sig
			dm.Params = make(map[string]string)
			if args, err := method.Inputs.Unpack(data[4:]); err == nil {
				for i, input := range method.Inputs {
					dm.Params[input.Name] = fmt.Sprintf("%v", args[i])
				}
			}
		}
	}

	return dm
}

func (f *Fx) decodeEvent(contractABI *abi.ABI, log *types.Log) DecodedEvent {
	de := DecodedEvent{
		Index:   log.Index,
		Address: log.Address,
	}

	if len(log.Topics) == 0 {
		return de
	}

	de.Topic = log.Topics[0]

	if contractABI != nil {
		event, err := contractABI.EventByID(log.Topics[0])
		if err == nil {
			de.Name = event.Name
			de.Signature = event.Sig
			de.Params = make(map[string]string)

			indexed := make([]abi.Argument, 0)
			nonIndexed := make([]abi.Argument, 0)
			for _, input := range event.Inputs {
				if input.Indexed {
					indexed = append(indexed, input)
				} else {
					nonIndexed = append(nonIndexed, input)
				}
			}

			for i, arg := range indexed {
				if i+1 < len(log.Topics) {
					de.Params[arg.Name] = log.Topics[i+1].Hex()
				}
			}

			if len(log.Data) > 0 {
				if args, err := abi.Arguments(nonIndexed).Unpack(log.Data); err == nil {
					for i, arg := range nonIndexed {
						de.Params[arg.Name] = fmt.Sprintf("%v", args[i])
					}
				}
			}
		}
	}

	return de
}

func (f *Fx) collectSelectors(block *Block) []string {
	funcSet := make(map[string]bool)

	for _, t := range block.Transactions {
		if t.Tx.To() != nil && len(t.Tx.Data()) >= 4 {
			funcSet[hexutil.Encode(t.Tx.Data()[:4])] = true
		}
	}

	funcs := make([]string, 0, len(funcSet))
	for k := range funcSet {
		funcs = append(funcs, k)
	}
	return funcs
}

func (f *Fx) Decode(ctx context.Context, block *Block) (*DecodedBlock, error) {
	selectors := f.collectSelectors(block)
	funcAbis := f.lookupFunctions(ctx, selectors)
	eventAbis := f.lookupEvents(ctx, block)

	decoded := make([]DecodedTx, len(block.Transactions))

	for i, t := range block.Transactions {
		dt := DecodedTx{
			Hash:   t.Tx.Hash(),
			From:   t.From,
			Events: make([]DecodedEvent, 0, len(t.Receipt.Logs)),
		}

		if t.Tx.To() == nil {
			dt.Deploy = true
			decoded[i] = dt
			continue
		}

		contractABI := f.lookupABI(ctx, *t.Tx.To())
		dt.Method = f.decodeMethod(contractABI, t.Tx.Data())

		// fallback: sig-provider
		if dt.Method != nil && dt.Method.Name == "" && funcAbis != nil {
			if sigAbi, ok := funcAbis[dt.Method.Selector]; ok {
				dt.Method.Name = sigAbi.GetName()
			}
		}

		for _, log := range t.Receipt.Logs {
			logABI := contractABI
			if log.Address != *t.Tx.To() {
				logABI = f.lookupABI(ctx, log.Address)
			}

			de := f.decodeEvent(logABI, log)

			// fallback: sig-provider
			if de.Name == "" && eventAbis != nil && len(log.Topics) > 0 {
				if sigAbi, ok := eventAbis[log.Topics[0].Hex()]; ok {
					de.Name = sigAbi.GetName()
				}
			}

			dt.Events = append(dt.Events, de)
		}

		if entryPoints[*t.Tx.To()] {
			ops := f.lookupUserOps(ctx, t.Tx.Hash())
			for j := range ops {
				if len(ops[j].CallData) >= 4 {
					opABI := f.lookupABI(ctx, ops[j].Sender)
					ops[j].Method = f.decodeMethod(opABI, ops[j].CallData)
				}
			}
			dt.UserOps = ops
		}

		decoded[i] = dt
	}

	return &DecodedBlock{Block: block, Decoded: decoded}, nil
}
