package universe

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zachklingbeil/factory"
)

type Universe struct {
	Factory *factory.Factory
	Map     map[common.Address]struct{}
}

func NewUniverse(factory *factory.Factory) *Universe {
	u := &Universe{
		Factory: factory,
		Map:     make(map[common.Address]struct{}),
	}
	return u
}

func (u *Universe) AddStructToMap(value Value) {
	address := value.GetAddress()
	if !common.IsHexAddress(address) {
		log.Printf("Invalid Ethereum address: %s. Entry not added to the map.", address)
		return
	}

	u.Map[common.HexToAddress(address)] = struct{}{}
	log.Printf("Added struct to map with key: %s", address)
}

func (u *Universe) GetStructFromMap(address string) (struct{}, bool) {
	if !common.IsHexAddress(address) {
		log.Printf("Invalid Ethereum address: %s. Entry not found in the map.", address)
		return struct{}{}, false
	}

	value, exists := u.Map[common.HexToAddress(address)]
	if !exists {
		log.Printf("No entry found in the map for address: %s", address)
	}
	return value, exists
}
