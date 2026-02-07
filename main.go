package main

import (
	"github.com/timefactoryio/block/fx"
)

func main() {
	f := fx.Init()
	f.Node()
	defer f.Rpc.Close()
}
