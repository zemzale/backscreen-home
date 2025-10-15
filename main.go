package main

import (
	"log"

	"github.com/zemzale/backscreen-home/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
