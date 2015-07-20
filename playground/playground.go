// Package playground is used to easily find the format of and explore the
// endpoints of the reddit api.
//
// Example useage from project root:
//
//   playground/playground --useragent=useragent.protobuf --get \
//   --url=https://oauth.reddit.com/api/v1/me
//
// This will output Reddit's JSON response to the shell.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/paytonturnage/graw/data"
	"github.com/paytonturnage/graw/nface"
)

var (
	requestURL = flag.String("url", "", "url to make a request to")
	get        = flag.Bool("get", false, "make a get request to the url")
	post       = flag.Bool("post", false, "make a post request to the url")
	userAgent  = flag.String("useragent", "", "user agent protobuf file")
)

func main() {
	flag.Parse()

	if *userAgent == "" {
		fmt.Printf("You must provide a user agent file.\n")
		os.Exit(-1)
	}

	agent, err := data.NewUserAgentFromFile(*userAgent)
	if err != nil {
		fmt.Printf("Failed to load user agent: %v\n", err)
		os.Exit(-1)
	}

	client, err := nface.NewClient(agent)
	if err != nil {
		fmt.Printf("Failed to create Graw entity: %v\n", err)
		os.Exit(-1)
	}

	if *get {
		raw, err := client.Raw(&nface.Request{
			Action: nface.GET,
			URL:    *requestURL,
		})
		if err != nil {
			fmt.Printf("GET request failed: %v\n", err)
		} else {
			fmt.Printf("GET response:\n\n%s\n", raw)
		}
	}

	if *post {
		raw, err := client.Raw(&nface.Request{
			Action: nface.POST,
			URL:    *requestURL,
		})
		if err != nil {
			fmt.Printf("POST request failed: %v\n", err)
		} else {
			fmt.Printf("POST response:\n\n%s\n", raw)
		}
	}
}
