package graw

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/turnage/redditproto"
)

// exec Executes a reddit api call and parses the returned json into the out
// interface.
func exec(c client, r *http.Request, out interface{}) error {
	rawResp, err := c.Do(r)
	if err != nil {
		return err
	}

	if rawResp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code in response")
	}

	defer func() {
		if rawResp.Body != nil {
			rawResp.Body.Close()
		}
	}()

	decoder := json.NewDecoder(rawResp.Body)
	return decoder.Decode(out)
}

// scrape scrapes a subreddit for its top posts, according to the sorting method
// specified, using cli. It returns a slice of posts. lim is only meaningful up
// to 100, because Reddit will not return more than that. To get more posts in
// order, use the "before" and "after" fields; they take fullnames of posts for
// reference. "after" gets posts posted before the provided post, and "before"
// gets posts posted after the provided post (Reddit implements and names them
// this way because generally scraping is done working backward in time.)
func scrape(cli client, sub, sort, after, before string,
	lim int) ([]*redditproto.Link, error) {
	response := &struct {
		Data struct {
			Children []struct {
				Data *redditproto.Link
			}
		}
	}{}
	req, err := newRequest(
		"GET",
		fmt.Sprintf("https://oauth.reddit.com/r/%s/%s.json", sub, sort),
		&url.Values{
			"limit":  []string{strconv.Itoa(lim)},
			"after":  []string{after},
			"before": []string{before},
		},
	)
	if err != nil {
		return nil, err
	}

	err = exec(cli, req, response)
	if err != nil {
		return nil, err
	}

	links := make([]*redditproto.Link, len(response.Data.Children))
	for i, child := range response.Data.Children {
		links[i] = child.Data
	}

	return links, nil
}
