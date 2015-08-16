package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestDo(t *testing.T) {
	expected := &struct {
		Key string
	}{Key: "value"}
	actual := &struct {
		Key string `"json:"key,omitempty"`
	}{}
	cli := &client{cli: http.DefaultClient}
	responseCode := 200
	responseBody := `{"key": "value"}`

	serv := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(responseCode)
				fmt.Fprintf(w, responseBody)
			},
		),
	)
	addr, err := url.Parse(serv.URL)
	if err != nil {
		t.Fatalf("failed to parse test server url")
	}
	req := &http.Request{URL: addr}

	if err := cli.Do(req, actual); err != nil {
		t.Errorf("exec failed: %v", err)
	}
	if actual.Key != expected.Key {
		t.Errorf(
			"response incorrect; got %v, wanted %v",
			actual,
			expected)
	}

	responseCode = 404
	if err := cli.Do(req, actual); err == nil {
		t.Error("bad status code did not return an error")
	}

	addr, err = url.Parse("http://notarealurl")
	if err != nil {
		t.Fatalf("failed to parse test server url")
	}

	responseCode = 200
	req = &http.Request{URL: addr}
	if err := cli.Do(req, actual); err == nil {
		t.Error("error in request did not return an error")
	}
}

func TestDoRaw(t *testing.T) {
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

	resp, err := cli.doRaw(&http.Request{URL: url})
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
