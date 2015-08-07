package graw

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func exec(c client, r *http.Request, out interface{}) error {
	rawResp, err := c.do(r)
	if err != nil {
		return err
	}

	if rawResp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code in response")
	}

	if rawResp.Body == nil {
		return fmt.Errorf("no body in response")
	}
	defer rawResp.Body.Close()

	buffer, err := ioutil.ReadAll(rawResp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(buffer, out)
}
