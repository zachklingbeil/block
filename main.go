package main

import (
	"github.com/zachklingbeil/block/loopring"
	"github.com/zachklingbeil/block/process"
	"github.com/zachklingbeil/factory"
	"github.com/zachklingbeil/peer"
)

func main() {
	factory := factory.Assemble("timefactory")
	peer.HelloPeers(factory)
	go loopring.Connect(factory)
	process.InitProcess(factory)
	// process := process.InitProcess(factory)
	// process.ProcessTransactions()
	select {}
}
