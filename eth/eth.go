package eth

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// DecodedTransaction represents a fully decoded transaction
type DecodedTransaction struct {
	Hash             common.Hash      `json:"hash"`
	Nonce            uint64           `json:"nonce"`
	BlockHash        *common.Hash     `json:"blockHash"`
	BlockNumber      *big.Int         `json:"blockNumber"`
	TransactionIndex uint             `json:"transactionIndex"`
	From             common.Address   `json:"from"`
	To               *common.Address  `json:"to"`
	Value            *big.Int         `json:"value"`
	Gas              uint64           `json:"gas"`
	GasPrice         *big.Int         `json:"gasPrice"`
	GasFeeCap        *big.Int         `json:"maxFeePerGas,omitempty"`
	GasTipCap        *big.Int         `json:"maxPriorityFeePerGas,omitempty"`
	Data             []byte           `json:"input"`
	Type             uint8            `json:"type"`
	ChainID          *big.Int         `json:"chainId,omitempty"`
	AccessList       types.AccessList `json:"accessList,omitempty"`
	V                *big.Int         `json:"v"`
	R                *big.Int         `json:"r"`
	S                *big.Int         `json:"s"`
	Size             uint64           `json:"size"`
}

// DecodedReceipt represents a fully decoded transaction receipt
type DecodedReceipt struct {
	TxHash            common.Hash     `json:"transactionHash"`
	TxIndex           uint            `json:"transactionIndex"`
	From              common.Address  `json:"from"`
	To                *common.Address `json:"to"`
	CumulativeGasUsed uint64          `json:"cumulativeGasUsed"`
	EffectiveGasPrice *big.Int        `json:"effectiveGasPrice"`
	GasUsed           uint64          `json:"gasUsed"`
	ContractAddress   *common.Address `json:"contractAddress,omitempty"`
	Logs              []*types.Log    `json:"logs"`
	LogsBloom         types.Bloom     `json:"logsBloom"`
	Type              uint8           `json:"type"`
	Status            uint64          `json:"status"`
	BlobGasUsed       uint64          `json:"blobGasUsed,omitempty"`
	BlobGasPrice      *big.Int        `json:"blobGasPrice,omitempty"`
}

// DecodedBlock represents a fully decoded Ethereum block
type DecodedBlock struct {
	Number           *big.Int             `json:"number"`
	Hash             common.Hash          `json:"hash"`
	ParentHash       common.Hash          `json:"parentHash"`
	Nonce            types.BlockNonce     `json:"nonce"`
	SHA3Uncles       common.Hash          `json:"sha3Uncles"`
	LogsBloom        types.Bloom          `json:"logsBloom"`
	TransactionsRoot common.Hash          `json:"transactionsRoot"`
	StateRoot        common.Hash          `json:"stateRoot"`
	ReceiptsRoot     common.Hash          `json:"receiptsRoot"`
	Miner            common.Address       `json:"miner"`
	Difficulty       *big.Int             `json:"difficulty"`
	TotalDifficulty  *big.Int             `json:"totalDifficulty"`
	ExtraData        []byte               `json:"extraData"`
	Size             uint64               `json:"size"`
	GasLimit         uint64               `json:"gasLimit"`
	GasUsed          uint64               `json:"gasUsed"`
	Timestamp        uint64               `json:"timestamp"`
	Transactions     []DecodedTransaction `json:"transactions"`
	Receipts         []DecodedReceipt     `json:"receipts"`
	Uncles           []common.Hash        `json:"uncles"`
	MixHash          common.Hash          `json:"mixHash"`
	BaseFee          *big.Int             `json:"baseFeePerGas,omitempty"`
	WithdrawalsRoot  *common.Hash         `json:"withdrawalsRoot,omitempty"`
	BlobGasUsed      *uint64              `json:"blobGasUsed,omitempty"`
	ExcessBlobGas    *uint64              `json:"excessBlobGas,omitempty"`
	ParentBeaconRoot *common.Hash         `json:"parentBeaconBlockRoot,omitempty"`
}

// DecodeOptions provides configuration for block decoding
type DecodeOptions struct {
	StartBlock       *big.Int // If nil, starts from latest
	EndBlock         *big.Int // If nil, goes to genesis (0)
	ProgressCallback func(blockNum *big.Int, total *big.Int)
	ErrorCallback    func(blockNum *big.Int, err error)
}

// Close closes the connection and cancels ongoing operations
func (e *Ethereum) Close() {
	e.cancel()
	e.Factory.Eth.Close()
}

