// Package api makes api calls to Reddit.
package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/turnage/redditproto"
)

const (
	// maxLinks is the amount of posts reddit will return for a scrape
	// query.
	maxLinks = 100
	// deletedAuthor is the author value if a post or comment was deleted.
	deletedAuthor = "[deleted]"
)

var (
	// formEncoding is the encoding format of parameters in the body of
	// requests sent to Reddit.
	formEncoding = map[string][]string{
		"content-type": {"application/x-www-form-urlencoded"},
	}
	// domain is the domain Reddit lives on.
	domain = "reddit.com"
	// oauth2Host is the hostname of Reddit's OAuth2 server.
	oauth2Host = "oauth." + domain
	// baseURL is the url all requests extend from.
	baseURL = "https://" + oauth2Host
)

// Requester executes and http.Request and returns the bytes in the body of the
// response.
type Requester func(*http.Request) ([]byte, error)

// SetTestDomain is a test hook for end to end tests to specify an alternate,
// test instance of Reddit to run against.
func SetTestDomain(domain string) {
	oauth2Host = "oauth." + domain
	baseURL = "https://" + domain
}

// Scrape returns slices with the content of a listing endpoint.
func Scrape(
	r Requester,
	path,
	after string,
	limit int,
) (
	[]*redditproto.Link,
	[]*redditproto.Comment,
	[]*redditproto.Message,
	error,
) {
	if limit <= 0 {
		limit = maxLinks
	}

	bytes, err := r(
		&http.Request{
			Method:     "GET",
			ProtoMajor: 2,
			ProtoMinor: 0,
			URL: &url.URL{
				Scheme: "https",
				Host:   oauth2Host,
				Path:   path,
				RawQuery: url.Values{
					"limit": []string{strconv.Itoa(int(limit))},
					// This looks like a mistake but it
					// isn't. Reddit thinks of listings
					// reverse-chronologically. What graw
					// calls "after" refers to things posted
					// at a time later than the reference
					// point. What Reddit calls "after"
					// refers to things that come up in the
					// listing later, assuming the head is
					// the latest post. I don't think the
					// Reddit naming makes sense so I hide
					// it except for this line.
					"before":   []string{after},
					"raw_json": []string{"1"},
				}.Encode(),
			},
			Host: oauth2Host,
		},
	)
	if err != nil {
		return nil, nil, nil, err
	}

	return redditproto.ParseListing(bytes)
}

// IsThereThing returns whether a thing by the given id exists.
func IsThereThing(r Requester, id string) (bool, error) {
	path := "/api/info.json"

	// api/info doesn't provide message types; these need to be fetched from
	// a different url.
	if strings.HasPrefix(id, "t4_") {
		id := strings.TrimPrefix(id, "t4_")
		path = fmt.Sprintf("/message/messages/%s", id)
	}

	bytes, err := r(
		&http.Request{
			Method:     "GET",
			ProtoMajor: 2,
			ProtoMinor: 0,
			URL: &url.URL{
				Scheme: "https",
				Host:   oauth2Host,
				Path:   path,
				RawQuery: url.Values{
					"id":       []string{id},
					"raw_json": []string{"1"},
				}.Encode(),
			},
			Host: oauth2Host,
		},
	)

	if err != nil {
		return false, err
	}

	links, comments, messages, err := redditproto.ParseListing(bytes)
	if err != nil {
		return false, err
	}

	if len(links) == 1 {
		return links[0].GetAuthor() != deletedAuthor, nil
	}

	if len(comments) == 1 {
		return comments[0].GetAuthor() != deletedAuthor, nil
	}

	if len(messages) == 1 {
		return true, nil
	}

	return false, nil
}

// Thread returns a link; the Comments field will be filled with the comment
// tree. Browse each comment's reply tree from the ReplyTree field.
func Thread(r Requester, permalink string) (*redditproto.Link, error) {
	bytes, err := r(
		&http.Request{
			Method:     "GET",
			ProtoMajor: 2,
			ProtoMinor: 0,
			URL: &url.URL{
				Scheme:   "https",
				Host:     oauth2Host,
				Path:     fmt.Sprintf("%s.json", permalink),
				RawQuery: "raw_json=1",
			},
			Host: oauth2Host,
		},
	)
	if err != nil {
		return nil, err
	}

	return redditproto.ParseThread(bytes)
}

// Reply replies to a post, message, or comment.
func Reply(r Requester, parent, content string) error {
	_, err := r(
		&http.Request{
			Method:     "POST",
			ProtoMajor: 2,
			ProtoMinor: 0,
			URL: &url.URL{
				Scheme: "https",
				Host:   oauth2Host,
				Path:   "/api/comment",
			},
			Header: formEncoding,
			Body: ioutil.NopCloser(
				bytes.NewBufferString(
					url.Values{
						"thing_id": []string{parent},
						"text":     []string{content},
					}.Encode(),
				),
			),
			Host: oauth2Host,
		},
	)
	return err
}

// Compose sends a private message to a user.
func Compose(r Requester, user, subject, content string) error {
	_, err := r(
		&http.Request{
			Method:     "POST",
			ProtoMajor: 2,
			ProtoMinor: 0,
			URL: &url.URL{
				Scheme: "https",
				Host:   oauth2Host,
				Path:   "/api/compose",
			},
			Header: formEncoding,
			Body: ioutil.NopCloser(
				bytes.NewBufferString(
					url.Values{
						"to":      []string{user},
						"subject": []string{subject},
						"text":    []string{content},
					}.Encode(),
				),
			),
			Host: oauth2Host,
		},
	)
	return err
}

// Submit submits a post.
func Submit(r Requester, subreddit, kind, title, content string) error {
	_, err := r(
		&http.Request{
			Method:     "POST",
			ProtoMajor: 2,
			ProtoMinor: 0,
			URL: &url.URL{
				Scheme: "https",
				Host:   oauth2Host,
				Path:   "/api/submit",
			},
			Header: formEncoding,
			Body: ioutil.NopCloser(
				bytes.NewBufferString(
					url.Values{
						"sr":    []string{subreddit},
						"kind":  []string{kind},
						"title": []string{title},
						"url":   []string{content},
						"text":  []string{content},
					}.Encode(),
				),
			),
			Host: oauth2Host,
		},
	)
	return err
}
