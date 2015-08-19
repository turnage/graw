// Package operator makes api calls to reddit.
package operator

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/turnage/graw/bot/internal/client"
	"github.com/turnage/graw/bot/internal/operator/internal/request"
	"github.com/turnage/redditproto"
)

// Operator makes api calls to reddit.
type Operator struct {
	cli client.Client
}

// New returns a new operator which uses cli as its client.
func New(cli client.Client) *Operator {
	return &Operator{cli: cli}
}

// Scrape returns posts from a subreddit, in the specified sort order, with the
// specified reference points for direction, up to lim. lims above 100 are
// ineffective because Reddit will return only 100 posts per query.
func (o *Operator) Scrape(sub, sort, after, before string, lim uint) ([]*redditproto.Link, error) {
	req, err := request.New(
		"GET",
		fmt.Sprintf("https://oauth.reddit.com/r/%s/%s.json", sub, sort),
		&url.Values{
			"limit":  []string{strconv.Itoa(int(lim))},
			"after":  []string{after},
			"before": []string{before},
		},
	)
	if err != nil {
		return nil, err
	}

	return o.getLinkListing(req)
}

// Threads returns specific threads, requested by their fullname (t[1-6]_[id]).
func (o *Operator) Threads(fullnames ...string) ([]*redditproto.Link, error) {
	ids := strings.Join(fullnames, ",")
	req, err := request.New(
		"GET",
		fmt.Sprintf("https://oauth.reddit.com/by_id/%s", ids),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return o.getLinkListing(req)
}

// getLinkListing executes a request and returns the reddit posts in the
// returned link listing.
func (o *Operator) getLinkListing(r *http.Request) ([]*redditproto.Link, error) {
	response := &redditproto.LinkListing{}
	err := o.cli.Do(r, response)
	return getLinks(response), err
}
