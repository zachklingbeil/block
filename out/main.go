package out

import (
	"github.com/zachklingbeil/factory"
)

type Output struct {
	Factory *factory.Factory
	Peers   []Peer
}
