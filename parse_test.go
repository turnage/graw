package graw

import (
	"strings"
	"testing"

	"github.com/turnage/graw/internal/testdata"
)

func TestParseThread(t *testing.T) {
	post, err := parseThread(testdata.MustAsset("thread.json"))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if post == nil {
		t.Fatalf("post data is nil")
	}

	if !strings.HasPrefix(post.Title, "my wife passed away") {
		t.Errorf("post title incorrect: %s", post.Title)
	}

	if !strings.HasPrefix(post.SelfText, "it feels like I'm losing") {
		t.Errorf("post body incorrect: %s", post.SelfText)
	}

	if post.Author != "hglkgkjd" {
		t.Errorf("post author incorrect: %s", post.Author)
	}

	if len(post.Replies) == 0 {
		t.Fatal("post has no replies but it should")
	}

	if post.Replies[0].Author != "bacon_cake" {
		t.Errorf(
			"first comment has incorrect author: %s",
			post.Replies[0].Author,
		)
	}

	if len(post.Replies[0].Replies) < 1 {
		t.Fatalf("bacon_cake should have replies but doesn't")
	}

	if post.Replies[0].Replies[0].Author != "hglkgkjd" {
		t.Errorf(
			"sub reply had incorrect author: %s",
			post.Replies[0].Replies[0].Author,
		)
	}
}
