package request

import (
	"net/url"
	"reflect"
	"testing"
)

func TestNewGet(t *testing.T) {
	expectedValues := &url.Values{
		"key": []string{
			"value1",
			"value2",
		},
	}
	req, err := New("GET", "http://web.com", expectedValues)
	if err != nil {
		t.Errorf("failed to build request: %v", err)
	}

	if err := req.ParseForm(); err != nil {
		t.Errorf("failed to parse form values: %v", err)
	}

	if !reflect.DeepEqual(req.Form, *expectedValues) {
		t.Errorf(
			"values incorrect; got %v, wanted %v",
			req.Form,
			*expectedValues)
	}
}

func TestNewPost(t *testing.T) {
	if _, err := New("POST", "http://web.com", nil); err == nil {
		t.Error("nil values for post did not return an error")
	}

	expectedValues := &url.Values{
		"key": []string{
			"value1",
			"value2",
		},
	}
	req, err := New("POST", "http://web.com", expectedValues)
	if err != nil {
		t.Error("failed to build request: %v", err)
	}
	if err := req.ParseForm(); err != nil {
		t.Errorf("failed to parse form values: %v", err)
	}
	if !reflect.DeepEqual(req.PostForm, *expectedValues) {
		t.Errorf(
			"values incorrect; got %v, wanted %v",
			req.PostForm,
			*expectedValues)
	}

	if _, err := New("POST", ":badurl", expectedValues); err == nil {
		t.Error("bad url did not return error")
	}
}

func TestNewUnsupported(t *testing.T) {
	if _, err := New("NOT", "http://web.com", nil); err == nil {
		t.Error("unsupported method did not return an error")
	}
}
