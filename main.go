package main

import (
	"log"

	input "github.com/zachklingbeil/block/in"
	"github.com/zachklingbeil/factory"
)

func main() {
	factory, err := factory.NewFactory("block")
	if err != nil {
		log.Fatalf("Error creating factory: %v", err)
	}
	input.NewInput(factory)
}
