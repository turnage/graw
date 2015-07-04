// Pacakge nface handles all communication between Go code and the Reddit api.
package nface

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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
	// contentType is a header flag for POST requests so the reddit api
	// knows how to read the request body.
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

// Exec executes a Request r and unmarshals the JSON response into resp.
// See godoc encoding/json Unmarshal for information on what to provide as resp.
// BasicAuth will override OAuth if those fields are set.
func Exec(r *Request, resp interface{}) error {
	httpReq, err := r.httpRequest()
	if err != nil {
		return err
	}

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
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
	reqUrl := fmt.Sprintf("%s?%s", url, vals.Encode())
	return http.NewRequest("GET", reqUrl, nil)
}
