package reddit

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

type mockParser struct {
	comments []*Comment
	posts    []*Post
	messages []*Message
}

func (m *mockParser) parse(
	blob json.RawMessage,
) ([]*Comment, []*Post, []*Message, error) {
	return m.comments, m.posts, m.messages, nil
}

func parserWhich(h harvest) parser {
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

func newMockClient() client {
	return &mockClient{}
}

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
		expected := harvest{
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

		harvest, err := r.reap(test.path, test.values)
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
