package fx

// output types

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

// type TokenTransfer struct {
// 	Token string `json:"token"`
// 	From  string `json:"from"`
// 	To    string `json:"to"`
// 	Value string `json:"value"`
// }

// type TokenApproval struct {
// 	Token   string `json:"token"`
// 	Owner   string `json:"owner"`
// 	Spender string `json:"spender"`
// 	Value   string `json:"value"`
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
// 	Input             string           `json:"input,omitempty"`
// 	Type              string           `json:"type"`
// 	Gas               uint64           `json:"gas"`
// 	GasPrice          string           `json:"gasPrice,omitempty"`
// 	MaxFeePerGas      string           `json:"maxFeePerGas,omitempty"`
// 	MaxPriorityFee    string           `json:"maxPriorityFeePerGas,omitempty"`
// 	ChainID           string           `json:"chainId,omitempty"`
// 	AccessList        []AccessListItem `json:"accessList,omitempty"`
// 	BlobGas           uint64           `json:"blobGas,omitempty"`
// 	BlobGasFeeCap     string           `json:"maxFeePerBlobGas,omitempty"`
// 	BlobHashes        []string         `json:"blobVersionedHashes,omitempty"`
// 	Status            string           `json:"status"`
// 	GasUsed           uint64           `json:"gasUsed"`
// 	EffectiveGasPrice string           `json:"effectiveGasPrice"`
// 	CumulativeGasUsed uint64           `json:"cumulativeGasUsed"`
// 	ContractAddress   string           `json:"contractAddress,omitempty"`
// 	Deploy            bool             `json:"deploy,omitempty"`
// 	Method            *DecodedMethod   `json:"method,omitempty"`
// 	Events            []DecodedEvent   `json:"events"`
// 	Transfers         []TokenTransfer  `json:"transfers,omitempty"`
// 	Approvals         []TokenApproval  `json:"approvals,omitempty"`
// 	UserOps           []DecodedUserOp  `json:"userOps,omitempty"`
// }

// type AccessListItem struct {
// 	Address     string   `json:"address"`
// 	StorageKeys []string `json:"storageKeys"`
// }

// type DecodedBlock struct {
// 	Number     string      `json:"number"`
// 	Hash       string      `json:"hash"`
// 	ParentHash string      `json:"parentHash"`
// 	Timestamp  string      `json:"timestamp"`
// 	GasUsed    uint64      `json:"gasUsed"`
// 	GasLimit   uint64      `json:"gasLimit"`
// 	BaseFee    string      `json:"baseFeePerGas,omitempty"`
// 	Miner      string      `json:"miner"`
// 	Txs        []DecodedTx `json:"transactions"`
// }

// var entryPoints = map[common.Address]bool{
// 	common.HexToAddress("0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789"): true,
// 	common.HexToAddress("0x0000000071727De22E5E9d8BAf0edAc6f37da032"): true,
// }

// // helpers

// func formatTimestamp(ts uint64) string {
// 	return fmt.Sprintf("%d", time.Unix(int64(ts), 0).Unix())
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

// func topicToAddress(topic common.Hash) string {
// 	return common.BytesToAddress(topic.Bytes()).Hex()
// }

// func topicToAmount(data []byte) string {
// 	if len(data) < 32 {
// 		return "0"
// 	}
// 	return new(big.Int).SetBytes(data[:32]).String()
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

// func extractTransfer(log *types.Log) *TokenTransfer {
// 	if len(log.Topics) != 3 || log.Topics[0] != TopicTransfer {
// 		return nil
// 	}
// 	return &TokenTransfer{
// 		Token: log.Address.Hex(),
// 		From:  topicToAddress(log.Topics[1]),
// 		To:    topicToAddress(log.Topics[2]),
// 		Value: topicToAmount(log.Data),
// 	}
// }

// func extractApproval(log *types.Log) *TokenApproval {
// 	if len(log.Topics) != 3 || log.Topics[0] != TopicApproval {
// 		return nil
// 	}
// 	return &TokenApproval{
// 		Token:   log.Address.Hex(),
// 		Owner:   topicToAddress(log.Topics[1]),
// 		Spender: topicToAddress(log.Topics[2]),
// 		Value:   topicToAmount(log.Data),
// 	}
// }

// // resolve — populate caches from services before decoding

// func (f *Fx) resolveABIs(ctx context.Context, block *Block) {
// 	if f.ByteDB == nil {
// 		return
// 	}

// 	seen := make(map[common.Address]bool)
// 	var addrs []common.Address

// 	for _, t := range block.Transactions {
// 		if t.Tx.To() != nil {
// 			addrs = append(addrs, *t.Tx.To())
// 		}
// 		for _, log := range t.Receipt.Logs {
// 			addrs = append(addrs, log.Address)
// 		}
// 	}

// 	for _, addr := range addrs {
// 		if seen[addr] {
// 			continue
// 		}
// 		seen[addr] = true

// 		f.RLock()
// 		_, cached := f.abis[addr]
// 		f.RUnlock()
// 		if cached {
// 			continue
// 		}

