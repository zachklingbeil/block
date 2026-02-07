package fx

// // output types

// type DecodedMethod struct {
// 	Selector  string            `json:"selector"`
// 	Name      string            `json:"name,omitempty"`
// 	Signature string            `json:"signature,omitempty"`
// 	Params    map[string]string `json:"params,omitempty"`
// }

// type DecodedEvent struct {
// 	Index     uint              `json:"index"`
// 	Address   string            `json:"address"`
// 	Topic     string            `json:"topic"`
// 	Name      string            `json:"name,omitempty"`
// 	Signature string            `json:"signature,omitempty"`
// 	Params    map[string]string `json:"params,omitempty"`
// }

// type DecodedUserOp struct {
// 	Sender   string         `json:"sender"`
// 	Nonce    string         `json:"nonce"`
// 	CallData string         `json:"callData,omitempty"`
// 	Method   *DecodedMethod `json:"method,omitempty"`
// }

// type DecodedTx struct {
// 	Hash              string           `json:"hash"`
// 	Nonce             uint64           `json:"nonce"`
// 	From              string           `json:"from"`
// 	To                string           `json:"to,omitempty"`
// 	Value             string           `json:"value"`
// 	ValueUnit         string           `json:"valueUnit"`
// 	Input             string           `json:"input,omitempty"`
// 	Type              string           `json:"type"`
// 	Gas               uint64           `json:"gas"`
// 	GasPrice          string           `json:"gasPrice,omitempty"`
// 	GasPriceUnit      string           `json:"gasPriceUnit,omitempty"`
// 	MaxFeePerGas      string           `json:"maxFeePerGas,omitempty"`
// 	MaxFeeUnit        string           `json:"maxFeeUnit,omitempty"`
// 	MaxPriorityFee    string           `json:"maxPriorityFeePerGas,omitempty"`
// 	MaxPriorityUnit   string           `json:"maxPriorityUnit,omitempty"`
// 	ChainID           string           `json:"chainId,omitempty"`
// 	AccessList        []AccessListItem `json:"accessList,omitempty"`
// 	BlobGas           uint64           `json:"blobGas,omitempty"`
// 	BlobGasFeeCap     string           `json:"maxFeePerBlobGas,omitempty"`
// 	BlobGasFeeUnit    string           `json:"blobGasFeeUnit,omitempty"`
// 	BlobHashes        []string         `json:"blobVersionedHashes,omitempty"`
// 	Status            string           `json:"status"`
// 	GasUsed           uint64           `json:"gasUsed"`
// 	EffectiveGasPrice string           `json:"effectiveGasPrice"`
// 	EffectiveGasUnit  string           `json:"effectiveGasUnit"`
// 	CumulativeGasUsed uint64           `json:"cumulativeGasUsed"`
// 	ContractAddress   string           `json:"contractAddress,omitempty"`
// 	Deploy            bool             `json:"deploy,omitempty"`
// 	Method            *DecodedMethod   `json:"method,omitempty"`
// 	Events            []DecodedEvent   `json:"events"`
// 	UserOps           []DecodedUserOp  `json:"userOps,omitempty"`
// }

// type AccessListItem struct {
// 	Address     string   `json:"address"`
// 	StorageKeys []string `json:"storageKeys"`
// }

// type DecodedBlock struct {
// 	Number      string      `json:"number"`
// 	Hash        string      `json:"hash"`
// 	ParentHash  string      `json:"parentHash"`
// 	Timestamp   string      `json:"timestamp"`
// 	GasUsed     uint64      `json:"gasUsed"`
// 	GasLimit    uint64      `json:"gasLimit"`
// 	BaseFee     string      `json:"baseFeePerGas,omitempty"`
// 	BaseFeeUnit string      `json:"baseFeeUnit,omitempty"`
// 	Miner       string      `json:"miner"`
// 	Txs         []DecodedTx `json:"transactions"`
// }

// var entryPoints = map[common.Address]bool{
// 	common.HexToAddress("0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789"): true,
// 	common.HexToAddress("0x0000000071727De22E5E9d8BAf0edAc6f37da032"): true,
// }

