package value

import (
	_ "embed"
	"fmt"
)

func (v *Value) Refresh() {
	for i := range v.Peers {
		fmt.Printf("%d\n", i)
		peer := v.Peers[i]
		v.Format(peer.Address)
		v.Save(peer)
	}
}

// get loopring ids
// func (v *Value) Refresh() {
// 	for i := range v.Peers {
// 		fmt.Printf("%d\n", i)
// 		peer := v.Peers[i]
// 		if peer.LoopringID == "." || peer.LoopringID == "!" || peer.LoopringID == "" {
// 			v.GetLoopringID(peer)
// 		}
// 	}
// }
