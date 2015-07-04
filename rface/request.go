package rface

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

type ReqAction int

const (
	GET = iota
	POST = iota
)

const (
	contentType = "application/x-www-form-urlencoded"
)

// Request describes how to build an http.Request for the reddit api.
type Request struct {
	// Action is the request type (e.g. "POST" or "GET").
	Action ReqAction
	// BasicAuthUser is the username to provide to basic auth prompts.
	BasicAuthUser string
	// BasicAuthPass is the password to provide to basic auth prompts.
	BasicAuthPass string
	// BaseUrl is the url of the api call, which values will be appended to.
	BaseUrl string
	// OAuth is the requests' oauth access token, which goes in the header.
	OAuth string
	// Values holds any parameters for the api call; encoded in url.
	Values *url.Values
}

// Do executes the request and writes the JSON response to the val interface.
// See the godoc on json.Unmarshal for information on what to provide as val
// and how to set it up for parsing.
func (r *Request) Do(val interface{}) error {
	req, err := r.httpRequest()
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return parseResponse(resp, val)
}

// httpRequest generates an http.Request from a Request struct.
func (r *Request) httpRequest() (*http.Request, error) {
	var req *http.Request
	var err error
	if r.Action == GET {
		req, err = getRequest(r.BaseUrl, r.Values)
	} else if r.Action == POST {
		req, err = postRequest(r.BaseUrl, r.Values)
	}
	if err != nil {
		return nil, err
	}

	if r.BasicAuthUser != "" && r.BasicAuthPass != "" {
		req.SetBasicAuth(r.BasicAuthUser, r.BasicAuthPass)
	} else if r.OAuth != "" {
		req.Header.Add(
			"Authorization",
			fmt.Sprintf("bearer %s", r.OAuth))
	}

	return req, nil
}

func postRequest(url string, vals *url.Values) (*http.Request, error) {
	reqBody := bytes.NewBufferString(vals.Encode())
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", contentType)
	return req, nil
}

func getRequest(url string, vals *url.Values) (*http.Request, error) {
	reqUrl := fmt.Sprintf("%s?%s", url, vals.Encode())
	return http.NewRequest("GET", reqUrl, nil)
}