// // helper functions

// func formatTimestamp(ts uint64) string {
// 	return fmt.Sprintf("%d", time.Unix(int64(ts), 0).UnixMicro())
// }

// func txTypeToString(txType uint8) string {
// 	switch txType {
// 	case 0:
// 		return "Legacy"
// 	case 1:
// 		return "EIP-2930 (Access List)"
// 	case 2:
// 		return "EIP-1559 (Dynamic Fee)"
// 	case 3:
// 		return "EIP-4844 (Blob)"
// 	default:
// 		return fmt.Sprintf("Unknown (%d)", txType)
// 	}
// }

// func statusToString(status uint64) string {
// 	if status == 1 {
// 		return "Success"
// 	}
// 	return "Failed"
// }

// func formatBigInt(value *big.Int) string {
// 	if value == nil {
// 		return "0"
// 	}
// 	return value.String()
// }

// func determineUnit(value *big.Int) string {
// 	if value == nil || value.Sign() == 0 {
// 		return "wei"
// 	}

// 	// If >= 1 ETH, use ETH
// 	oneEth := new(big.Int).SetUint64(params.Ether)
// 	if value.Cmp(oneEth) >= 0 {
// 		return "wei" // Keep as wei for precision
// 	}

// 	// If >= 1 Gwei, use Gwei
// 	oneGwei := new(big.Int).SetUint64(params.GWei)
// 	if value.Cmp(oneGwei) >= 0 {
// 		return "wei"
// 	}

// 	return "wei"
// }

// func firstAbi(abis []*sigprovider.Abi) *sigprovider.Abi {
// 	if len(abis) > 0 {
// 		return abis[0]
// 	}
// 	return nil
// }

// func unpackArgs(args abi.Arguments, data []byte) (map[string]string, error) {
// 	values, err := args.Unpack(data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	params := make(map[string]string, len(args))
// 	for i, arg := range args {
// 		if i < len(values) {
// 			params[arg.Name] = fmt.Sprintf("%v", values[i])
// 		}
// 	}
// 	return params, nil
// }

// func splitEventInputs(event *abi.Event) (indexed, nonIndexed abi.Arguments) {
// 	for _, input := range event.Inputs {
// 		if input.Indexed {
// 			indexed = append(indexed, input)
// 		} else {
// 			nonIndexed = append(nonIndexed, input)
// 		}
// 	}
// 	return
// }

// func decodeTopics(indexed abi.Arguments, topics []common.Hash) map[string]string {
// 	params := make(map[string]string)
// 	for i, arg := range indexed {
// 		if i+1 < len(topics) {
// 			params[arg.Name] = topics[i+1].Hex()
// 		}
// 	}
// 	return params
// }

// func applyFallbackName(name *string, abis map[string]*sigprovider.Abi, key string) {
// 	if *name == "" && abis != nil {
// 		if sigAbi, ok := abis[key]; ok {
// 			*name = sigAbi.GetName()
// 		}
// 	}
// }

// // gRPC lookups

// func (f *Fx) lookupFunctions(ctx context.Context, selectors []string) map[string]*sigprovider.Abi {
// 	if len(selectors) == 0 || f.Sig == nil {
// 		return nil
// 	}

// 	out := make(map[string]*sigprovider.Abi)
// 	for _, sel := range selectors {
// 		resp, err := f.Sig.GetFunctionAbi(ctx, &sigprovider.GetFunctionAbiRequest{
// 			TxInput: sel,
// 		})
// 		if err == nil {
// 			if abi := firstAbi(resp.GetAbi()); abi != nil {
// 				out[sel] = abi
// 			}
// 		}
// 	}
// 	return out
// }

// func (f *Fx) lookupEvents(ctx context.Context, block *Block) map[string]*sigprovider.Abi {
// 	if f.Sig == nil {
// 		return nil
// 	}

// 	var reqs []*sigprovider.GetEventAbiRequest
// 	keys := make([]string, 0)
// 	seen := make(map[string]bool)

