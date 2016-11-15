// Package testbot uses all of the features of graw and can be controlled
// remotely with messages (gently).
package main

import (
	"log"

	"github.com/turnage/graw"
)

type bot struct {
	graw.Account
	graw.Stopper
}

func main() {
	log.Printf("Error: %v", graw.Run(config, &bot{}))
}
