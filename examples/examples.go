package main

import (
	"fmt"
	"os"

	"github.com/paytonturnage/graw"
)

func main() {
	agent, err := graw.NewAgentFromFile("useragent.protobuf")
	if err != nil {
		fmt.Printf("Failed to make agent: %v", err)
		os.Exit(-1)
	}

	redditor, err := agent.Me()
	if err != nil {
		fmt.Printf("Failed to get own account: %v", err)
		os.Exit(-1)
	}
	fmt.Printf("Account:\n%s", redditor.String())

	karmaList, err := agent.MeKarma()
	if err != nil {
		fmt.Printf("Failed to get own account: %v", err)
		os.Exit(-1)
	}
	fmt.Printf("Karma breakdown:\n%v", karmaList)
}
