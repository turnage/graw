package testutil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

// bytesCloser implements io.ReadCloser around bytes.Buffer. This makes it
// easier to inject expected body content into http.Response structs.
type bytesCloser struct {
	*bytes.Buffer
}

func (b bytesCloser) Close() error {
	return nil
}

func TestNewProxyClient(t *testing.T) {
	expectedString := "sample response"
	writeResponse := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, expectedString)
	}
	server := httptest.NewServer(http.HandlerFunc(writeResponse))
	client := NewProxyClient(server.URL)

	dummyURL, err := url.Parse("http://www.google.com")
	if err != nil {
		t.Fatalf("failed to create url: %v", err)
	}

	resp, err := client.Do(&http.Request{URL: dummyURL})
	if err != nil {
		t.Fatalf("failed request to proxy server: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("request returned bad status: %d", resp.StatusCode)
	}

	if resp.Body == nil {
		t.Fatalf("no body in response")
	}

	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	expected := []byte(expectedString)
	if !reflect.DeepEqual(respBytes, expected) {
		t.Errorf("response incorrect; expected %s, got %s", respBytes, expected)
	}
}

func TestNewServerFromResponse(t *testing.T) {
	expected := []byte("10101010101___")
	server := NewServerFromResponse(200, expected)

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("request to server failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("request returned bad status: %d", resp.StatusCode)
	}

	if resp.Body == nil {
		t.Fatalf("no body in response")
	}

	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if !reflect.DeepEqual(respBytes, expected) {
		t.Errorf("response incorrect; expected %s, got %s", respBytes, expected)
	}
}

func TestResponseIs(t *testing.T) {
	expected := []byte("ksjdnksbf")
	if ResponseIs(&http.Response{StatusCode: 200, Body: nil}, 200, expected) {
		t.Error("failed to identify nil body")
	}

	if !ResponseIs(&http.Response{StatusCode: 200, Body: nil}, 200, nil) {
		t.Error("failed to accept nil body with nil expectation")
	}

	if ResponseIs(&http.Response{
		StatusCode: 201,
		Body:       bytesCloser{Buffer: bytes.NewBuffer(expected)},
	}, 200, expected) {
		t.Error("failed to identify status code difference")
	}

	if ResponseIs(&http.Response{
		StatusCode: 200,
		Body:       bytesCloser{Buffer: bytes.NewBuffer(expected)},
	}, 200, []byte("sdfsdj")) {
		t.Error("body comparison failed; should have returned false")
	}

	if !ResponseIs(&http.Response{
		StatusCode: 200,
		Body:       bytesCloser{Buffer: bytes.NewBuffer(expected)},
	}, 200, expected) {
		t.Error("body comparison failed; should have returned true")
	}
}
