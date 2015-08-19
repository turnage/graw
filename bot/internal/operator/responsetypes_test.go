package operator

import (
	"encoding/json"
	"testing"

	"github.com/turnage/redditproto"
)

func TestLinks(t *testing.T) {
	resp := &redditproto.LinkListing{}
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
	if len(getLinks(resp)) != 3 {
		t.Errorf("wanted to find 3 links; resp is %v", resp)
	}
}
