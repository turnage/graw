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
	"io/ioutil"
	"os"

	"github.com/paytonturnage/graw"
	"github.com/paytonturnage/graw/internal/request"
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

	pilot, err := graw.NewUserFromFile(*userAgent)
	if err != nil {
		fmt.Printf("Failed to load user: %v\n", err)
	}

	if err := pilot.Auth(); err != nil {
		fmt.Printf("Failed to log user in: %v\n", err)
		os.Exit(-1)
	}

	method := ""
	if *get {
		method = "GET"
	} else if *post {
		method = "POST"
	}

	req, err := request.New(method, *requestURL, nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		os.Exit(-1)
	}

	resp, err := pilot.ExecRaw(req)
	if err != nil {
		fmt.Printf("Failed to execute request: %v\n", err)
		os.Exit(-1)
	}

	if resp.Body == nil {
		fmt.Printf("Response body was empty.\n")
		os.Exit(-1)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("There was an error reading the response: %v\n", err)
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Body: %s\n", buf)
}
