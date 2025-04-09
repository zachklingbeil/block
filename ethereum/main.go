package ethereum

import (
	"github.com/zachklingbeil/factory"
)

type Ethereum struct {
	Factory *factory.Factory
}

func NewEthereum(factory *factory.Factory) *Ethereum {
	return &Ethereum{
		Factory: factory,
	}
}

// type Ethereum struct {
// 	Status       string  `json:"status,omitempty"`
// 	Message      string  `json:"message,omitempty"`
// 	Transactions []EthTx `json:"result,omitempty"`
// }

// type EthTx struct {
// 	TxId        int    `json:"id,omitempty"`
// 	BlockNumber string `json:"blockNumber,omitempty"`
// 	Index       string `json:"transactionIndex,omitempty"`
// 	TimeStamp   string `json:"timeStamp,omitempty"`
// 	Zero        string `json:"from,omitempty"`
// 	One         string `json:"to,omitempty"`
// 	Value       string `json:"value,omitempty"`
// 	Gas         string `json:"gas,omitempty"`
// 	GasPrice    string `json:"gasPrice,omitempty"`
// 	GasUsed     string `json:"gasUsed,omitempty"`
// 	Token       string `json:"tokenSymbol,omitempty"`
// }
