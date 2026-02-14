package main

import (
	"log"
	"os"

	"github.com/timefactoryio/block/one"
)

func main() {
	pw := os.Getenv("PASSWORD")
	f := one.Init(pw)
	defer f.Close()
	if err := f.Test(); err != nil {
		log.Fatalf("Test failed: %v", err)
	}
}
