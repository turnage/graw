package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestDo(t *testing.T) {
	expected := "expected"
	actual := ""
	serv := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				actual = r.Header.Get("User-Agent")
				fmt.Fprintf(w, "ok")
			},
		),
	)
	cli := &client{
		useragent: "expected",
		cli:       http.DefaultClient,
	}
	url, err := url.Parse(serv.URL)
	if err != nil {
		t.Fatalf("failed to parse test url: %v", err)
	}

	resp, err := cli.Do(&http.Request{URL: url})
	if err != nil {
		t.Fatalf("failed to execute request: %v")
	}
	if resp.StatusCode != 200 {
		t.Errorf("an anomaly occurred")
	}
	if actual != expected {
		t.Errorf("wanted %s; got %s", expected, actual)
	}
}