// 	for _, t := range block.Transactions {
// 		for _, log := range t.Receipt.Logs {
// 			if len(log.Topics) == 0 {
// 				continue
// 			}
// 			key := log.Topics[0].Hex()
// 			if seen[key] {
// 				continue
// 			}
// 			seen[key] = true
// 			keys = append(keys, key)

// 			topics := make([]string, len(log.Topics))
// 			for i, t := range log.Topics {
// 				topics[i] = t.Hex()
// 			}

// 			reqs = append(reqs, &sigprovider.GetEventAbiRequest{
// 				Data:   hexutil.Encode(log.Data),
// 				Topics: strings.Join(topics, ","),
// 			})
// 		}
// 	}

// 	if len(reqs) == 0 {
// 		return nil
// 	}

// 	resp, err := f.Sig.BatchGetEventAbis(ctx, &sigprovider.BatchGetEventAbisRequest{
// 		Requests: reqs,
// 	})
// 	if err != nil {
// 		return nil
// 	}

// 	out := make(map[string]*sigprovider.Abi)
// 	responses := resp.GetResponses()
// 	for i, key := range keys {
// 		if i < len(responses) {
// 			if abi := firstAbi(responses[i].GetAbi()); abi != nil {
// 				out[key] = abi
// 			}
// 		}
// 	}
// 	return out
// }

// func (f *Fx) lookupABI(ctx context.Context, addr common.Address) *abi.ABI {
// 	if f.ByteDB == nil {
// 		return nil
// 	}

// 	f.RLock()
// 	if cached, ok := f.abiCache[addr]; ok {
// 		f.RUnlock()
// 		return cached
// 	}
// 	f.RUnlock()

// 	resp, err := f.ByteDB.SearchSources(ctx, &bytecodedb.SearchSourcesRequest{
// 		Bytecode:     addr.Hex(),
// 		BytecodeType: bytecodedb.BytecodeType_DEPLOYED_BYTECODE,
// 	})
// 	if err != nil {
// 		return nil
// 	}

// 	sources := resp.GetSources()
// 	if len(sources) == 0 {
// 		return nil
// 	}

// 	raw := sources[0].GetAbi()
// 	if raw == "" {
// 		return nil
// 	}

// 	parsed, err := abi.JSON(strings.NewReader(raw))
// 	if err != nil {
// 		return nil
// 	}

// 	f.Lock()
// 	f.abiCache[addr] = &parsed
// 	f.Unlock()

// 	return &parsed
// }

// func (f *Fx) lookupUserOps(ctx context.Context, txHash common.Hash) []DecodedUserOp {
// 	if f.Ops == nil {
// 		return nil
// 	}

// 	hex := txHash.Hex()
// 	resp, err := f.Ops.ListUserOps(ctx, &userops.ListUserOpsRequest{
// 		TransactionHash: &hex,
// 	})
// 	if err != nil {
// 		return nil
// 	}

// 	items := resp.GetItems()
// 	if len(items) == 0 {
// 		return nil
// 	}

// 	ops := make([]DecodedUserOp, 0, len(items))
// 	for _, item := range items {
// 		full, err := f.Ops.GetUserOp(ctx, &userops.GetUserOpRequest{
// 			Hash: item.GetHash(),
// 		})
// 		if err != nil {
// 			continue
// 		}

// 		callData := full.GetCallData()
// 		op := DecodedUserOp{
// 			Sender: full.GetSender(),
// 			Nonce:  full.GetNonce(),
// 		}

// 		if callData != "" && callData != "0x" {
// 			op.CallData = callData
// 		}

// 		ops = append(ops, op)
// 	}
// 	return ops
// }

// // decode logic

// func (f *Fx) decodeMethod(contractABI *abi.ABI, data []byte) *DecodedMethod {
// 	if len(data) < 4 {
// 		return nil
// 	}

// 	selector := hexutil.Encode(data[:4])
// 	dm := &DecodedMethod{Selector: selector}

