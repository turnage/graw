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
