package graw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/paytonturnage/graw/data"
	"github.com/paytonturnage/graw/nface"
)

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
func newServerFromResponse(status int, response []byte) *httptest.Server {
	responseString := bytes.NewBuffer(response).String()
	responseWriter := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Header().Set("content-type", "application/json")
		fmt.Fprintln(w, responseString)
	}
	return httptest.NewServer(http.HandlerFunc(responseWriter))
}

// Tests the me api call and that all response data is properly represented.
func TestMe(t *testing.T) {
	resp, err := json.Marshal(&data.Account{})
	if err != nil {
		t.Fatalf("preparing response failed: %v", err)
	}
	server := newServerFromResponse(200, resp)
	agent := &Agent{
		client: nface.NewClient(
			newProxyClient(server.URL),
			"test-client",
			server.URL),
	}
	if _, err := agent.Me(); err != nil {
		t.Fatalf("failed to get self: %v", err)
	}
}

func TestMeKarma(t *testing.T) {
	resp, err := json.Marshal(&data.KarmaList{})
	if err != nil {
		t.Fatalf("preparing response failed: %v", err)
	}
	server := newServerFromResponse(200, resp)
	agent := &Agent{
		client: nface.NewClient(
			newProxyClient(server.URL),
			"test-client",
			server.URL),
	}
	if _, err := agent.MeKarma(); err != nil {
		t.Fatalf("failed to get karma: %v", err)
	}
}
