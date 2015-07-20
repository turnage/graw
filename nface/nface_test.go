package nface

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/paytonturnage/graw/data"
	"github.com/paytonturnage/graw/testutil"
)

func TestBuildPost(t *testing.T) {
	expectedAgent := "test"
	userAgent := &data.UserAgent{UserAgent: &expectedAgent}
	client := &Client{userAgent: userAgent}
	vals := &url.Values{
		"food":   []string{"pancake"},
		"animal": []string{"lynx"},
	}
	req, err := client.buildRequest(&Request{
		Action: POST,
		Values: vals,
	})

	if err != nil {
		t.Errorf("failed to build http request: %v", err)
	}

	if req == nil {
		t.Fatal("returned http.Request is nil")
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Error("request body not readable")
	}

	expectedBody := vals.Encode()
	actualBody := bytes.NewBuffer(body).String()
	if actualBody != expectedBody {
		t.Errorf("bad POST body; expected %s, got %s", expectedBody, actualBody)
	}

	actualContentType := req.Header.Get("content-type")
	if req.Header.Get("content-type") != contentType {
		t.Errorf(
			"bad content-type; expected %s, got %s",
			contentType,
			actualContentType)
	}

	actualUserAgent := req.Header.Get("user-agent")
	if actualUserAgent != expectedAgent {
		t.Errorf(
			"bad user-agent; expected %s, got %s",
			expectedAgent,
			actualUserAgent)
	}
}

func TestBuildGet(t *testing.T) {
	client := &Client{}
	vals := &url.Values{
		"food":   []string{"pancake"},
		"animal": []string{"lynx"},
	}
	req, err := client.buildRequest(&Request{
		Action: GET,
		Values: vals,
	})

	if err != nil {
		t.Errorf("failed to build http request: %v", err)
	}

	if req == nil {
		t.Fatal("returned http.Request is nil")
	}

	if req.Body != nil {
		t.Error("GET request has body")
	}

	if !strings.Contains(req.URL.String(), vals.Encode()) {
		t.Errorf("GET values not written to url: %s", req.URL.String())
	}
}

func TestBuildPostNilValues(t *testing.T) {
	client := &Client{}
	req, err := client.buildRequest(&Request{Action: POST})
	if err == nil {
		t.Error("no error on nil values in POST request")
	}
	if req != nil {
		t.Error("returned http.Request is not nil")
	}
}

func TestBuildGetNilValues(t *testing.T) {
	client := &Client{}
	req, err := client.buildRequest(&Request{Action: GET})
	if err != nil {
		t.Error("error on nil values in GET request")
	}
	if req == nil {
		t.Error("returned http.Request is nil")
	}
}

func TestDo(t *testing.T) {
	expected := &data.UserAgent{}
	if err := proto.UnmarshalText(`
		user_agent: "agent"
		client_id: "id"
		client_secret: "secret"
		username: "username"
		password: "password"
	`, expected); err != nil {
		t.Fatalf("failed to build expectation proto: %v", err)
	}
	serv, _ := testutil.NewServerFromResponse(200, []byte(`{
		"user_agent": "agent",
		"client_id": "id",
		"client_secret": "secret",
		"username": "username",
		"password": "password"
	}`))

	actual := &data.UserAgent{}
	client := &Client{client: http.DefaultClient}
	if err := client.Do(&Request{URL: serv.URL}, actual); err != nil {
		t.Errorf("executing request failed: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", actual, expected)
	}
}

func TestOAuth(t *testing.T) {
	serv, _ := testutil.NewServerFromResponse(200, []byte(`{
			"access_token": "sjkhefwhf383nfjkf",
			"token_type": "bearer",
			"expires_in": 3600
			"scope": "*",
			"refresh_token": "akjfbkfjhksdjhf"
	}`))
	userAgent := data.NewUserAgent("test", "id", "secret", "user", "pass")
	client := &Client{userAgent: userAgent}
	if err := client.oauth(serv.URL); err != nil {
		t.Errorf("failed to make oauth client: %v", err)
	}

	if client == nil {
		t.Error("client not returned")
	}
}

func TestSend(t *testing.T) {
	expected := []byte("sample response")
	serv, _ := testutil.NewServerFromResponse(200, expected)
	client := &Client{client: http.DefaultClient}

	url, err := url.Parse(serv.URL)
	if err != nil {
		t.Errorf("failed to parse test server url: %v", err)
	}

	actual, err := client.doRequest(&http.Request{URL: url})
	if err != nil {
		t.Errorf("failed to send request: %v", err)
	}

	if actual == nil {
		t.Error("failed to extract response; body is nil")
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("bad response; expected %s, got %s", expected, actual)
	}
}

func TestSendError(t *testing.T) {
	makeError := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "an error", http.StatusInternalServerError)
	}
	serv := httptest.NewServer(http.HandlerFunc(makeError))
	client := &Client{client: http.DefaultClient}

	url, err := url.Parse(serv.URL)
	if err != nil {
		t.Errorf("failed to parse test server url: %v", err)
	}

	actualResponse, err := client.doRequest(&http.Request{URL: url})
	if err == nil {
		t.Error("no error on server error")
	}

	if actualResponse != nil {
		t.Errorf("error from unexpected path; response: %s", actualResponse)
	}
}
