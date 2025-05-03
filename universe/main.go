package universe

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zachklingbeil/block/universe/peer"
	"github.com/zachklingbeil/block/universe/token"
	"github.com/zachklingbeil/factory"
)

type Universe struct {
	Factory *factory.Factory
	Peer    *peer.Peers
	Token   *token.Tokens
	Map     map[common.Address]struct{}
}

func (u *Universe) AddStructToMap(value Value) {
	if !common.IsHexAddress(value.Address) {
		log.Printf("Invalid Ethereum address: %s. Entry not added to the map.", value.Address)
		return
	}
	u.Map[common.HexToAddress(value.Address)] = struct{}{}
	log.Printf("Added struct to map with key: %s", value.Address)
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

type Value struct {
	ENS   string `json:"ens,omitempty"`
	Token string `json:"token,omitempty"`

	LoopringENS string `json:"loopringEns,omitempty"`
	LoopringID  string `json:"loopringId,omitempty"`
	Address     string `json:"address,omitempty"`
	FirstBlock  string `json:"firstBlock,omitempty"`
	Decimals    string `json:"decimals,omitempty"`
	TokenId     string `json:"tokenId,omitempty"`
}

// func NewUniverse(factory *factory.Factory) *Universe {
// 	u := &Universe{
// 		Factory: factory,
// 		Map:     make(map[string]*any),
// 		Slice:   make([]*any, 0),
// 	}
// 	return u
// }

// func (u *Universe) AddToMap(key string, value any) {
// 	if !common.IsHexAddress(strings.ToLower(key)) {
// 		log.Printf("Invalid key: %s. Entry not added to the map.", key)
// 		return
// 	}
// 	var one any = value
// 	u.Map[key] = &one
// }

// func (u *Universe) GetFromMap(key string) (*any, bool) {
// 	if !common.IsHexAddress(strings.ToLower(key)) {
// 		log.Printf("Invalid key: %s. Entry not found in the map.", key)
// 		return nil, false
// 	}
// 	value, exists := u.Map[key]
// 	return value, exists
// }

// func (u *Universe) AppendToSlice(value *any) {
// 	u.Slice = append(u.Slice, value)
// }

// // func (u *Universe) AppendToSlice(value any) {
// // 	var obj any = value
// // 	u.Slice = append(u.Slice, &obj)
// // }

// func (u *Universe) BuildMapFromSlices() {
// 	for _, slice := range u.Slice {
// 		if slicePtr, ok := (*slice).([]any); ok {
// 			for _, item := range slicePtr {
// 				switch v := item.(type) {
// 				case Peer:
// 					if v.Address != "" {
// 						u.Map[v.Address] = slice
// 					}
// 					if v.LoopringID != "" {
// 						u.Map[v.LoopringID] = slice
// 					}
// 				case Token:
// 					if v.Address != "" {
// 						u.Map[v.Address] = slice
// 					}
// 					if v.TokenId != "" {
// 						u.Map[v.Token] = slice
// 					}
// 				default:
// 					log.Printf("Unsupported type in slice: %+v", v)
// 				}
// 			}
// 		} else {
// 			log.Printf("Invalid slice type in Universe.Slice: %+v", slice)
// 		}
// 	}
// 	fmt.Printf("Map built with %d entries.\n", len(u.Map))
// }