// GetLatestBlockNumber returns the latest block number
func (e *Ethereum) GetLatestBlockNumber() (*big.Int, error) {
	header, err := e.Factory.Eth.HeaderByNumber(e.Factory.Ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest header: %w", err)
	}
	return header.Number, nil
}

// DecodeBlock decodes a single block with all its transactions and receipts
func (e *Ethereum) DecodeBlock(blockNum *big.Int, opts DecodeOptions) (*DecodedBlock, error) {
	// Get the block
	block, err := e.Factory.Eth.BlockByNumber(e.Factory.Ctx, blockNum)
	if err != nil {
		return nil, fmt.Errorf("failed to get block %s: %w", blockNum.String(), err)
	}

	decodedBlock := &DecodedBlock{
		Number:           block.Number(),
		Hash:             block.Hash(),
		ParentHash:       block.ParentHash(),
		SHA3Uncles:       block.UncleHash(),
		LogsBloom:        block.Bloom(),
		TransactionsRoot: block.TxHash(),
		StateRoot:        block.Root(),
		ReceiptsRoot:     block.ReceiptHash(),
		Miner:            block.Coinbase(),
		Difficulty:       block.Difficulty(),
		ExtraData:        block.Extra(),
		Size:             block.Size(),
		GasLimit:         block.GasLimit(),
		GasUsed:          block.GasUsed(),
		Timestamp:        block.Time(),
		MixHash:          block.MixDigest(),
		BaseFee:          block.BaseFee(),
		BlobGasUsed:      block.BlobGasUsed(),
		ExcessBlobGas:    block.ExcessBlobGas(),
		ParentBeaconRoot: block.BeaconRoot(),
	}

	// Decode transactions
	transactions := block.Transactions()
	decodedBlock.Transactions = make([]DecodedTransaction, len(transactions))
	for i, tx := range transactions {
		decodedTx, err := e.decodeTransaction(tx, block)
		if err != nil {
			return nil, fmt.Errorf("failed to decode transaction %d in block %s: %w", i, blockNum.String(), err)
		}
		decodedBlock.Transactions[i] = *decodedTx
	}

	// Always include receipts
	if len(transactions) > 0 {
		receipts, err := e.getBlockReceipts(blockNum, block.Time(), transactions)
		if err != nil {
			return nil, fmt.Errorf("failed to get receipts for block %s: %w", blockNum.String(), err)
		}
		decodedBlock.Receipts = receipts
	}

	// Always include uncles
	uncles := block.Uncles()
	decodedBlock.Uncles = make([]common.Hash, len(uncles))
	for i, uncle := range uncles {
		decodedBlock.Uncles[i] = uncle.Hash()
	}

	return decodedBlock, nil
}

// DecodeBlockRange decodes a range of blocks from start to end (inclusive) and stores each block in Redis as it is decoded.
func (e *Ethereum) DecodeBlockRange(opts DecodeOptions) error {
	var startBlock, endBlock *big.Int
	var err error

	// Determine start block
	if opts.StartBlock == nil {
		startBlock, err = e.GetLatestBlockNumber()
		if err != nil {
			return fmt.Errorf("failed to get latest block number: %w", err)
		}
	} else {
		startBlock = new(big.Int).Set(opts.StartBlock)
	}

	// Determine end block
	if opts.EndBlock == nil {
		endBlock = big.NewInt(0) // Genesis
	} else {
		endBlock = new(big.Int).Set(opts.EndBlock)
	}

	// Validate range
	if startBlock.Cmp(endBlock) < 0 {
		return fmt.Errorf("start block (%s) must be >= end block (%s)", startBlock.String(), endBlock.String())
	}

	// Calculate total blocks to process
	totalBlocks := new(big.Int).Sub(startBlock, endBlock)
	totalBlocks.Add(totalBlocks, big.NewInt(1))

	current := new(big.Int).Set(startBlock)

	for current.Cmp(endBlock) >= 0 {
		select {
		case <-e.Factory.Ctx.Done():
			return fmt.Errorf("operation cancelled")
		default:
		}

		block, err := e.DecodeBlock(current, opts)
		if err != nil {
			if opts.ErrorCallback != nil {
				opts.ErrorCallback(current, err)
			}
			log.Printf("Error decoding block %s: %v", current.String(), err)
		} else {
			if err := e.StoreDecodedBlock(block); err != nil {
				log.Printf("Error storing block %s: %v", current.String(), err)
			}
		}

		if opts.ProgressCallback != nil {
			processed := new(big.Int).Sub(startBlock, current)
			processed.Add(processed, big.NewInt(1))
			opts.ProgressCallback(current, totalBlocks)
		}

		current.Sub(current, big.NewInt(1))
	}

	return nil
}

