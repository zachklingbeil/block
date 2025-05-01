package ethereum

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

func (e *Ethereum) processTransaction(ctx context.Context, tx *types.Transaction) *Transactions {
	txInfo := &Transactions{
		Hash:       tx.Hash().Hex(),
		Value:      tx.Value(),
		Gas:        tx.Gas(),
		GasPrice:   tx.GasPrice(),
		Nonce:      tx.Nonce(),
		DataLength: len(tx.Data()),
		Type:       tx.Type(),
	}

	// Get sender address
	var from string
	var err error
	if tx.ChainId() == nil || tx.ChainId().Sign() == 0 {
		addr, err := types.Sender(types.HomesteadSigner{}, tx)
		if err == nil {
			from = addr.Hex()
		}
	} else {
		signer := types.LatestSignerForChainID(tx.ChainId())
		addr, err2 := types.Sender(signer, tx)
		if err2 == nil {
			from = addr.Hex()
		}
		err = err2
	}
	if err == nil {
		txInfo.From = from
	}

	// To address (contract creation if nil)
	if tx.To() == nil {
		txInfo.To = "Contract Creation"
	} else {
		txInfo.To = tx.To().Hex()
	}

	// Get receipt for transaction status, logs, etc.
	receipt, err := e.Factory.Eth.TransactionReceipt(ctx, tx.Hash())
	if err == nil {
		txInfo.Status = receipt.Status
		txInfo.CumulativeGasUsed = receipt.CumulativeGasUsed
		for _, log := range receipt.Logs {
			logInfo := &LogInfo{
				Address:    log.Address.Hex(),
				DataLength: len(log.Data),
			}
			for _, topic := range log.Topics {
				logInfo.Topics = append(logInfo.Topics, topic.Hex())
			}

			switch {
			case isERC20or721Transfer(log):
				processERC20or721Transfer(log, logInfo)
			case isERC1155TransferSingle(log):
				processERC1155TransferSingle(log, logInfo)
			case isERC1155TransferBatch(log):
				processERC1155TransferBatch(log, logInfo)
			}

			txInfo.Logs = append(txInfo.Logs, logInfo)
		}
	}

	return txInfo
}

func isERC20or721Transfer(log *types.Log) bool {
	const transferEvent = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	return len(log.Topics) == 3 && log.Topics[0].Hex() == transferEvent && len(log.Data) == 32
}

func processERC20or721Transfer(log *types.Log, logInfo *LogInfo) {
	fromAddr := "0x" + log.Topics[1].Hex()[26:]
	toAddr := "0x" + log.Topics[2].Hex()[26:]
	amount := new(big.Int).SetBytes(log.Data)
	logInfo.Topics = append(logInfo.Topics,
		fmt.Sprintf("Transfer: from %s to %s value %s", fromAddr, toAddr, amount.String()))
}

func isERC1155TransferSingle(log *types.Log) bool {
	const transfer1155SingleEvent = "0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62"
	return len(log.Topics) == 4 && log.Topics[0].Hex() == transfer1155SingleEvent && len(log.Data) == 64
}

func processERC1155TransferSingle(log *types.Log, logInfo *LogInfo) {
	operator := "0x" + log.Topics[1].Hex()[26:]
	fromAddr := "0x" + log.Topics[2].Hex()[26:]
	toAddr := "0x" + log.Topics[3].Hex()[26:]
	id := new(big.Int).SetBytes(log.Data[:32])
	value := new(big.Int).SetBytes(log.Data[32:])
	logInfo.Topics = append(logInfo.Topics,
		fmt.Sprintf("ERC1155 TransferSingle: operator %s from %s to %s id %s value %s", operator, fromAddr, toAddr, id.String(), value.String()))
}

func isERC1155TransferBatch(log *types.Log) bool {
	const transfer1155BatchEvent = "0x4a39dc06d4c0dbc64b70b1b5fdcf9a43c3b840ecb9c7aafb5c62c0124c6a16e3"
	return len(log.Topics) == 4 && log.Topics[0].Hex() == transfer1155BatchEvent && len(log.Data) >= 64
}

func processERC1155TransferBatch(log *types.Log, logInfo *LogInfo) {
	operator := "0x" + log.Topics[1].Hex()[26:]
	fromAddr := "0x" + log.Topics[2].Hex()[26:]
	toAddr := "0x" + log.Topics[3].Hex()[26:]

	if len(log.Data) >= 128 {
		idsOffset := new(big.Int).SetBytes(log.Data[:32]).Int64()
		valuesOffset := new(big.Int).SetBytes(log.Data[32:64]).Int64()
		idsStart := int(idsOffset)
		valuesStart := int(valuesOffset)

		// ids array
		idsLen := new(big.Int).SetBytes(log.Data[idsStart : idsStart+32]).Int64()
		var ids []string
		for i := int64(0); i < idsLen; i++ {
			id := new(big.Int).SetBytes(log.Data[idsStart+32+int(i)*32 : idsStart+32+int(i+1)*32])
			ids = append(ids, id.String())
		}

		// values array
		valuesLen := new(big.Int).SetBytes(log.Data[valuesStart : valuesStart+32]).Int64()
		var values []string
		for i := int64(0); i < valuesLen; i++ {
			val := new(big.Int).SetBytes(log.Data[valuesStart+32+int(i)*32 : valuesStart+32+int(i+1)*32])
			values = append(values, val.String())
		}

		logInfo.Topics = append(logInfo.Topics,
			fmt.Sprintf("ERC1155 TransferBatch: operator %s from %s to %s ids %v values %v", operator, fromAddr, toAddr, ids, values))
	} else {
		logInfo.Topics = append(logInfo.Topics,
			fmt.Sprintf("ERC1155 TransferBatch: operator %s from %s to %s (unable to decode ids/values)", operator, fromAddr, toAddr))
	}
}
