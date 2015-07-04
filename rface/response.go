package rface

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// parseResponse parses the JSON body of an http.Response into a type.
// See the godoc on json.Unmarshal for information on what to provide as a
// a val and how to set it up for parsing.
func parseResponse(resp *http.Response, val interface{}) error {
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(buf, val)
}