// decodeTransaction converts a types.Transaction to DecodedTransaction
func (e *Ethereum) decodeTransaction(tx *types.Transaction, block *types.Block) (*DecodedTransaction, error) {
	// Get sender address
	signer := e.Signer(block.Number(), block.Time())
	from, err := types.Sender(signer, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender: %w", err)
	}

	// Extract signature values
	v, r, s := tx.RawSignatureValues()

	blockHash := block.Hash()
	decoded := &DecodedTransaction{
		Hash:        tx.Hash(),
		Nonce:       tx.Nonce(),
		BlockHash:   &blockHash,
		BlockNumber: block.Number(),
		From:        from,
		To:          tx.To(),
		Value:       tx.Value(),
		Gas:         tx.Gas(),
		Data:        tx.Data(),
		Type:        tx.Type(),
		V:           v,
		R:           r,
		S:           s,
		Size:        tx.Size(),
	}

	// Set transaction index
	for i, blockTx := range block.Transactions() {
		if blockTx.Hash() == tx.Hash() {
			decoded.TransactionIndex = uint(i)
			break
		}
	}

	// Handle different transaction types
	switch tx.Type() {
	case types.LegacyTxType:
		decoded.GasPrice = tx.GasPrice()
		decoded.ChainID = tx.ChainId()
	case types.AccessListTxType:
		decoded.GasPrice = tx.GasPrice()
		decoded.ChainID = tx.ChainId()
		decoded.AccessList = tx.AccessList()
	case types.DynamicFeeTxType:
		decoded.GasFeeCap = tx.GasFeeCap()
		decoded.GasTipCap = tx.GasTipCap()
		decoded.ChainID = tx.ChainId()
		decoded.AccessList = tx.AccessList()
	case types.BlobTxType:
		decoded.GasFeeCap = tx.GasFeeCap()
		decoded.GasTipCap = tx.GasTipCap()
		decoded.ChainID = tx.ChainId()
		decoded.AccessList = tx.AccessList()
	}

	return decoded, nil
}

// Signer returns a signer for Ethereum mainnet at the given block number and time.
func (e *Ethereum) Signer(blockNumber *big.Int, blockTime uint64) types.Signer {
	return types.MakeSigner(e.Chain, blockNumber, blockTime)
}

// getBlockReceipts gets all receipts for transactions in a block
func (e *Ethereum) getBlockReceipts(blockNum *big.Int, blockTime uint64, transactions types.Transactions) ([]DecodedReceipt, error) {
	receipts := make([]DecodedReceipt, len(transactions))

	for i, tx := range transactions {
		receipt, err := e.Factory.Eth.TransactionReceipt(e.Factory.Ctx, tx.Hash())
		if err != nil {
			return nil, fmt.Errorf("failed to get receipt for tx %s: %w", tx.Hash().Hex(), err)
		}
		// Use the chain's configured signer method
		signer := e.Signer(blockNum, blockTime)
		from, err := types.Sender(signer, tx)
		if err != nil {
			return nil, fmt.Errorf("failed to get sender for receipt: %w", err)
		}
		decodedReceipt := DecodedReceipt{
			TxHash:            receipt.TxHash,
			TxIndex:           receipt.TransactionIndex,
			From:              from,
			To:                tx.To(),
			CumulativeGasUsed: receipt.CumulativeGasUsed,
			EffectiveGasPrice: receipt.EffectiveGasPrice,
			GasUsed:           receipt.GasUsed,
			ContractAddress:   &receipt.ContractAddress,
			Logs:              receipt.Logs,
			LogsBloom:         receipt.Bloom,
			Type:              receipt.Type,
			Status:            receipt.Status,
			BlobGasUsed:       receipt.BlobGasUsed,
			BlobGasPrice:      receipt.BlobGasPrice,
		}

		receipts[i] = decodedReceipt
	}

	return receipts, nil
}

// StoreDecodedBlock stores a decoded block in Redis as JSON.
func (e *Ethereum) StoreDecodedBlock(block *DecodedBlock) error {
	blockJSON, err := json.Marshal(block)
	if err != nil {
		return fmt.Errorf("failed to marshal decoded block: %w", err)
	}
	key := fmt.Sprintf("ethereum:block:%s", block.Number.String())
	if err := e.Factory.Data.RB.Set(e.Factory.Ctx, key, blockJSON, 0).Err(); err != nil {
		return fmt.Errorf("failed to store decoded block in Redis: %w", err)
	}
	return nil
}
