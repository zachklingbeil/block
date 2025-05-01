package value

import (
	"fmt"
	"time"

	"github.com/zachklingbeil/factory"
	"golang.org/x/time/rate"
)

type Value struct {
	Factory  *factory.Factory
	Peers    []*Peer
	Tokens   []*Token
	Map      map[string]*Peer
	TokenMap map[any]*Token
}

func NewValue(factory *factory.Factory) *Value {
	v := &Value{
		Factory:  factory,
		Map:      make(map[string]*Peer),
		TokenMap: make(map[any]*Token),
	}

	v.LoadTokens()
	v.LoadPeers()
	v.rebuildMap()
	v.DotLoop()
	return v
}

func (v *Value) Refresh() {
	for i := range v.Peers {
		fmt.Printf("%d", i)
		peer := v.Peers[i]
		// Clear invalid fields if they are "." or "!"
		if peer.ENS == "." || peer.ENS == "!" {
			peer.ENS = ""
		}
		if peer.LoopringID == "." || peer.LoopringID == "!" {
			peer.LoopringID = ""
		}
		if peer.LoopringENS == "." || peer.LoopringENS == "!" {
			peer.LoopringENS = ""
		}
		if peer.Address == "." || peer.Address == "!" {
			peer.Address = ""
		}
		v.HelloUniverse(peer.Address)
	}
}

// func (v *Value) DotLoop() {
// 	for i := range v.Peers {
// 		fmt.Printf("%d", i)
// 		peer := v.Peers[i]
// 		// Clear invalid fields if they are "." or "!"
// 		if peer.ENS == "." || peer.ENS == "!" {
// 			peer.ENS = ""
// 		}
// 		if peer.LoopringID == "." || peer.LoopringID == "!" {
// 			peer.LoopringID = ""
// 		}
// 		if peer.LoopringENS == "." || peer.LoopringENS == "!" {
// 			peer.LoopringENS = ""
// 		}
// 		v.GetLoopringENS(peer)
// 		fmt.Printf("%s %s %s\n", peer.ENS, peer.LoopringENS, peer.LoopringID)
// 	}
// }

func (v *Value) DotLoop() {
	// Create a rate limiter for 20 requests per second
	limiter := rate.NewLimiter(rate.Every(50*time.Millisecond), 1) // 50ms per request = 20 RPS

	for i := range v.Peers {
		// Wait for permission to proceed
		if err := limiter.Wait(v.Factory.Ctx); err != nil {
			fmt.Printf("Rate limiter error: %v\n", err)
			continue
		}

		fmt.Printf("%d", i)
		peer := v.Peers[i]

		// Clear invalid fields if they are "." or "!"
		if peer.ENS == "." || peer.ENS == "!" {
			peer.ENS = ""
		}
		if peer.LoopringID == "." || peer.LoopringID == "!" {
			peer.LoopringID = ""
		}
		if peer.LoopringENS == "." || peer.LoopringENS == "!" {
			peer.LoopringENS = ""
		}

		// Fetch Loopring ENS
		v.GetLoopringENS(peer)

		// Print peer details
		fmt.Printf("%s %s %s\n", peer.ENS, peer.LoopringENS, peer.LoopringID)
	}
}
