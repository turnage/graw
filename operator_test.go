package graw

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

func TestExec(t *testing.T) {
	expected := &struct {
		Key string
	}{Key: "value"}
	actual := &struct {
		Key string `"json:"key,omitempty"`
	}{}
	jsonAgent := `{"key": "value"}`

	if err := exec(&mockClient{
		&http.Response{
			StatusCode: 200,
			Body: &bytesCloser{
				buffer: bytes.NewBufferString(jsonAgent),
				err:    nil,
			},
		},
		fmt.Errorf("A BAD THING HAPPENED"),
	}, &http.Request{}, actual); err == nil {
		t.Error("failed &http.Request{}uest did not return an error")
	}

	if err := exec(&mockClient{
		&http.Response{
			StatusCode: 200,
			Body: &bytesCloser{
				buffer: bytes.NewBufferString(jsonAgent),
				err:    fmt.Errorf("misbehavior bad stuff"),
			},
		},
		nil,
	}, &http.Request{}, actual); err == nil {
		t.Error("corrupt body did not return an error")
	}

	if err := exec(&mockClient{
		&http.Response{
			StatusCode: 201,
			Body: &bytesCloser{
				buffer: bytes.NewBufferString(jsonAgent),
				err:    nil,
			},
		},
		nil,
	}, &http.Request{}, actual); err == nil {
		t.Error("bad status code did not return an error")
	}

	if err := exec(&mockClient{
		&http.Response{
			StatusCode: 200,
			Body:       nil,
		},
		nil,
	}, &http.Request{}, actual); err == nil {
		t.Error("nil body did not return an error")
	}

	err := exec(&mockClient{
		&http.Response{
			StatusCode: 200,
			Body: &bytesCloser{
				buffer: bytes.NewBufferString(jsonAgent),
				err:    nil,
			},
		},
		nil,
	}, &http.Request{}, actual)
	if err != nil {
		t.Errorf("exec failed: %v", err)
	}
	if actual.Key != expected.Key {
		t.Errorf(
			"response incorrect; got %v, wanted %v",
			actual,
			expected)
	}
}

func TestScrape(t *testing.T) {
	listingJSON := `{
		"data": {
			"children": [
				{"data": {"title": "1"}},
				{"data": {"title": "2"}}
			]
		}
	}`
	listing, err := scrape(&mockClient{
		&http.Response{
			StatusCode: 200,
			Body: &bytesCloser{
				bytes.NewBufferString(listingJSON),
				nil,
			},
		},
		nil,
	}, "relationships", "hot", "", "", 3)
	if err != nil {
		t.Fatalf("failed to scrape: %v", err)
	}

	if len(listing) != 2 {
		t.Errorf(
			"unexpected listing length; got %d, wanted 2",
			len(listing))
	}

	if listing[0].GetTitle() != "1" {
		t.Errorf("first title incorrect; link: %v", listing[0])
	}

	if listing[1].GetTitle() != "2" {
		t.Errorf("second title incorrect; link: %v", listing[0])
	}
}
