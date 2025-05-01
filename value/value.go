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
	// v.DotLoop()
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
// 		if peer.LoopringENS == "." || peer.LoopringENS == "!" {
// 			peer.LoopringENS = ""
// 		}
// 		v.GetLoopringENS(peer)
// 		fmt.Printf("%s %s %s\n", peer.ENS, peer.LoopringENS, peer.LoopringID)
// 	}
// }

func (v *Value) DotLoop() {
	// No burst: set burst to 1, and use a conservative rate (adjust as needed)
	limiter := rate.NewLimiter(rate.Every(334*time.Millisecond), 1)

	// Count how many peers need processing
	toProcess := 0
	for _, peer := range v.Peers {
		if peer.LoopringENS == "." || peer.LoopringENS == "!" {
			toProcess++
		}
	}

	for _, peer := range v.Peers {
		// Only process peers where LoopringENS is "." or "!"
		if peer.LoopringENS != "." && peer.LoopringENS != "!" {
			continue
		}

		if err := limiter.Wait(v.Factory.Ctx); err != nil {
			fmt.Printf("Rate limiter error: %v\n", err)
			continue
		}

		peer.LoopringENS = ""

		// Fetch Loopring ENS
		v.GetLoopringENS(peer)

		// Print peer details, decrementing the count each time
		toProcess--
		fmt.Printf("%d %s %s %s\n", toProcess, peer.ENS, peer.LoopringENS, peer.LoopringID)
	}
}
