package api

import (
	"fmt"

	"github.com/paytonturnage/graw/nface"
)

const (
	// baseURL is the base url for all api calls.
	baseURL = "https://oauth.reddit.com/api"
	// meURL is the url exension for the /v1/me api call.
	meURL = "/v1/me"
	// meKarmaURL is the url extension /v1/me/karma api call.
	meKarmaURL = "/v1/me/karma"
)

// MeRequest returns an nface.Request representing a /v1/me call.
func MeRequest() *nface.Request {
	return &nface.Request{
		Action:  nface.GET,
		BaseURL: fmt.Sprintf("%s%s", baseURL, meURL),
	}
}

// MeKarmaRequest returns an nface.Request representing a /v1/me/karma call.
func MeKarmaRequest() *nface.Request {
	return &nface.Request{
		Action:  nface.GET,
		BaseURL: fmt.Sprintf("%s%s", baseURL, meKarmaURL),
	}
}
