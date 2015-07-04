package rface

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type closer struct {
	io.Reader
}

func (c closer) Close() error {
	return nil
}

func TestHttpRequest(t *testing.T) {
	req := &Request{
		Action: POST,
		BasicAuthUser: "Robert",
		BasicAuthPass: "1234",
		Values: &url.Values{},
	}

	httpReq, err := req.httpRequest()
	if err != nil {
		t.Errorf("failed to generate http request: %v", err)
	}

	user, pass, ok := httpReq.BasicAuth()
	if !ok {
		t.Error("basic auth not set")
	}

	if user != "Robert" || pass != "1234" {
		t.Error("basic auth set incorrectly")
	}
}

func TestHttpRequestPostValues(t *testing.T) {
	vals := &url.Values{}
	vals.Add("food", "pancake")
	vals.Add("animal", "lynx")

	req := &Request{
		Action: POST,
		Values: vals,
	}

	httpReq, err := req.httpRequest()
	if err != nil {
		t.Error("failed to generate http request: %v", err)
	}

	body, err := ioutil.ReadAll(httpReq.Body)
	if err != nil {
		t.Error("request body not readable")
	}

	buf := bytes.NewBuffer(body)
	if buf.String() != vals.Encode() {
		t.Error(
			"POST data incorrect; expected %s, got %s",
			vals.Encode(),
			buf.String())
	}

	if httpReq.Header.Get("content-type") != contentType {
		t.Error("content type incorrect or unset")
	}
}

func TestHttpRequestGetValues(t *testing.T) {
	vals := &url.Values{}
	vals.Add("food", "pancake")
	vals.Add("animal", "lynx")

	req := &Request{
		Action: GET,
		Values: vals,
		OAuth: "token",
	}

	httpReq, err := req.httpRequest()
	if err != nil {
		t.Error("failed to generate http request: %v", err)
	}

	if httpReq.Body != nil {
		t.Error("body set in GET request")
	}

	if !strings.Contains(httpReq.URL.String(), vals.Encode()) {
		t.Error("GET values not written to url")
	}

	if httpReq.Header.Get("Authorization") != "bearer token" {
		t.Errorf(
			"http oauth wrong; expected bearer token, got %s",
			httpReq.Header.Get("Authorization"))
	}
}

func TestParseResponse(t *testing.T) {
	body := closer{bytes.NewBufferString(`{"id":"manning", "age":18}`)}
	dummy := struct {
		ID string `json:"id"`
		Age int `json:"age"`
	}{"jacob", 29}
	resp := &http.Response{Body: body}
	if err := parseResponse(resp, &dummy); err != nil {
		t.Errorf("parsing response failed: %v\n", err)
	}

	if dummy.ID != "manning" || dummy.Age != 18 {
		t.Error("field incorrectly unmarshaled")
	}
}
