package operator

import (
	"net/http"
	"testing"

	"github.com/turnage/graw/bot/internal/client"
)

func TestGetListing(t *testing.T) {
	op := &Operator{
		cli: client.NewMock(
			`{"data":{"children":[{"data":{"title":"hey"}}]}}`,
		),
	}
	posts, err := op.getLinkListing(&http.Request{})
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if len(posts) != 1 {
		t.Errorf("wanted one post with title 'hey'; got %v", posts)
	}
}
