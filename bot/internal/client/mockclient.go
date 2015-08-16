package client

import (
	"encoding/json"
	"net/http"
)

type mockClient struct {
	response []byte
}

// Do returns the preconfigured response in the mock client.
func (m *mockClient) Do(r *http.Request, out interface{}) error {
	return json.Unmarshal(m.response, out)
}
