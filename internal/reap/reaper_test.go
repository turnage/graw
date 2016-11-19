package reap

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/kylelemons/godebug/pretty"

	"github.com/turnage/graw/internal/client"
	"github.com/turnage/graw/reddit"
)

type mockParser struct {
	comments []*reddit.Comment
	posts    []*reddit.Post
	messages []*reddit.Message
}

func (m *mockParser) Parse(
	blob json.RawMessage,
) ([]*reddit.Comment, []*reddit.Post, []*reddit.Message, error) {
	return m.comments, m.posts, m.messages, nil
}

func parserWhich(h Harvest) reddit.Parser {
	return &mockParser{
		comments: h.Comments,
		posts:    h.Posts,
		messages: h.Messages,
	}
}

type mockClient struct {
	request *http.Request
}

func (m *mockClient) Do(r *http.Request) ([]byte, error) {
	m.request = r
	return nil, nil
}

func newMockClient() client.Client {
	return &mockClient{}
}

func TestNew(t *testing.T) {
	cli := &mockClient{}
	par := &mockParser{}
	cfg := Config{
		Client:   cli,
		Parser:   par,
		Hostname: "reddit.com",
		TLS:      true,
	}
	expected := &reaper{
		cli:      cli,
		parser:   par,
		hostname: "reddit.com",
		scheme:   "https",
	}

	if diff := pretty.Compare(New(cfg), expected); diff != "" {
		t.Errorf("reaper construction incorrect; diff: %s", diff)
	}
}

func TestReap(t *testing.T) {
	for i, test := range []struct {
		path    string
		values  map[string]string
		correct http.Request
	}{
		{"", nil, http.Request{
			Method: "GET",
			Host:   "reddit.com",
			URL: &url.URL{
				Scheme:   "http",
				Host:     "reddit.com",
				Path:     "",
				RawQuery: "",
			},
		}},
		{"", map[string]string{"key": "value"}, http.Request{
			Method: "GET",
			Host:   "reddit.com",
			URL: &url.URL{
				Scheme:   "http",
				Host:     "reddit.com",
				Path:     "",
				RawQuery: "key=value",
			},
		}},
		{"path", nil, http.Request{
			Method: "GET",
			Host:   "reddit.com",
			URL: &url.URL{
				Scheme:   "http",
				Host:     "reddit.com",
				Path:     "path",
				RawQuery: "",
			},
		}},
	} {
		expected := Harvest{
			Comments: []*reddit.Comment{
				&reddit.Comment{
					Body: "comment",
				},
			},
			Posts: []*reddit.Post{
				&reddit.Post{
					SelfText: "post",
				},
			},
			Messages: []*reddit.Message{
				&reddit.Message{
					Body: "message",
				},
			},
		}
		c := &mockClient{}
		r := &reaper{
			cli:      c,
			parser:   parserWhich(expected),
			hostname: "reddit.com",
			scheme:   "http",
		}

		harvest, err := r.Reap(test.path, test.values)
		if err != nil {
			t.Errorf("Error reaping input %d: %v", i, err)
		}

		if diff := pretty.Compare(harvest, expected); diff != "" {
			t.Errorf("harvest incorrect; diff: %s", diff)
		}

		if diff := pretty.Compare(c.request, test.correct); diff != "" {
			t.Errorf("request incorrect; diff: %s", diff)
		}
	}
}

func TestSow(t *testing.T) {
	for i, test := range []struct {
		path    string
		values  map[string]string
		correct http.Request
	}{
		{"", nil, http.Request{
			Method: "POST",
			Header: formEncoding,
			Host:   "reddit.com",
			URL: &url.URL{
				Scheme:   "http",
				Host:     "reddit.com",
				Path:     "",
				RawQuery: "",
			},
		}},
		{"", map[string]string{"key": "value"}, http.Request{
			Method: "POST",
			Header: formEncoding,
			Host:   "reddit.com",
			URL: &url.URL{
				Scheme:   "http",
				Host:     "reddit.com",
				Path:     "",
				RawQuery: "key=value",
			},
		}},
		{"path", nil, http.Request{
			Method: "POST",
			Header: formEncoding,
			Host:   "reddit.com",
			URL: &url.URL{
				Scheme:   "http",
				Host:     "reddit.com",
				Path:     "path",
				RawQuery: "",
			},
		}},
	} {
		c := &mockClient{}
		r := &reaper{
			cli:      c,
			parser:   &mockParser{},
			hostname: "reddit.com",
			scheme:   "http",
		}

		if err := r.Sow(test.path, test.values); err != nil {
			t.Errorf("Error reaping input %d: %v", i, err)
		}

		if diff := pretty.Compare(c.request, test.correct); diff != "" {
			t.Errorf("request incorrect; diff: %s", diff)
		}
	}
}
