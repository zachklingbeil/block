package ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

func (e *Ethereum) processTransaction(ctx context.Context, tx *types.Transaction, signer types.Signer) *Transactions {
	txInfo := &Transactions{
		Hash:       tx.Hash().Hex(),
		Value:      tx.Value(),
		Gas:        tx.Gas(),
		GasPrice:   tx.GasPrice(),
		Nonce:      tx.Nonce(),
		DataLength: len(tx.Data()),
		Type:       tx.Type(),
	}

	if addr, err := types.Sender(signer, tx); err == nil {
		txInfo.From = addr.Hex()
	}

	if tx.To() == nil {
		txInfo.To = "Contract Creation"
	} else {
		txInfo.To = tx.To().Hex()
	}

	receipt, err := e.Factory.Eth.TransactionReceipt(ctx, tx.Hash())
	if err == nil {
		txInfo.Status = receipt.Status
		txInfo.CumulativeGasUsed = receipt.CumulativeGasUsed
		for _, log := range receipt.Logs {
			switch {
			case isERC20or721Transfer(log):
				txInfo.Logs = append(txInfo.Logs, &LogInfo{
					Address:   log.Address.Hex(),
					EventType: "ERC20/ERC721 Transfer",
					From:      "0x" + log.Topics[1].Hex()[26:],
					To:        "0x" + log.Topics[2].Hex()[26:],
					Value:     new(big.Int).SetBytes(log.Data),
				})
			case isERC1155TransferSingle(log):
				txInfo.Logs = append(txInfo.Logs, &LogInfo{
					Address:   log.Address.Hex(),
					EventType: "ERC1155 TransferSingle",
					Operator:  "0x" + log.Topics[1].Hex()[26:],
					From:      "0x" + log.Topics[2].Hex()[26:],
					To:        "0x" + log.Topics[3].Hex()[26:],
					ID:        new(big.Int).SetBytes(log.Data[:32]),
					Value:     new(big.Int).SetBytes(log.Data[32:]),
				})
			case isERC1155TransferBatch(log):
				logInfo := &LogInfo{
					Address:   log.Address.Hex(),
					EventType: "ERC1155 TransferBatch",
					Operator:  "0x" + log.Topics[1].Hex()[26:],
					From:      "0x" + log.Topics[2].Hex()[26:],
					To:        "0x" + log.Topics[3].Hex()[26:],
				}
				if len(log.Data) >= 128 {
					logInfo.IDs, logInfo.Values = decode1155Batch(log.Data)
				}
				txInfo.Logs = append(txInfo.Logs, logInfo)
			}
		}
	}

	return txInfo
}

func isERC20or721Transfer(log *types.Log) bool {
	const transferEvent = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	return len(log.Topics) == 3 && log.Topics[0].Hex() == transferEvent && len(log.Data) == 32
}

func isERC1155TransferSingle(log *types.Log) bool {
	const transfer1155SingleEvent = "0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62"
	return len(log.Topics) == 4 && log.Topics[0].Hex() == transfer1155SingleEvent && len(log.Data) == 64
}

func isERC1155TransferBatch(log *types.Log) bool {
	const transfer1155BatchEvent = "0x4a39dc06d4c0dbc64b70b1b5fdcf9a43c3b840ecb9c7aafb5c62c0124c6a16e3"
	return len(log.Topics) == 4 && log.Topics[0].Hex() == transfer1155BatchEvent && len(log.Data) >= 64
}

func decode1155Batch(data []byte) ([]string, []string) {
	ids := []string{}
	values := []string{}
	if len(data) < 128 {
		return ids, values
	}
	idsOffset := new(big.Int).SetBytes(data[:32]).Int64()
	valuesOffset := new(big.Int).SetBytes(data[32:64]).Int64()
	idsStart := int(idsOffset)
	valuesStart := int(valuesOffset)

	idsLen := new(big.Int).SetBytes(data[idsStart : idsStart+32]).Int64()
	for i := int64(0); i < idsLen; i++ {
		id := new(big.Int).SetBytes(data[idsStart+32+int(i)*32 : idsStart+32+int(i+1)*32])
		ids = append(ids, id.String())
	}

	valuesLen := new(big.Int).SetBytes(data[valuesStart : valuesStart+32]).Int64()
	for i := int64(0); i < valuesLen; i++ {
		val := new(big.Int).SetBytes(data[valuesStart+32+int(i)*32 : valuesStart+32+int(i+1)*32])
		values = append(values, val.String())
	}
	return ids, values
}
