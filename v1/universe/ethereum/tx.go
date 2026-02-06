package ethereum

// func (e *Ethereum) processBlock(ctx context.Context, block *types.Block) *Block {
// 	signer := e.Signer(block.Number(), block.Time())
// 	txs := block.Transactions()
// 	transactions := make([]*Transactions, 0, len(txs))
// 	for _, tx := range txs {
// 		// Fetch the transaction hash
// 		txHash := tx.Hash()
// 		// Fetch the full transaction data using the hash (if needed)
// 		fullTx, _, err := e.Factory.Eth.TransactionByHash(ctx, txHash)
// 		if err != nil {
// 			// fallback to using the tx from block if not found
// 			fullTx = tx
// 		}
// 		if txInfo := e.processTransaction(ctx, fullTx, signer); txInfo != nil {
// 			transactions = append(transactions, txInfo)
// 		}
// 	}
// 	return &Block{
// 		Number: block.NumberU64(),
// 	}
// }

// func (e *Ethereum) processTransaction(ctx context.Context, tx *types.Transaction, signer types.Signer) *Transactions {
// 	txInfo := &Transactions{
// 		Value:    tx.Value(),
// 		Gas:      tx.Gas(),
// 		GasPrice: tx.GasPrice(),
// 		Nonce:    tx.Nonce(),
// 	}

// 	// Set From address (prefer ENS/Token/Address)
// 	if addr, err := types.Sender(signer, tx); err == nil {
// 		txInfo.From = e.Who(addr.Hex())
// 	}

// 	// Set To address or contract creation (prefer ENS/Token/Address)
// 	if to := tx.To(); to == nil {
// 		txInfo.To = "Contract Creation"
// 	} else {
// 		txInfo.To = e.Who(to.Hex())
// 	}

// 	// Populate receipt info (logs, cumulative gas used, etc.)
// 	if receipt, err := e.Factory.Eth.TransactionReceipt(ctx, tx.Hash()); err == nil {
// 		e.populateReceiptInfo(txInfo, receipt)
// 	}
// 	return txInfo
// }

// func (e *Ethereum) populateReceiptInfo(txInfo *Transactions, receipt *types.Receipt) {
// 	txInfo.CumulativeGasUsed = receipt.CumulativeGasUsed
// 	for _, log := range receipt.Logs {
// 		if logInfo := e.ParseLogSafe(log); logInfo != nil {
// 			txInfo.Logs = append(txInfo.Logs, logInfo)
// 		}
// 	}
// }

// // Who returns the ENS, Token, or Address for a given hex address.
// func (e *Ethereum) Who(hex string) string {
// 	one := e.Zero.Source(hex)
// 	if one == nil {
// 		return hex
// 	}
// 	if one.ENS != "" && one.ENS != "." {
// 		return one.ENS
// 	}
// 	if one.Token != "" {
// 		return one.Token
// 	}
// 	return one.Address
// }
