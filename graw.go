// Package graw is the Golang Reddit API Wrapper.
package graw

import (
	"strings"
	"time"

	"github.com/paytonturnage/graw/data"
	"github.com/paytonturnage/graw/nface"
)

const (
	// maxQueriesPerWindow is the amount of queries allowed per minute according
	// to the rules of the reddit api.
	maxQueriesPerWindow = 60
	// holdSymbol is the symbol in api urls to be replaced with values.
	holdSymbol = "#"
)

// Reddit api call urls.
const (
	// meURL is the url exension for the /v1/me api call.
	meURL = "/v1/me"
	// meKarmaURL is the url extension /v1/me/karma api call.
	meKarmaURL = "/v1/me/karma"
	// userURL is the url for fetching user info.
	userURL = "/user/#/about"
)

// Graw wraps the reddit api; all api calls go through Graw.
type Graw struct {
	// client manages the Graw's connection to reddit.
	client *nface.Client
	// windowStart is the start of Graw's query-limit window.
	windowStart time.Time
	// queriesInWindow is the number of queries made in the active window.
	queriesInWindow int
}

// NewGraw returns an graw ready to use with reddit. An error is returned
// if the graw cannot authenticate.
func NewGraw(userAgent *data.UserAgent) (*Graw, error) {
	client, err := nface.NewClient(userAgent)
	if err != nil {
		return nil, err
	}

	return &Graw{
		client: client,
		windowStart: time.Now(),
	}, nil
}

// Me wraps /v1/me. See https://www.reddit.com/dev/api#GET_api_v1_me
func (g *Graw) Me() (*data.Account, error) {
	resp := &data.Account{}
	err := g.do(&nface.Request{Action: nface.GET, URL: meURL}, resp)
	return resp, err
}

// MeKarma wraps /v1/me/karma. See
// https://www.reddit.com/dev/api#GET_api_v1_me_karma
func (g *Graw) MeKarma() ([]*data.SubredditKarma, error) {
	resp := &data.KarmaList{}
	err := g.do(&nface.Request{Action: nface.GET, URL: meKarmaURL}, resp)
	return resp.GetData(), err
}

// User wraps /user/username/about. See
// https://www.reddit.com/dev/api#GET_user_{username}_about
func (g *Graw) User(username string) (*data.Account, error) {
	resp := &data.Account{}
	err := g.do(&nface.Request{
		Action: nface.GET,
		URL:    strings.Replace(userURL, holdSymbol, username, 1),
	}, resp)
	return resp, err
}

// do executes a request using the Graw's client, and abides the query limit.
func (g *Graw) do(r *nface.Request, resp interface{}) error {
	g.clearWindowByBlocking()
	return g.client.Do(r, resp)
}

// clearWindowByBlocking keeps track of the Graw's query limit and blocks until
// it is allowed to make new queries if it reaches it.
func (g *Graw) clearWindowByBlocking() {
	timeInWindow := time.Now().Sub(g.windowStart)
	if timeInWindow < time.Minute {
		if g.queriesInWindow >= maxQueriesPerWindow {
			time.Sleep(time.Minute - timeInWindow)
			g.resetWindow()
		}
		g.queriesInWindow++
	} else {
		g.resetWindow()
	}
}

func (g *Graw) resetWindow() {
	g.windowStart = time.Now()
	g.queriesInWindow = 0
}
