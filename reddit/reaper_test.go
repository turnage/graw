package reddit

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestNew(t *testing.T) {
	cli := &mockClient{}
	par := &mockParser{}
	cfg := reaperConfig{
		client:   cli,
		parser:   par,
		hostname: "com",
		tls:      true,
	}
	expected := &reaperImpl{
		cli:      cli,
		parser:   par,
		hostname: "com",
		scheme:   "https",
	}

	if diff := pretty.Compare(newReaper(cfg), expected); diff != "" {
		t.Errorf("reaper construction incorrect; diff: %s", diff)
	}
}

func Testreap(t *testing.T) {
	for i, test := range []struct {
		path    string
		values  map[string]string
		correct http.Request
	}{
		{"", nil, http.Request{
			Method: "GET",
			Host:   "com",
			URL: &url.URL{
				Scheme:   "http",
				Host:     "com",
				Path:     "",
				RawQuery: "",
			},
		}},
		{"", map[string]string{"key": "value"}, http.Request{
			Method: "GET",
			Host:   "com",
			URL: &url.URL{
				Scheme:   "http",
				Host:     "com",
				Path:     "",
				RawQuery: "key=value",
			},
		}},
		{"path", nil, http.Request{
			Method: "GET",
			Host:   "com",
			URL: &url.URL{
				Scheme:   "http",
				Host:     "com",
				Path:     "path",
				RawQuery: "",
			},
		}},
	} {
		expected := Harvest{
			Comments: []*Comment{
				&Comment{
					Body: "comment",
				},
			},
			Posts: []*Post{
				&Post{
					SelfText: "post",
				},
			},
			Messages: []*Message{
				&Message{
					Body: "message",
				},
			},
		}
		c := &mockClient{}
		r := &reaperImpl{
			cli:      c,
			parser:   parserWhich(expected),
			hostname: "com",
			scheme:   "http",
		}

		Harvest, err := r.reap(test.path, test.values)
		if err != nil {
			t.Errorf("Error reaping input %d: %v", i, err)
		}

		if diff := pretty.Compare(Harvest, expected); diff != "" {
			t.Errorf("Harvest incorrect; diff: %s", diff)
		}

		if diff := pretty.Compare(c.request, test.correct); diff != "" {
			t.Errorf("request incorrect; diff: %s", diff)
		}
	}
}

func Testsow(t *testing.T) {
	for i, test := range []struct {
		path    string
		values  map[string]string
		correct http.Request
	}{
		{"", nil, http.Request{
			Method: "POST",
			Header: formEncoding,
			Host:   "com",
			URL: &url.URL{
				Scheme:   "http",
				Host:     "com",
				Path:     "",
				RawQuery: "",
			},
		}},
		{"", map[string]string{"key": "value"}, http.Request{
			Method: "POST",
			Header: formEncoding,
			Host:   "com",
			URL: &url.URL{
				Scheme:   "http",
				Host:     "com",
				Path:     "",
				RawQuery: "key=value",
			},
		}},
		{"path", nil, http.Request{
			Method: "POST",
			Header: formEncoding,
			Host:   "com",
			URL: &url.URL{
				Scheme:   "http",
				Host:     "com",
				Path:     "path",
				RawQuery: "",
			},
		}},
	} {
		c := &mockClient{}
		r := &reaperImpl{
			cli:      c,
			parser:   &mockParser{},
			hostname: "com",
			scheme:   "http",
		}

		if err := r.sow(test.path, test.values); err != nil {
			t.Errorf("Error reaping input %d: %v", i, err)
		}

		if diff := pretty.Compare(c.request, test.correct); diff != "" {
			t.Errorf("request incorrect; diff: %s", diff)
		}
	}
}
