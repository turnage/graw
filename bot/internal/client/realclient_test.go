package client

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestExec(t *testing.T) {
	cli := &client{cli: http.DefaultClient}
	responseCode := 200
	expectedBody := `{"key": "value"}`

	serv := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(responseCode)
				fmt.Fprintf(w, expectedBody)
			},
		),
	)
	addr, err := url.Parse(serv.URL)
	if err != nil {
		t.Fatalf("failed to parse test server url")
	}
	req := &http.Request{URL: addr}

	body, err := cli.exec(req)
	if err != nil {
		t.Errorf("exec failed: %v", err)
	}
	bodyBuffer := new(bytes.Buffer)
	_, err = bodyBuffer.ReadFrom(body)
	if err != nil {
		t.Errorf("failed to ready response body: %v", err)
	}

	bodyString := bodyBuffer.String()
	if bodyString != expectedBody {
		t.Errorf("got %s; wanted %s", bodyString, expectedBody)
	}

	responseCode = 404
	if _, err := cli.exec(req); err == nil {
		t.Error("bad status code did not return an error")
	}

	addr, err = url.Parse("http://notarealurl")
	if err != nil {
		t.Fatalf("failed to parse test server url")
	}

	responseCode = 200
	req = &http.Request{URL: addr}
	if _, err := cli.exec(req); err == nil {
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
		agent: "expected",
		cli:   http.DefaultClient,
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
