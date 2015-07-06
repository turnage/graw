// Pacakge nface handles all communication between Go code and the Reddit api.
package nface

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ReqAction int

const (
	GET  = iota
	POST = iota
)

const (
	// contentType is a header flag for POST requests so the reddit api
	// knows how to read the request body.
	contentType = "application/x-www-form-urlencoded"
)

// Request describes how to build an http.Request for the reddit api.
type Request struct {
	// Action is the request type (e.g. "POST" or "GET").
	Action ReqAction
	// BaseUrl is the url of the api call, which values will be appended to.
	BaseUrl string
	// Values holds any parameters for the api call; encoded in url.
	Values *url.Values
}

// Exec executes a Request r and unmarshals the JSON response into resp.
// See godoc encoding/json Unmarshal for information on what to provide as resp.
// BasicAuth will override OAuth if those fields are set.
func Exec(client *http.Client, agent string, r *Request, resp interface{}) error {
	httpReq, err := r.httpRequest()
	if err != nil {
		return err
	}
	httpReq.Header.Add("user-agent", agent)

	httpResp, err := client.Do(httpReq)
	if err != nil {
		return err
	}

	if resp == nil {
		return nil
	}

	return parseResponse(httpResp, resp)
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

	return req, err
}

// parseResponse parses the JSON body of an http.Response into a type.
func parseResponse(resp *http.Response, val interface{}) error {
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(buf, val)
}

// postRequest returns a template http.Request with the given url and POST form
// values set.
func postRequest(url string, vals *url.Values) (*http.Request, error) {
	if vals == nil {
		return nil, errors.New("no values for POST body")
	}

	reqBody := bytes.NewBufferString(vals.Encode())
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", contentType)
	return req, nil
}

// getRequest returns a template http.Request with the given url and GET form
// values set.
func getRequest(url string, vals *url.Values) (*http.Request, error) {
	reqURL := url
	if vals != nil {
		reqURL = fmt.Sprintf("%s?%s", reqURL, vals.Encode())
	}
	return http.NewRequest("GET", reqURL, nil)
}
