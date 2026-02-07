package main

import (
	"log"

	"github.com/timefactoryio/block/fx"
)

func main() {
	f := fx.Init("")
	defer f.Close()
	if err := f.Test(); err != nil {
		log.Fatalf("Test failed: %v", err)
	}
}