// 	if contractABI != nil {
// 		method, err := contractABI.MethodById(data[:4])
// 		if err == nil {
// 			dm.Name = method.Name
// 			dm.Signature = method.Sig
// 			dm.Params, _ = unpackArgs(method.Inputs, data[4:])
// 			return dm
// 		}
// 	}

// 	dm.Name = "Unknown"
// 	return dm
// }

// func (f *Fx) decodeEvent(contractABI *abi.ABI, log *types.Log) DecodedEvent {
// 	de := DecodedEvent{
// 		Index:   log.Index,
// 		Address: log.Address.Hex(),
// 	}

// 	if len(log.Topics) == 0 {
// 		return de
// 	}

// 	de.Topic = log.Topics[0].Hex()

// 	if contractABI != nil {
// 		event, err := contractABI.EventByID(log.Topics[0])
// 		if err == nil {
// 			de.Name = event.Name
// 			de.Signature = event.Sig

// 			indexed, nonIndexed := splitEventInputs(event)
// 			de.Params = decodeTopics(indexed, log.Topics)

// 			if len(log.Data) > 0 {
// 				if dataParams, err := unpackArgs(nonIndexed, log.Data); err == nil {
// 					maps.Copy(de.Params, dataParams)
// 				}
// 			}
// 			return de
// 		}
// 	}

// 	de.Name = "Unknown"
// 	return de
// }

// func (f *Fx) collectSelectors(block *Block) []string {
// 	funcSet := make(map[string]bool)
// 	for _, t := range block.Transactions {
// 		if t.Tx.To() != nil && len(t.Tx.Data()) >= 4 {
// 			funcSet[hexutil.Encode(t.Tx.Data()[:4])] = true
// 		}
// 	}

// 	funcs := make([]string, 0, len(funcSet))
// 	for k := range funcSet {
// 		funcs = append(funcs, k)
// 	}
// 	return funcs
// }

// func (f *Fx) decodeTx(t Transaction, contractABI *abi.ABI, funcAbis, eventAbis map[string]*sigprovider.Abi, ctx context.Context) DecodedTx {
// 	tx := t.Tx

// 	dt := DecodedTx{
// 		Hash:              tx.Hash().Hex(),
// 		Nonce:             tx.Nonce(),
// 		From:              t.From.Hex(),
// 		Value:             formatBigInt(tx.Value()),
// 		ValueUnit:         determineUnit(tx.Value()),
// 		Type:              txTypeToString(tx.Type()),
// 		Gas:               tx.Gas(),
// 		Status:            statusToString(t.Receipt.Status),
// 		GasUsed:           t.Receipt.GasUsed,
// 		EffectiveGasPrice: formatBigInt(t.Receipt.EffectiveGasPrice),
// 		EffectiveGasUnit:  determineUnit(t.Receipt.EffectiveGasPrice),
// 		CumulativeGasUsed: t.Receipt.CumulativeGasUsed,
// 		Events:            make([]DecodedEvent, 0, len(t.Receipt.Logs)),
// 	}

// 	if tx.To() != nil {
// 		dt.To = tx.To().Hex()
// 	}

// 	if tx.ChainId() != nil {
// 		dt.ChainID = tx.ChainId().String()
// 	}

// 	if tx.GasPrice() != nil {
// 		dt.GasPrice = formatBigInt(tx.GasPrice())
// 		dt.GasPriceUnit = determineUnit(tx.GasPrice())
// 	}

// 	if tx.Type() >= 1 {
// 		accessList := tx.AccessList()
// 		dt.AccessList = make([]AccessListItem, len(accessList))
// 		for i, item := range accessList {
// 			keys := make([]string, len(item.StorageKeys))
// 			for j, key := range item.StorageKeys {
// 				keys[j] = key.Hex()
// 			}
// 			dt.AccessList[i] = AccessListItem{
// 				Address:     item.Address.Hex(),
// 				StorageKeys: keys,
// 			}
// 		}
// 	}

