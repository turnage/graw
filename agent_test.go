package graw

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/paytonturnage/graw/nface"
)

// expected formats a string for error reporting expected versus actual values.
func expected(field string, expected, actual interface{}) string {
	return fmt.Sprintf(
		"%s incorrect; expected %v, got %v",
		field,
		expected,
		actual)
}

// newProxyClient returns an http.Client that redirects all requests to the
// redirect url.
func newProxyClient(redirect string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: func(r *http.Request) (*url.URL, error) {
				return url.Parse(redirect)
			},
		},
	}
}

// newServerFromResponse returns an httptest.Server that always responds with
// response and status.
func newServerFromResponse(status int, response string) *httptest.Server {
	responseWriter := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Header().Set("content-type", "application/json")
		fmt.Fprintln(w, response)
	}
	return httptest.NewServer(http.HandlerFunc(responseWriter))
}

// Tests the me api call and that all response data is properly represented.
func TestMe(t *testing.T) {
	server := newServerFromResponse(200, `{
		"has_mail": true, 
		"inbox_count": 2,
		"name": "fooBar", 
		"created": 123456789.0, 
		"created_utc": 1315269998.0, 
		"link_karma": 31, 
		"comment_karma": 557, 
		"is_gold": true,
		"gold_credits": 40,
		"gold_expiration": 23273468.0, 
		"is_mod": true, 
		"over_18": true,
		"has_verified_email": true,
		"hide_from_robots": true,
		"id": "5sryd", 
		"has_mod_mail": true
	}`)
	agent := &Agent{
		client: nface.NewClient(
			newProxyClient(server.URL),
			"test-client",
			server.URL),
	}
	me, err := agent.Me()
	if err != nil {
		t.Fatalf("failed to get self: %v", err)
	}

	if !me.GetHasMail() {
		t.Error(expected("has_mail", true, false))
	}

	if me.GetInboxCount() != 2 {
		t.Error(expected("inbox_count", 2, me.GetInboxCount()))
	}

	if me.GetName() != "fooBar" {
		t.Error(expected("name", "fooBar", me.GetName()))
	}

	if me.GetCreated() != 123456789.0 {
		t.Error(expected("created", 123456789.0, me.GetCreated()))
	}

	if me.GetCreatedUtc() != 1315269998.0 {
		t.Error(expected("created_utc", 1315269998.0, me.GetCreatedUtc()))
	}

	if me.GetLinkKarma() != 31 {
		t.Error(expected("link_karma", 31, me.GetLinkKarma()))
	}

	if me.GetCommentKarma() != 557 {
		t.Error(expected("comment_karma", 557, me.GetCommentKarma()))
	}

	if !me.GetIsGold() {
		t.Error(expected("is_gold", true, false))
	}

	if me.GetGoldCredits() != 40 {
		t.Error(expected("gold_credits", 40, me.GetGoldCredits()))
	}

	if me.GetGoldExpiration() != 23273468.0 {
		t.Error(expected("gold_expiration", 23273468.0, me.GetGoldExpiration()))
	}

	if !me.GetOver_18() {
		t.Error(expected("over_18", true, false))
	}

	if !me.GetHasVerifiedEmail() {
		t.Error(expected("has_verified_email", true, false))
	}

	if !me.GetHideFromRobots() {
		t.Error(expected("hide_from_robots", true, false))
	}

	if me.GetId() != "5sryd" {
		t.Error(expected("id", "5sryd", me.GetId()))
	}

	if !me.GetHasModMail() {
		t.Error(expected("has_mod_mail", true, false))
	}
}

func TestMeKarma(t *testing.T) {
	server := newServerFromResponse(200, `{
		"data": [
			{
				"sr": "relationships",
				"comment_karma": 80,
				"link_karma": 60
			},
			{
				"sr": "self",
				"comment_karma": 30,
				"link_karma": 20
			}
		]
	}`)
	agent := &Agent{
		client: nface.NewClient(
			newProxyClient(server.URL),
			"test-client",
			server.URL),
	}
	karma, err := agent.MeKarma()
	if err != nil {
		t.Fatalf("failed to get self: %v", err)
	}

	if len(karma) != 2 {
		t.Error(expected("karma breakdown length", 2, len(karma)))
	}

	if karma[0].GetSr() != "relationships" {
		t.Error(expected("subreddit name", "relationships", karma[0].GetSr()))
	}

	if karma[0].GetCommentKarma() != 80 {
		t.Error(expected("comment karma", 80, karma[0].GetCommentKarma()))
	}

	if karma[0].GetLinkKarma() != 60 {
		t.Error(expected("link karma", 60, karma[0].GetLinkKarma()))
	}
}
