// parse_test uses actual Reddit  selected at human-random, to test that
// the parser can extract the information from them. The selected information I
// test isn't a full check because that would take a long time to write and I
// don't think it is that worthwhile. Instead some things are poked where I
// think there are likely to be issues.
//
// Find the result expecations in internal/test*.json
package reddit

import (
	"strings"
	"testing"

	"github.com/turnage/graw/reddit/internal/testdata"
)

func TestParse(t *testing.T) {
	p := newParser()
	for i, input := range [][]byte{
		testdata.MustAsset("thread.json"),
		testdata.MustAsset("user.json"),
		testdata.MustAsset("subreddit.json"),
		testdata.MustAsset("inbox.json"),
	} {
		if _, _, _, _, err := p.parse(input); err != nil {
			t.Errorf("failed to parse input %d: %v", i, err)
		}
	}
}

func TestParseThread(t *testing.T) {
	post, _, err := parseThread(testdata.MustAsset("thread.json"))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if post == nil {
		t.Fatalf("post is nil")
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

func TestParseUserFeed(t *testing.T) {
	comments, posts, _, _, err := parseRawListing(
		testdata.MustAsset("user.json"),
	)
	if err != nil {
		t.Fatalf("failed to parse user feed: %v", err)
	}

	if len(comments) < 1 {
		t.Fatalf("found no comments in user feed")
	}

	if len(posts) < 1 {
		t.Fatalf("found no posts in user feed")
	}

	if comments[0].LinkTitle != "Dreamworks LLC!!!" {
		t.Errorf(
			"user feed comment had unexpected link title: %s",
			comments[0].LinkTitle,
		)
	}

	if posts[0].Score != 417 {
		t.Errorf(
			"user feed post had unexpected score: %d",
			posts[0].Score,
		)
	}
}

func TestParseSubredditFeed(t *testing.T) {
	_, posts, _, _, err := parseRawListing(testdata.MustAsset("subreddit.json"))
	if err != nil {
		t.Fatalf("failed to parse subreddit feed: %v", err)
	}

	if len(posts) != 27 {
		t.Fatalf(
			"failed to parse all posts; found %d; wanted %d",
			len(posts), 27,
		)
	}

	if posts[0].Name != "t3_552rz1" {
		t.Errorf("failed to parse post name; found: %s", posts[0].Name)
	}

	if posts[26].LinkFlairCSSClass != "black" {
		t.Errorf(
			"failed to parse link flair css; found: %s",
			posts[26].LinkFlairCSSClass,
		)
	}
}

func TestParseInboxFeed(t *testing.T) {
	_, _, msgs, _, err := parseRawListing(testdata.MustAsset("inbox.json"))
	if err != nil {
		t.Fatalf("failed to parse inbox feed: %v", err)
	}

	if len(msgs) != 5 {
		t.Fatalf("found unexpected number of messages: %v", len(msgs))
	}

	if msgs[0].Name != "t1_cwup4dd" {
		t.Errorf("first message had unexpected name: %s", msgs[0].Name)
	}
}
