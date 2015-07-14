package graw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/paytonturnage/graw/data"
	"github.com/paytonturnage/graw/nface"
)

// newProxyClient returns an http.Client that redirects all requests to the
// redirect url.
func newProxyClient(redirect string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: func(r *http.Request) (*url.URL, error) {
				return url.Parse(redirect)
			},
		},
	}
}

// newServerFromResponse returns an httptest.Server that always responds with
// response and status.
func newServerFromResponse(status int, response []byte) *httptest.Server {
	responseString := bytes.NewBuffer(response).String()
	responseWriter := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Header().Set("content-type", "application/json")
		fmt.Fprintln(w, responseString)
	}
	return httptest.NewServer(http.HandlerFunc(responseWriter))
}

func TestNewUserAgentFromtFile(t *testing.T) {
	expected := &data.UserAgent{}
	if err := proto.UnmarshalText(`
		user_agent: "test"
		client_id: "id"
		client_secret: "secret"
		username: "user"
		password: "1234"
	`, expected); err != nil {
		t.Errorf("could not build expectation proto: %v", err)
	}

	agentFile, err := ioutil.TempFile("", "user_agent")
	if err != nil {
		t.Errorf("could not make user_agent file: %v", err)
	}

	_, err = agentFile.WriteString(proto.MarshalTextString(expected))
	if err != nil {
		t.Errorf("could not write to user_agent file: %v", err)
	}

	actual, err := newUserAgentFromFile(agentFile.Name())
	if err != nil {
		t.Errorf("could not build user agent from file: %v", err)
	}

	if !proto.Equal(expected, actual) {
		t.Errorf(
			"user agent incorrect; expected %v, got %v",
			expected,
			actual)
	}
}

func TestNewOAuthClient(t *testing.T) {
	server := newServerFromResponse(200, []byte(`{
			"access_token": "sjkhefwhf383nfjkf",
			"token_type": "bearer",
			"expires_in": 3600
			"scope": "*",
			"refresh_token": "akjfbkfjhksdjhf"
	}`))

	client, err := newOAuthClient(
		newUserAgent("test", "id", "secret", "user", "pass"), server.URL)
	if err != nil {
		t.Errorf("failed to make oauth client: %v", err)
	}

	if client == nil {
		t.Error("client not returned")
	}
}

func TestMe(t *testing.T) {
	resp, err := json.Marshal(&data.Account{})
	if err != nil {
		t.Fatalf("preparing response failed: %v", err)
	}
	server := newServerFromResponse(200, resp)
	agent := &Graw{
		client: nface.NewClient(
			newProxyClient(server.URL),
			"test-client",
			server.URL),
	}
	if _, err := agent.Me(); err != nil {
		t.Fatalf("failed to get self: %v", err)
	}
}

func TestMeKarma(t *testing.T) {
	resp, err := json.Marshal(&data.KarmaList{})
	if err != nil {
		t.Fatalf("preparing response failed: %v", err)
	}
	server := newServerFromResponse(200, resp)
	agent := &Graw{
		client: nface.NewClient(
			newProxyClient(server.URL),
			"test-client",
			server.URL),
	}
	if _, err := agent.MeKarma(); err != nil {
		t.Fatalf("failed to get karma: %v", err)
	}
}
