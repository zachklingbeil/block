package other

// type Ethereum struct {
// 	Chain      *params.ChainConfig
// 	SigService *SigProvider
// }

// func New(ethURL, sigProviderURL string) (*Ethereum, error) {

// 	ethereum := &Ethereum{

// 		Chain:      params.MainnetChainConfig,
// 		SigService: NewSigProvider(sigProviderURL),
// 	}
// 	return ethereum, nil
// }

// func (e *Ethereum) FetchAndDecodeBlock(ctx context.Context, blockNum int64) (*types.Block, []map[string]interface{}, []map[string]interface{}, error) {
// 	block, err := e.Factory.Eth.BlockByNumber(ctx, big.NewInt(blockNum))
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}

// 	var decodedInputs []map[string]interface{}
// 	var decodedEvents []map[string]interface{}

// 	for _, tx := range block.Transactions() {
// 		// Decode transaction input
// 		if len(tx.Data()) > 0 {
// 			input, err := e.SigService.DecodeFunctionInput(fmt.Sprintf("0x%x", tx.Data()))
// 			if err == nil && len(input) > 0 {
// 				decodedInputs = append(decodedInputs, input...)
// 			}
// 		}

// 		// Fetch receipt and decode logs/events
// 		receipt, err := e.Factory.Eth.TransactionReceipt(ctx, tx.Hash())
// 		if err != nil {
// 			continue // skip if receipt not found
// 		}
// 		for _, logEntry := range receipt.Logs {
// 			topics := make([]string, len(logEntry.Topics))
// 			for i, t := range logEntry.Topics {
// 				topics[i] = t.Hex()
// 			}
// 			dataHex := fmt.Sprintf("0x%x", logEntry.Data)
// 			var topic string
// 			if len(topics) > 0 {
// 				topic = topics[0]
// 			} else {
// 				topic = ""
// 			}
// 			decoded, err := e.SigService.DecodeEvent(dataHex, topic)
// 			if err == nil && len(decoded) > 0 {
// 				decodedEvents = append(decodedEvents, decoded...)
// 			}
// 		}
// 	}

// 	return block, decodedInputs, decodedEvents, nil
// }
