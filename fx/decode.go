package fx

import (
	"context"
	"fmt"
	"maps"
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
	// transaction fields
	Hash           common.Hash      `json:"hash"`
	Nonce          uint64           `json:"nonce"`
	From           common.Address   `json:"from"`
	To             *common.Address  `json:"to,omitempty"`
	Value          *hexutil.Big     `json:"value"`
	Input          hexutil.Bytes    `json:"input"`
	Type           uint8            `json:"type"`
	Gas            uint64           `json:"gas"`
	GasPrice       *hexutil.Big     `json:"gasPrice,omitempty"`
	MaxFeePerGas   *hexutil.Big     `json:"maxFeePerGas,omitempty"`
	MaxPriorityFee *hexutil.Big     `json:"maxPriorityFeePerGas,omitempty"`
	ChainID        *hexutil.Big     `json:"chainId,omitempty"`
	AccessList     types.AccessList `json:"accessList,omitempty"`
	BlobGas        uint64           `json:"blobGas,omitempty"`
	BlobGasFeeCap  *hexutil.Big     `json:"maxFeePerBlobGas,omitempty"`
	BlobHashes     []common.Hash    `json:"blobVersionedHashes,omitempty"`
	V              *hexutil.Big     `json:"v"`
	R              *hexutil.Big     `json:"r"`
	S              *hexutil.Big     `json:"s"`

	// receipt fields
	Status            uint64          `json:"status"`
	GasUsed           uint64          `json:"gasUsed"`
	EffectiveGasPrice *hexutil.Big    `json:"effectiveGasPrice"`
	CumulativeGasUsed uint64          `json:"cumulativeGasUsed"`
	ContractAddress   *common.Address `json:"contractAddress,omitempty"`

	// decoded fields
	Deploy  bool            `json:"deploy,omitempty"`
	Method  *DecodedMethod  `json:"method,omitempty"`
	Events  []DecodedEvent  `json:"events"`
	UserOps []DecodedUserOp `json:"userOps,omitempty"`
}

type DecodedBlock struct {
	Number     *hexutil.Big   `json:"number"`
	Hash       common.Hash    `json:"hash"`
	ParentHash common.Hash    `json:"parentHash"`
	Timestamp  hexutil.Uint64 `json:"timestamp"`
	GasUsed    hexutil.Uint64 `json:"gasUsed"`
	GasLimit   hexutil.Uint64 `json:"gasLimit"`
	BaseFee    *hexutil.Big   `json:"baseFeePerGas,omitempty"`
	Miner      common.Address `json:"miner"`
	Txs        []DecodedTx    `json:"transactions"`
	Signer     types.Signer   `json:"-"`
}

var entryPoints = map[common.Address]bool{
	common.HexToAddress("0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789"): true,
	common.HexToAddress("0x0000000071727De22E5E9d8BAf0edAc6f37da032"): true,
}

// helper functions

func firstAbi(abis []*sigprovider.Abi) *sigprovider.Abi {
	if len(abis) > 0 {
		return abis[0]
	}
	return nil
}

func unpackArgs(args abi.Arguments, data []byte) (map[string]string, error) {
	values, err := args.Unpack(data)
	if err != nil {
		return nil, err
	}

	params := make(map[string]string, len(args))
	for i, arg := range args {
		params[arg.Name] = fmt.Sprintf("%v", values[i])
	}
	return params, nil
}

func splitEventInputs(event *abi.Event) (indexed, nonIndexed abi.Arguments) {
	for _, input := range event.Inputs {
		if input.Indexed {
			indexed = append(indexed, input)
		} else {
			nonIndexed = append(nonIndexed, input)
		}
	}
	return
}

func decodeTopics(indexed abi.Arguments, topics []common.Hash) map[string]string {
	params := make(map[string]string)
	for i, arg := range indexed {
		if i+1 < len(topics) {
			params[arg.Name] = topics[i+1].Hex()
		}
	}
	return params
}

func applyFallbackName(name *string, abis map[string]*sigprovider.Abi, key string) {
	if *name == "" && abis != nil {
		if sigAbi, ok := abis[key]; ok {
			*name = sigAbi.GetName()
		}
	}
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
		if err == nil {
			if abi := firstAbi(resp.GetAbi()); abi != nil {
				out[sel] = abi
			}
		}
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
			if abi := firstAbi(responses[i].GetAbi()); abi != nil {
				out[key] = abi
			}
		}
	}
	return out
}

