package main

import (
	"github.com/zachklingbeil/block/loopring"
	"github.com/zachklingbeil/block/process"
	"github.com/zachklingbeil/factory"
	"github.com/zachklingbeil/peer"
)

func main() {
	factory := factory.Assemble("timefactory")
	p := peer.HelloPeers(factory)

	go loopring.Connect(factory)
	go p.HelloUniverse()

	process := process.InitProcess(factory, p)
	process.ProcessTransactions()
	process.PrintExampleTxForEachType()
	select {}
}
