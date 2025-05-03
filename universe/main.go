package universe

import (
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zachklingbeil/block/value"
	"github.com/zachklingbeil/factory"
)

type Universe struct {
	Factory *factory.Factory
	Map     map[string]*any
	Value   *value.Value
}

func NewUniverse(factory *factory.Factory) *Universe {
	u := &Universe{
		Factory: factory,
		Map:     make(map[string]*any),
		Value:   value.NewValue(factory),
	}

	return u
}

func (u *Universe) AddToMap(key string, value any) {
	if !common.IsHexAddress(strings.ToLower(key)) {
		log.Printf("Invalid key: %s. Entry not added to the map.", key)
		return
	}
	var one any = value
	u.Map[key] = &one
}

func (u *Universe) GetFromMap(key string) (*any, bool) {
	if !common.IsHexAddress(strings.ToLower(key)) {
		log.Printf("Invalid key: %s. Entry not found in the map.", key)
		return nil, false
	}
	value, exists := u.Map[key]
	return value, exists
}
