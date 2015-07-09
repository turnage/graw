package nface

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestBuildPost(t *testing.T) {
	expectedUserAgent := "test"
	client := &Client{userAgent: expectedUserAgent}
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
		t.Error("bad POST body; expected %s, got %s", expectedBody, actualBody)
	}

	actualContentType := req.Header.Get("content-type")
	if req.Header.Get("content-type") != contentType {
		t.Error(
			"bad content-type; expected %s, got %s",
			contentType,
			actualContentType)
	}

	actualUserAgent := req.Header.Get("user-agent")
	if actualUserAgent != expectedUserAgent {
		t.Error(
			"bad user-agent; expected %s, got %s",
			expectedUserAgent,
			actualUserAgent)
	}
}

func TestBuildGet(t *testing.T) {
	expectedUserAgent := "test"
	client := &Client{userAgent: expectedUserAgent}
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

func TestSend(t *testing.T) {
	expectedResponse := "sample response"
	writeResponse := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, expectedResponse)
	}
	serv := httptest.NewServer(http.HandlerFunc(writeResponse))
	client := &Client{client: http.DefaultClient}

	url, err := url.Parse(serv.URL)
	if err != nil {
		t.Errorf("failed to parse test server url: %v", err)
	}

	buf, err := client.doRequest(&http.Request{URL: url})
	if err != nil {
		t.Errorf("failed to send request: %v", err)
	}

	if buf == nil {
		t.Error("failed to extract response; body is nil")
	}

	actualResponse := bytes.NewBuffer(buf).String()
	if actualResponse != expectedResponse {
		t.Errorf(
			"bad response; expected %s, got %s",
			expectedResponse,
			actualResponse)
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