func (f *Fx) lookupABI(ctx context.Context, addr common.Address) *abi.ABI {
	if f.ByteDB == nil {
		return nil
	}

	f.RLock()
	if cached, ok := f.abiCache[addr]; ok {
		f.RUnlock()
		return cached
	}
	f.RUnlock()

	resp, err := f.ByteDB.SearchSources(ctx, &bytecodedb.SearchSourcesRequest{
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
			dm.Params, _ = unpackArgs(method.Inputs, data[4:])
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

			indexed, nonIndexed := splitEventInputs(event)
			de.Params = decodeTopics(indexed, log.Topics)

			if len(log.Data) > 0 {
				if dataParams, err := unpackArgs(nonIndexed, log.Data); err == nil {
					maps.Copy(de.Params, dataParams)
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

func (f *Fx) decodeTx(t Transaction, contractABI *abi.ABI, funcAbis, eventAbis map[string]*sigprovider.Abi, ctx context.Context) DecodedTx {
	tx := t.Tx
	v, r, s := tx.RawSignatureValues()

	dt := DecodedTx{
		Hash:              tx.Hash(),
		Nonce:             tx.Nonce(),
		From:              t.From,
		To:                tx.To(),
		Value:             (*hexutil.Big)(tx.Value()),
		Input:             tx.Data(),
		Type:              tx.Type(),
		Gas:               tx.Gas(),
		GasPrice:          (*hexutil.Big)(tx.GasPrice()),
		ChainID:           (*hexutil.Big)(tx.ChainId()),
		V:                 (*hexutil.Big)(v),
		R:                 (*hexutil.Big)(r),
		S:                 (*hexutil.Big)(s),
		Status:            t.Receipt.Status,
		GasUsed:           t.Receipt.GasUsed,
		EffectiveGasPrice: (*hexutil.Big)(t.Receipt.EffectiveGasPrice),
		CumulativeGasUsed: t.Receipt.CumulativeGasUsed,
		Events:            make([]DecodedEvent, 0, len(t.Receipt.Logs)),
	}

	if tx.Type() >= 1 {
		dt.AccessList = tx.AccessList()
	}
	if tx.Type() >= 2 {
		dt.MaxFeePerGas = (*hexutil.Big)(tx.GasFeeCap())
		dt.MaxPriorityFee = (*hexutil.Big)(tx.GasTipCap())
	}
	if tx.Type() == 3 {
		dt.BlobGas = tx.BlobGas()
		dt.BlobGasFeeCap = (*hexutil.Big)(tx.BlobGasFeeCap())
		dt.BlobHashes = tx.BlobHashes()
	}

	if t.Receipt.ContractAddress != (common.Address{}) {
		dt.ContractAddress = &t.Receipt.ContractAddress
	}

	if tx.To() == nil {
		dt.Deploy = true
		return dt
	}

	dt.Method = f.decodeMethod(contractABI, tx.Data())
	if dt.Method != nil {
		applyFallbackName(&dt.Method.Name, funcAbis, dt.Method.Selector)
	}

	for _, log := range t.Receipt.Logs {
		logABI := contractABI
		if log.Address != *tx.To() {
			logABI = f.lookupABI(ctx, log.Address)
		}

		de := f.decodeEvent(logABI, log)
		if len(log.Topics) > 0 {
			applyFallbackName(&de.Name, eventAbis, log.Topics[0].Hex())
		}

		dt.Events = append(dt.Events, de)
	}

	if entryPoints[*tx.To()] {
		ops := f.lookupUserOps(ctx, tx.Hash())
		for j := range ops {
			if len(ops[j].CallData) >= 4 {
				opABI := f.lookupABI(ctx, ops[j].Sender)
				ops[j].Method = f.decodeMethod(opABI, ops[j].CallData)
			}
		}
		dt.UserOps = ops
	}

	return dt
}

func (f *Fx) Decode(ctx context.Context, block *Block) (*DecodedBlock, error) {
	selectors := f.collectSelectors(block)
	funcAbis := f.lookupFunctions(ctx, selectors)
	eventAbis := f.lookupEvents(ctx, block)
	decoded := make([]DecodedTx, len(block.Transactions))

	for i, t := range block.Transactions {
		var contractABI *abi.ABI
		if t.Tx.To() != nil {
			contractABI = f.lookupABI(ctx, *t.Tx.To())
		}
		decoded[i] = f.decodeTx(t, contractABI, funcAbis, eventAbis, ctx)
	}

	return &DecodedBlock{
		Number:     (*hexutil.Big)(block.Header.Number),
		Hash:       block.Header.Hash(),
		ParentHash: block.Header.ParentHash,
		Timestamp:  hexutil.Uint64(block.Header.Time),
		GasUsed:    hexutil.Uint64(block.Header.GasUsed),
		GasLimit:   hexutil.Uint64(block.Header.GasLimit),
		BaseFee:    (*hexutil.Big)(block.Header.BaseFee),
		Miner:      block.Header.Coinbase,
		Txs:        decoded,
		Signer:     block.Signer,
	}, nil
}
