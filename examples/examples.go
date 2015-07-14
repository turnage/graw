package main

import (
	"fmt"
	"os"

	"github.com/paytonturnage/graw"
)

func main() {
	graw, err := graw.NewGrawFromFile("useragent.protobuf")
	if err != nil {
		fmt.Printf("Failed to make graw: %v", err)
		os.Exit(-1)
	}

	redditor, err := graw.Me()
	if err != nil {
		fmt.Printf("Failed to get own account: %v", err)
		os.Exit(-1)
	}
	fmt.Printf("Account:\n%s", redditor.String())

	karmaList, err := graw.MeKarma()
	if err != nil {
		fmt.Printf("Failed to get own account: %v", err)
		os.Exit(-1)
	}
	fmt.Printf("Karma breakdown:\n%v", karmaList)
}
