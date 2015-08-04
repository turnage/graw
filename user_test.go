package graw

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/paytonturnage/graw/internal/auth"
	"github.com/paytonturnage/graw/internal/client"
	"github.com/paytonturnage/graw/internal/ratelimiter"
	"github.com/paytonturnage/graw/internal/testutil"
	"github.com/paytonturnage/redditproto"
)

func TestNewUser(t *testing.T) {
	agent := &redditproto.UserAgent{}
	if err := proto.UnmarshalText(`
		user_agent: "agent"
		client_id: "id"
		client_secret: "secret"
		username: "username"
		password: "password"
	`, agent); err != nil {
		t.Fatalf("failed to build expectation proto: %v", err)
	}
	expected := &User{
		agent: "agent",
		authorizer: auth.NewOAuth2Authorizer(
			"id",
			"secret",
			"username",
			"password"),
		client: nil,
		limiter: ratelimiter.NewTimeRateLimiter(
			time.Minute,
			maxQueriesPerMinute,
		),
	}
	actual := NewUser(agent)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf(
			"user built incorrectly; got %v, wanted %v",
			actual,
			expected)
	}

}

func TestNewUserFromFile(t *testing.T) {
	if _, err := NewUserFromFile("not_a_real_file"); err == nil {
		t.Error("bad filename did not return error")
	}

	emptyFile, err := ioutil.TempFile("", "empty_user_agent")
	if err != nil {
		t.Fatalf("could not create temporary file: %v", err)
	}

	if _, err := NewUserFromFile(emptyFile.Name()); err == nil {
		t.Error("failed to err with bad user agent file")
	}

	expectedAgent := &redditproto.UserAgent{}
	if err := proto.UnmarshalText(`
		user_agent: "agent"
		client_id: "id"
		client_secret: "secret"
		username: "username"
		password: "password"
	`, expectedAgent); err != nil {
		t.Errorf("failed to build expectation proto: %v", err)
	}
	expected := &User{
		agent: "agent",
		authorizer: auth.NewOAuth2Authorizer(
			"id",
			"secret",
			"username",
			"password"),
		client: nil,
		limiter: ratelimiter.NewTimeRateLimiter(
			time.Minute,
			maxQueriesPerMinute,
		),
	}

	agentFile, err := ioutil.TempFile("", "user_agent")
	if err != nil {
		t.Fatalf("could not create temporary file: %v", err)
	}

	_, err = agentFile.WriteString(proto.MarshalTextString(expectedAgent))
	if err != nil {
		t.Errorf("could not write to user_agent file: %v", err)
	}

	actual, err := NewUserFromFile(agentFile.Name())
	if err != nil {
		t.Errorf("could not build user from file: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"user incorrect; expected %v, got %v",
			expected,
			actual)
	}
}

func TestAuth(t *testing.T) {
	user := &User{}
	user.authorizer = auth.NewMockAuth(nil, fmt.Errorf("BAD THING"))
	if err := user.Auth(); err == nil {
		t.Error("auth did not return an error on failure")
	}

	user.authorizer = auth.NewMockAuth(http.DefaultClient, nil)
	if err := user.Auth(); err != nil {
		t.Error("auth returned error on success")
	}
	if user.client != http.DefaultClient {
		t.Error("client not set upon successful auth")
	}
}

func TestExec(t *testing.T) {
	user := &User{}
	req := &http.Request{}
	expected := &struct {
		Key string
	}{Key: "value"}
	actual := &struct {
		Key string `"json:"key,omitempty"`
	}{}
	jsonAgent := `{"key": "value"}`

	user.client = client.NewMockClient(
		&http.Response{
			StatusCode: 200,
			Body:       testutil.NewReadCloser(jsonAgent, nil),
		},
		fmt.Errorf("A BAD THING HAPPENED"),
	)
	if err := user.Exec(req, actual); err == nil {
		t.Error("failed request did not return an error")
	}

	user.client = client.NewMockClient(
		&http.Response{
			StatusCode: 200,
			Body: testutil.NewReadCloser(
				jsonAgent,
				fmt.Errorf("misbehavior bad stuff")),
		},
		nil,
	)
	if err := user.Exec(req, actual); err == nil {
		t.Error("corrupt body did not return an error")
	}

	user.client = client.NewMockClient(
		&http.Response{
			StatusCode: 201,
			Body:       testutil.NewReadCloser(jsonAgent, nil),
		},
		nil,
	)
	if err := user.Exec(req, actual); err == nil {
		t.Error("bad status code did not return an error")
	}

	user.client = client.NewMockClient(
		&http.Response{
			StatusCode: 200,
			Body:       nil,
		},
		nil,
	)
	if err := user.Exec(req, actual); err == nil {
		t.Error("nil body did not return an error")
	}

	user.client = client.NewMockClient(
		&http.Response{
			StatusCode: 200,
			Body:       testutil.NewReadCloser(jsonAgent, nil),
		},
		nil,
	)
	err := user.Exec(req, actual)
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

func TestMe(t *testing.T) {
	acct := `{"name":"Rob"}`
	user := &User{}
	user.client = client.NewMockClient(
		&http.Response{
			StatusCode: 200,
			Body:       testutil.NewReadCloser(acct, nil),
		},
		nil,
	)
	actual, err := user.Me()
	if err != nil {
		t.Fatalf("getting self failed: %v", err)
	}
	if actual.GetName() != "Rob" {
		t.Error("name not extracted from response")
	}
}

func TestScrape(t *testing.T) {
	listing := `{
		"data": {
			"children": [
				{"data": {"title": "1"}},
				{"data": {"title": "2"}}
			]
		}
	}`
	user := &User{}
	user.client = client.NewMockClient(
		&http.Response{
			StatusCode: 200,
			Body:       testutil.NewReadCloser(listing, nil),
		},
		nil,
	)

	actualListing, err := user.Scrape("relationships", "hot", "", "", 3)
	if err != nil {
		t.Fatalf("failed to scrape: %v", err)
	}

	if len(actualListing) != 2 {
		t.Errorf(
			"unexpected listing length; got %d, wanted 2",
			len(actualListing))
	}

	if actualListing[0].GetTitle() != "1" {
		t.Errorf("first title incorrect; link: %v", actualListing[0])
	}

	if actualListing[1].GetTitle() != "2" {
		t.Errorf("second title incorrect; link: %v", actualListing[0])
	}
}
