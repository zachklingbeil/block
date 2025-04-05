package loop

import (
	"fmt"

	"github.com/zachklingbeil/factory"
)

type Block struct {
	Created      int64         `json:"createdAt"`
	Number       int64         `json:"blockId"`
	Size         int64         `json:"blockSize"`
	TxHash       string        `json:"txHash"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	TxType    string `json:"txType"`
	From      int64  `json:"accountId"`
	To        int64  `json:"toAccountId"`
	ToAddress string `json:"toAccountAddress"`
}

type Loopring struct {
	Factory *factory.Factory
	Delta   int64
}

func NewLoopring(factory *factory.Factory) *Loopring {
	return &Loopring{
		Factory: factory,
	}
}

func (l *Loopring) FetchBlocks() {
	current := l.currentBlock()
	blockHeight := l.blockHeight()
	if blockHeight == current {
		fmt.Println("blockHeight = currentBlock")
		return
	}
	for i := blockHeight + 1; i <= current; i++ {
		if err := l.GetBlock(int(i)); err != nil {
			fmt.Printf("Failed to fetch block %d: %v\n", i, err)
			continue
		}
	}

}
