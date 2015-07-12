package api

import (
	"fmt"

	"github.com/paytonturnage/graw/nface"
)

const (
	// baseURL is the base url for all api calls.
	baseURL = "https://oauth.reddit.com/api"
	// meURL is the url exension for the me api call.
	meURL = "/v1/me"
)

// MeRequest returns an nface.Request representing a Me call.
func MeRequest() *nface.Request {
	return &nface.Request{
		Action:  nface.GET,
		BaseURL: fmt.Sprintf("%s%s", baseURL, meURL),
	}
}