// 		resp, err := f.ByteDB.SearchSources(ctx, &bytecodedb.SearchSourcesRequest{
// 			Bytecode:     addr.Hex(),
// 			BytecodeType: bytecodedb.BytecodeType_DEPLOYED_BYTECODE,
// 		})
// 		if err != nil {
// 			continue
// 		}

// 		sources := resp.GetSources()
// 		if len(sources) == 0 {
// 			continue
// 		}

// 		raw := sources[0].GetAbi()
// 		if raw == "" {
// 			continue
// 		}

// 		parsed, err := abi.JSON(strings.NewReader(raw))
// 		if err != nil {
// 			continue
// 		}

// 		f.Lock()
// 		f.abis[addr] = &parsed
// 		for _, event := range parsed.Events {
// 			e := event
// 			if _, exists := f.events[e.ID]; !exists {
// 				f.events[e.ID] = &e
// 			}
// 		}
// 		for _, method := range parsed.Methods {
// 			m := method
// 			sel := hexutil.Encode(m.ID)
// 			if _, exists := f.methods[sel]; !exists {
// 				f.methods[sel] = &m
// 			}
// 		}
// 		f.Unlock()
// 	}
// }

// func (f *Fx) resolveMethods(ctx context.Context, block *Block) {
// 	if f.Sig == nil {
// 		return
// 	}

// 	seen := make(map[string]bool)

// 	for _, t := range block.Transactions {
// 		data := t.Tx.Data()
// 		if t.Tx.To() == nil || len(data) < 4 {
// 			continue
// 		}
// 		sel := hexutil.Encode(data[:4])
// 		if seen[sel] {
// 			continue
// 		}
// 		seen[sel] = true

// 		f.RLock()
// 		_, cached := f.methods[sel]
// 		f.RUnlock()
// 		if cached {
// 			continue
// 		}

// 		resp, err := f.Sig.GetFunctionAbi(ctx, &sigprovider.GetFunctionAbiRequest{
// 			TxInput: sel,
// 		})
// 		if err != nil {
// 			continue
// 		}

// 		sigAbi := firstAbi(resp.GetAbi())
// 		if sigAbi == nil {
// 			continue
// 		}

// 		args := make(abi.Arguments, 0, len(sigAbi.GetInputs()))
// 		var sigParts []string
// 		for _, input := range sigAbi.GetInputs() {
// 			argType, err := abi.NewType(input.GetType(), "", nil)
// 			if err != nil {
// 				continue
// 			}
// 			args = append(args, abi.Argument{
// 				Name: input.GetName(),
// 				Type: argType,
// 			})
// 			sigParts = append(sigParts, input.GetType())
// 		}

// 		name := sigAbi.GetName()
// 		sig := fmt.Sprintf("%s(%s)", name, strings.Join(sigParts, ","))
// 		method := abi.NewMethod(name, name, abi.Function, "", false, false, args, nil)
// 		method.Sig = sig

// 		f.Lock()
// 		f.methods[sel] = &method
// 		f.Unlock()
// 	}
// }

// func (f *Fx) resolveEvents(ctx context.Context, block *Block) {
// 	if f.Sig == nil {
// 		return
// 	}

// 	type pending struct {
// 		topic common.Hash
// 	}

// 	var reqs []*sigprovider.GetEventAbiRequest
// 	var items []pending
// 	seen := make(map[common.Hash]bool)

// 	for _, t := range block.Transactions {
// 		for _, log := range t.Receipt.Logs {
// 			if len(log.Topics) == 0 {
// 				continue
// 			}
// 			topic := log.Topics[0]
// 			if seen[topic] {
// 				continue
// 			}
// 			seen[topic] = true

// 			f.RLock()
// 			_, cached := f.events[topic]
// 			f.RUnlock()
// 			if cached {
// 				continue
// 			}

// 			topics := make([]string, len(log.Topics))
// 			for i, t := range log.Topics {
// 				topics[i] = t.Hex()
// 			}

// 			reqs = append(reqs, &sigprovider.GetEventAbiRequest{
// 				Data:   hexutil.Encode(log.Data),
// 				Topics: strings.Join(topics, ","),
// 			})
// 			items = append(items, pending{topic: topic})
// 		}
// 	}

// 	if len(reqs) == 0 {
// 		return
// 	}

// 	resp, err := f.Sig.BatchGetEventAbis(ctx, &sigprovider.BatchGetEventAbisRequest{
// 		Requests: reqs,
// 	})
// 	if err != nil {
// 		return
// 	}

// 	responses := resp.GetResponses()
// 	f.Lock()
// 	for i, p := range items {
// 		if i >= len(responses) {
// 			break
// 		}
// 		sigAbi := firstAbi(responses[i].GetAbi())
// 		if sigAbi == nil {
// 			continue
// 		}

// 		args := make(abi.Arguments, 0, len(sigAbi.GetInputs()))
// 		for _, input := range sigAbi.GetInputs() {
// 			argType, err := abi.NewType(input.GetType(), "", nil)
// 			if err != nil {
// 				continue
// 			}
// 			args = append(args, abi.Argument{
// 				Name:    input.GetName(),
// 				Type:    argType,
// 				Indexed: input.GetIndexed(),
// 			})
// 		}

