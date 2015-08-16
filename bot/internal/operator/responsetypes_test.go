package operator

import (
	"encoding/json"
	"testing"
)

func TestLinks(t *testing.T) {
	resp := &linkListing{}
	if err := json.Unmarshal([]byte(`{
		"data": {
			"children": [
				{"data": {"title": "hello"}},
				{"data": {"title": "hola"}},
				{"data": {"title": "bye"}}
			]
		}
	}`), resp); err != nil {
		t.Fatalf("failed to prepare test input struct: %v", err)
	}
	if len(resp.Links()) != 3 {
		t.Errorf("wanted to find 3 links; resp is %v", resp)
	}
}
