package testutil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestNewReadCloser(t *testing.T) {
	expected := "internet"
	resp := &http.Response{Body: NewReadCloser(expected, nil)}
	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("failed to read body: %v", err)
	}

	actual := bytes.NewBuffer(buffer).String()
	if actual != expected {
		t.Errorf(
			"content not correct; got %v, wanted %v",
			actual,
			expected)
	}

	resp.Body = NewReadCloser(expected, fmt.Errorf("an error"))
	if _, err := ioutil.ReadAll(resp.Body); err == nil {
		t.Error("requested error not returned by Read() calls")
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
		Body:       bytesCloser{buffer: bytes.NewBuffer(expected)},
	}, 200, expected) {
		t.Error("failed to identify status code difference")
	}

	if ResponseIs(&http.Response{
		StatusCode: 200,
		Body: bytesCloser{
			buffer: bytes.NewBuffer(expected),
			err:    fmt.Errorf("AN ERROR"),
		},
	}, 200, expected) {
		t.Error("faulty read of response body did not become a diff")
	}

	if ResponseIs(&http.Response{
		StatusCode: 200,
		Body:       bytesCloser{buffer: bytes.NewBuffer(expected)},
	}, 200, []byte("sdfsdj")) {
		t.Error("body comparison failed; should have returned false")
	}

	if !ResponseIs(&http.Response{
		StatusCode: 200,
		Body:       bytesCloser{buffer: bytes.NewBuffer(expected)},
	}, 200, expected) {
		t.Error("body comparison failed; should have returned true")
	}
}
