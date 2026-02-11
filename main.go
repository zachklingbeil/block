package main

import (
	"log"
	"os"

	"github.com/timefactoryio/block/fx"
)

func main() {
	f := fx.Init(os.Getenv("PASSWORD"))
	defer f.Close()
	if err := f.Test(); err != nil {
		log.Fatalf("Test failed: %v", err)
	}
}