// 		name := sigAbi.GetName()
// 		event := abi.NewEvent(name, name, false, args)
// 		f.events[p.topic] = &event
// 	}
// 	f.Unlock()
// }

// // decode — read from caches only

// func (f *Fx) decodeMethod(data []byte) *DecodedMethod {
// 	if len(data) < 4 {
// 		return nil
// 	}

// 	sel := hexutil.Encode(data[:4])
// 	dm := &DecodedMethod{Selector: sel}

// 	f.RLock()
// 	method, ok := f.methods[sel]
// 	f.RUnlock()

// 	if !ok {
// 		dm.Name = "Unknown"
// 		return dm
// 	}

// 	dm.Name = method.Name
// 	dm.Signature = method.Sig
// 	if len(data) > 4 {
// 		dm.Params, _ = unpackArgs(method.Inputs, data[4:])
// 	}
// 	return dm
// }

// func (f *Fx) decodeEvent(log *types.Log) DecodedEvent {
// 	de := DecodedEvent{
// 		Index:   log.Index,
// 		Address: log.Address.Hex(),
// 	}

// 	if len(log.Topics) == 0 {
// 		return de
// 	}

// 	de.Topic = log.Topics[0].Hex()

// 	f.RLock()
// 	event, ok := f.events[log.Topics[0]]
// 	f.RUnlock()

// 	if !ok {
// 		de.Name = "Unknown"
// 		return de
// 	}

// 	de.Name = event.Name
// 	de.Signature = event.Sig

// 	indexed, nonIndexed := splitEventInputs(event)
// 	de.Params = decodeTopics(indexed, log.Topics)

// 	if len(log.Data) > 0 {
// 		if dataParams, err := unpackArgs(nonIndexed, log.Data); err == nil {
// 			maps.Copy(de.Params, dataParams)
// 		}
// 	}
// 	return de
// }

// func (f *Fx) decodeUserOps(ctx context.Context, txHash common.Hash) []DecodedUserOp {
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

// 		op := DecodedUserOp{
// 			Sender: full.GetSender(),
// 			Nonce:  full.GetNonce(),
// 		}

// 		callData := full.GetCallData()
// 		if callData != "" && callData != "0x" {
// 			op.CallData = callData
// 			cd := common.FromHex(callData)
// 			if len(cd) >= 4 {
// 				op.Method = f.decodeMethod(cd)
// 			}
// 		}

// 		ops = append(ops, op)
// 	}
// 	return ops
// }

// func (f *Fx) decodeTx(ctx context.Context, t Transaction) DecodedTx {
// 	tx := t.Tx

// 	dt := DecodedTx{
// 		Hash:              tx.Hash().Hex(),
// 		Nonce:             tx.Nonce(),
// 		From:              t.From.Hex(),
// 		Value:             formatBigInt(tx.Value()),
// 		Type:              txTypeToString(tx.Type()),
// 		Gas:               tx.Gas(),
// 		Status:            statusToString(t.Receipt.Status),
// 		GasUsed:           t.Receipt.GasUsed,
// 		EffectiveGasPrice: formatBigInt(t.Receipt.EffectiveGasPrice),
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
// 		}
// 		if tx.GasTipCap() != nil {
// 			dt.MaxPriorityFee = formatBigInt(tx.GasTipCap())
// 		}
// 	}

// 	if tx.Type() == 3 {
// 		dt.BlobGas = tx.BlobGas()
// 		if tx.BlobGasFeeCap() != nil {
// 			dt.BlobGasFeeCap = formatBigInt(tx.BlobGasFeeCap())
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
// 		dt.Method = f.decodeMethod(txData)
// 		if dt.Method != nil && dt.Method.Name == "Unknown" {
// 			dt.Input = hexutil.Encode(txData)
// 		}
// 	}

// 	for _, log := range t.Receipt.Logs {
// 		if transfer := extractTransfer(log); transfer != nil {
// 			dt.Transfers = append(dt.Transfers, *transfer)
// 		}
// 		if approval := extractApproval(log); approval != nil {
// 			dt.Approvals = append(dt.Approvals, *approval)
// 		}
// 		dt.Events = append(dt.Events, f.decodeEvent(log))
// 	}

// 	if tx.To() != nil && entryPoints[*tx.To()] {
// 		dt.UserOps = f.decodeUserOps(ctx, tx.Hash())
// 	}

// 	return dt
// }

// // Decode resolves everything upfront, then decodes.
// func (f *Fx) Decode(ctx context.Context, block *Block) (*DecodedBlock, error) {
// 	// Phase 1: resolve — populate all caches from services
// 	f.resolveABIs(ctx, block)
// 	f.resolveMethods(ctx, block)
// 	f.resolveEvents(ctx, block)

// 	// Phase 2: decode — read from caches only
// 	decoded := make([]DecodedTx, len(block.Transactions))
// 	for i, t := range block.Transactions {
// 		decoded[i] = f.decodeTx(ctx, t)
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
// 	}

// 	return db, nil
// }
