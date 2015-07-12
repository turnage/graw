package api

import (
	"fmt"
	"testing"

	"github.com/paytonturnage/graw/nface"
)

func TestMeRequest(t * testing.T) {
	req := MeRequest()
	if req.Action != nface.GET {
		t.Errorf("action incorrect; expected %v, got %v", nface.GET, req.Action)
	}

	expectedURL := fmt.Sprintf("%s%s", baseURL, meURL)
	if req.BaseURL != expectedURL {
		t.Errorf("url incorrect; expected %s, got %s", expectedURL, req.BaseURL)
	}
}
