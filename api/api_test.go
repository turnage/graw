package api

import (
	"testing"

	"github.com/paytonturnage/graw/nface"
)

func TestMeRequest(t *testing.T) {
	req := MeRequest()
	if req.Action != nface.GET {
		t.Errorf("action incorrect; expected %v, got %v", nface.GET, req.Action)
	}

	if req.URL != meURL {
		t.Errorf("url incorrect; expected %s, got %s", meURL, req.URL)
	}
}

func TestMeKarmaRequest(t *testing.T) {
	req := MeKarmaRequest()
	if req.Action != nface.GET {
		t.Errorf("action incorrect; expected %v, got %v", nface.GET, req.Action)
	}

	if req.URL != meKarmaURL {
		t.Errorf("url incorrect; expected %s, got %s", meKarmaURL, req.URL)
	}
}