// 	if tx.Type() >= 2 {
// 		if tx.GasFeeCap() != nil {
// 			dt.MaxFeePerGas = formatBigInt(tx.GasFeeCap())
// 			dt.MaxFeeUnit = determineUnit(tx.GasFeeCap())
// 		}
// 		if tx.GasTipCap() != nil {
// 			dt.MaxPriorityFee = formatBigInt(tx.GasTipCap())
// 			dt.MaxPriorityUnit = determineUnit(tx.GasTipCap())
// 		}
// 	}

// 	if tx.Type() == 3 {
// 		dt.BlobGas = tx.BlobGas()
// 		if tx.BlobGasFeeCap() != nil {
// 			dt.BlobGasFeeCap = formatBigInt(tx.BlobGasFeeCap())
// 			dt.BlobGasFeeUnit = determineUnit(tx.BlobGasFeeCap())
// 		}
// 		hashes := tx.BlobHashes()
// 		dt.BlobHashes = make([]string, len(hashes))
// 		for i, h := range hashes {
// 			dt.BlobHashes[i] = h.Hex()
// 		}
// 	}

// 	if t.Receipt.ContractAddress != (common.Address{}) {
// 		dt.ContractAddress = t.Receipt.ContractAddress.Hex()
// 		dt.Deploy = true
// 		return dt
// 	}

// 	txData := tx.Data()
// 	if len(txData) > 0 {
// 		dt.Method = f.decodeMethod(contractABI, txData)
// 		if dt.Method != nil {
// 			applyFallbackName(&dt.Method.Name, funcAbis, dt.Method.Selector)
// 			if dt.Method.Name == "Unknown" {
// 				dt.Input = hexutil.Encode(txData)
// 			}
// 		}
// 	}

// 	for _, log := range t.Receipt.Logs {
// 		logABI := contractABI
// 		if tx.To() != nil && log.Address != *tx.To() {
// 			logABI = f.lookupABI(ctx, log.Address)
// 		}

// 		de := f.decodeEvent(logABI, log)
// 		if len(log.Topics) > 0 {
// 			applyFallbackName(&de.Name, eventAbis, log.Topics[0].Hex())
// 		}

// 		dt.Events = append(dt.Events, de)
// 	}

// 	if tx.To() != nil && entryPoints[*tx.To()] {
// 		ops := f.lookupUserOps(ctx, tx.Hash())
// 		for j := range ops {
// 			if ops[j].CallData != "" {
// 				callData := common.FromHex(ops[j].CallData)
// 				if len(callData) >= 4 {
// 					opABI := f.lookupABI(ctx, common.HexToAddress(ops[j].Sender))
// 					ops[j].Method = f.decodeMethod(opABI, callData)
// 				}
// 			}
// 		}
// 		dt.UserOps = ops
// 	}

// 	return dt
// }

// func (f *Fx) Decode(ctx context.Context, block *Block) (*DecodedBlock, error) {
// 	selectors := f.collectSelectors(block)
// 	funcAbis := f.lookupFunctions(ctx, selectors)
// 	eventAbis := f.lookupEvents(ctx, block)
// 	decoded := make([]DecodedTx, len(block.Transactions))

// 	for i, t := range block.Transactions {
// 		var contractABI *abi.ABI
// 		if t.Tx.To() != nil {
// 			contractABI = f.lookupABI(ctx, *t.Tx.To())
// 		}
// 		decoded[i] = f.decodeTx(t, contractABI, funcAbis, eventAbis, ctx)
// 	}

// 	db := &DecodedBlock{
// 		Number:     block.Header.Number.String(),
// 		Hash:       block.Header.Hash().Hex(),
// 		ParentHash: block.Header.ParentHash.Hex(),
// 		Timestamp:  formatTimestamp(block.Header.Time),
// 		GasUsed:    block.Header.GasUsed,
// 		GasLimit:   block.Header.GasLimit,
// 		Miner:      block.Header.Coinbase.Hex(),
// 		Txs:        decoded,
// 	}

// 	if block.Header.BaseFee != nil {
// 		db.BaseFee = formatBigInt(block.Header.BaseFee)
// 		db.BaseFeeUnit = determineUnit(block.Header.BaseFee)
// 	}

// 	return db, nil
// }
