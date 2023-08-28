// parse_test uses actual Reddit  selected at human-random, to test that
// the parser can extract the information from them. The selected information I
// test isn't a full check because that would take a long time to write and I
// don't think it is that worthwhile. Instead some things are poked where I
// think there are likely to be issues.
// Find the result expectations in internal/test*.json
package reddit

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/mix/graw/reddit/internal/testdata"
)

func TestParse(t *testing.T) {
	p := newParser()
	for i, input := range [][]byte{
		testdata.MustAsset("thread.json"),
		testdata.MustAsset("user.json"),
		testdata.MustAsset("subreddit.json"),
		testdata.MustAsset("inbox.json"),
		testdata.MustAsset("more.json"),
	} {
		if _, _, _, _, err := p.parse(input); err != nil {
			t.Errorf("failed to parse input %d: %v", i, err)
		}
	}
	if _, _, _, _, err := p.parse(json.RawMessage(ThreadWithPreview)); err != nil {
		t.Errorf("failed to parse input ThreadWithPreview: %v", err)
	}
	if _, _, _, _, err := p.parse(json.RawMessage(ThreadWithRedditVideoPreview)); err != nil {
		t.Errorf("failed to parse input ThreadWithRedditVideoPreview: %v", err)
	}
	if _, _, _, _, err := p.parse(json.RawMessage(ThreadRemovedByCategoryModerator)); err != nil {
		t.Errorf("failed to parse input ThreadRemovedByCategoryModerator: %v", err)
	}
}

func TestParseThread(t *testing.T) {
	post, err := parseThread(testdata.MustAsset("thread.json"))
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
	if post.RemovedByCategory != "" {
		t.Errorf("post removedByCategory incorrect: %s", post.RemovedByCategory)
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

	if post.Replies[0].Edited != 0 {
		t.Errorf(
			"first comment's edit state is incorrect: %d",
			post.Replies[0].Edited,
		)
	}

	if post.Replies[1].Replies[0].Replies[1].Edited != 1366216653 {
		t.Errorf(
			"comment %s edit timestamp is incorrect: %d",
			post.Replies[1].Replies[0].Replies[1].ID,
			post.Replies[1].Replies[0].Replies[1].Edited,
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

	if post.Replies[0].Replies[0].Replies[0].Replies[0].More.Name != "t1_c9hrw24" {
		t.Errorf(
			"sub reply had incorrect more: %s",
			post.Replies[0].Replies[0].Replies[0].Replies[0].More.Name,
		)
	}

	if post.More.Name != "t1_c9gzz0k" {
		t.Errorf(
			"post had incorrect more: %s",
			post.More.Name,
		)
	}

	if len(post.More.Children) != 649 {
		t.Errorf(
			"post more had incorrect number of children: %d",
			len(post.More.Children),
		)
	}
}

func TestParseThreadWithPreview(t *testing.T) {
	post, err := parseThread(json.RawMessage(ThreadWithPreview))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}
	if post == nil {
		t.Fatalf("post is nil")
	}

	if !strings.HasPrefix(post.Title, "🔥😆") {
		t.Errorf("post title incorrect: %s", post.Title)
	}

	if post.Author != "low-vibe" {
		t.Errorf("post author incorrect: %s", post.Author)
	}
	if post.RemovedByCategory != "" {
		t.Errorf("post removedByCategory incorrect: %s", post.RemovedByCategory)
	}

	if post.Preview.Images == nil {
		t.Errorf("post preview.images is nil")
	}
	if len(post.Preview.Images) != 1 {
		t.Errorf("len(post.Preview.Images) != 1 : equals %d", len(post.Preview.Images))
	}
	if post.Preview.Images[0].Source.URL != "https://external-preview.redd.it/ZaH9VJhfj6gd_dwIyD3sWJRqlpu3c8zTEhhM1nS9ifo.png?auto=webp&amp;v=enabled&amp;s=8823e5882da7fe88a33ac849802fb2f97787346a" {
		t.Errorf("post.Preview.Images[0].Source.URL"+
			" https://external-preview.redd.it/ZaH9VJhfj6gd_dwIyD3sWJRqlpu3c8zTEhhM1nS9ifo.png?auto=webp&amp;v=enabled&amp;s=8823e5882da7fe88a33ac849802fb2f97787346a != %s",
			post.Preview.Images[0].Source.URL)
	}
	if post.Preview.Images[0].Source.Width != 2160 {
		t.Errorf("post.Preview.Images[0].Source.Width"+
			" 2160 !=  %d",
			post.Preview.Images[0].Source.Width)
	}
	if post.Preview.Images[0].Source.Height != 3840 {
		t.Errorf("post.Preview.Images[0].Source.Height"+
			" 3840 !=  %d",
			post.Preview.Images[0].Source.Height)
	}
	if post.Preview.RedditVideoPreview != nil {
		t.Errorf("post preview.RedditVideoPreview is not nil")
	}
}

func TestParseThreadWithRedditVideoPreview(t *testing.T) {
	post, err := parseThread(json.RawMessage(ThreadWithRedditVideoPreview))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}
	if post == nil {
		t.Fatalf("post is nil")
	}

	if !strings.HasPrefix(post.Title, "Inside artificial horizon (or attitude indicator) that informs pilot of the aircraft orientation relative to Earth's horizon") {
		t.Errorf("post title incorrect: %s", post.Title)
	}

	if post.Author != "toolgifs" {
		t.Errorf("post author incorrect: %s", post.Author)
	}
	if post.RemovedByCategory != "" {
		t.Errorf("post removedByCategory incorrect: %s", post.RemovedByCategory)
	}

	if post.Preview.Images == nil {
		t.Errorf("post preview.images is nil")
	}
	if len(post.Preview.Images) != 1 {
		t.Errorf("len(post.Preview.Images) != 1 : equals %d", len(post.Preview.Images))
	}
	if post.Preview.Images[0].Source.URL != "https://external-preview.redd.it/WqhRau_Lo2LJrtlZfv-iMXTwy4pq1Zkc4gMdREumfzI.jpg?auto=webp&amp;v=enabled&amp;s=2c6ae912de30addfbaf0a0b5566607f50fcbe2e5" {
		t.Errorf("post.Preview.Images[0].Source.URL"+
			" https://external-preview.redd.it/WqhRau_Lo2LJrtlZfv-iMXTwy4pq1Zkc4gMdREumfzI.jpg?auto=webp&amp;v=enabled&amp;s=2c6ae912de30addfbaf0a0b5566607f50fcbe2e5 != %s",
			post.Preview.Images[0].Source.URL)
	}
	if post.Preview.Images[0].Source.Width != 960 {
		t.Errorf("post.Preview.Images[0].Source.Width"+
			" 960 !=  %d",
			post.Preview.Images[0].Source.Width)
	}
	if post.Preview.Images[0].Source.Height != 540 {
		t.Errorf("post.Preview.Images[0].Source.Height"+
			" 540 !=  %d",
			post.Preview.Images[0].Source.Height)
	}
	if post.Preview.RedditVideoPreview == nil {
		t.Errorf("post preview.RedditVideoPreview is nil")
	}
	if post.Preview.RedditVideoPreview.BitrateKPBS != 1200 {
		t.Errorf("post.Preview.RedditVideoPreview.BitrateKPBS"+
			" 1200 !=  %d",
			post.Preview.RedditVideoPreview.BitrateKPBS)
	}
	if post.Preview.RedditVideoPreview.FallbackURL != "https://v.redd.it/d9voiiq7e3e91/DASH_480.mp4" {
		t.Errorf("post.Preview.RedditVideoPreview.FallbackURL"+
			" https://v.redd.it/d9voiiq7e3e91/DASH_480.mp4 !=  %s",
			post.Preview.RedditVideoPreview.FallbackURL)
	}
	if post.Preview.RedditVideoPreview.Width != 853 {
		t.Errorf("post.Preview.RedditVideoPreview.Width"+
			" 853 !=  %d",
			post.Preview.RedditVideoPreview.Width)
	}
	if post.Preview.RedditVideoPreview.Height != 480 {
		t.Errorf("post.Preview.RedditVideoPreview.Height"+
			" 480 !=  %d",
			post.Preview.RedditVideoPreview.Height)
	}
	if post.Preview.RedditVideoPreview.ScrubberMediaURL != "https://v.redd.it/d9voiiq7e3e91/DASH_96.mp4" {
		t.Errorf("post.Preview.RedditVideoPreview.ScrubberMediaURL"+
			"https://v.redd.it/d9voiiq7e3e91/DASH_96.mp4 !=  %s",
			post.Preview.RedditVideoPreview.ScrubberMediaURL)
	}
	if post.Preview.RedditVideoPreview.DashURL != "https://v.redd.it/d9voiiq7e3e91/DASHPlaylist.mpd" {
		t.Errorf("post.Preview.RedditVideoPreview.DashURL"+
			"https://v.redd.it/d9voiiq7e3e91/DASHPlaylist.mpd !=  %s",
			post.Preview.RedditVideoPreview.DashURL)
	}
	if post.Preview.RedditVideoPreview.Duration != 31 {
		t.Errorf("post.Preview.RedditVideoPreview.Duration"+
			" 31 !=  %d",
			post.Preview.RedditVideoPreview.Duration)
	}
	if post.Preview.RedditVideoPreview.HLSURL != "https://v.redd.it/d9voiiq7e3e91/HLSPlaylist.m3u8" {
		t.Errorf("post.Preview.RedditVideoPreview.HLSURL"+
			"https://v.redd.it/d9voiiq7e3e91/HLSPlaylist.m3u8 !=  %s",
			post.Preview.RedditVideoPreview.HLSURL)
	}
	if post.Preview.RedditVideoPreview.IsGIF != true {
		t.Errorf("post.Preview.RedditVideoPreview.IsGIF"+
			" true !=  %t",
			post.Preview.RedditVideoPreview.IsGIF)
	}
	if post.Preview.RedditVideoPreview.TranscodingStatus != "completed" {
		t.Errorf("post.Preview.RedditVideoPreview.HLSURL"+
			"completed !=  %s",
			post.Preview.RedditVideoPreview.TranscodingStatus)
	}
}

func TestParseThreadRemovedByModerator(t *testing.T) {
	post, err := parseThread(json.RawMessage(ThreadRemovedByCategoryModerator))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}
	if post == nil {
		t.Fatalf("post is nil")
	}

	if !strings.HasPrefix(post.Title, "Bronze casting in a sand mold") {
		t.Errorf("post title incorrect: %s", post.Title)
	}

	if post.Author != "neonroli47" {
		t.Errorf("post author incorrect: %s", post.Author)
	}
	if post.RemovedByCategory != "moderator" {
		t.Errorf("post removedByCategory incorrect: %s", post.RemovedByCategory)
	}
}
func TestParserParseRemovedByModerator(t *testing.T) {
	p := newParser()
	_, posts, _, _, err := p.parse([]byte(ThreadRemovedByCategoryModerator))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}
	if posts[0] == nil {
		t.Fatalf("posts[0] is nil")
	}

	if !strings.HasPrefix(posts[0].Title, "Bronze casting in a sand mold") {
		t.Errorf("posts[0] title incorrect: %s", posts[0].Title)
	}

	if posts[0].Author != "neonroli47" {
		t.Errorf("posts[0] author incorrect: %s", posts[0].Author)
	}
	if posts[0].RemovedByCategory != "moderator" {
		t.Errorf("posts[0] removedByCategory incorrect: %s", posts[0].RemovedByCategory)
	}
}

func TestParserParseThreadGallery(t *testing.T) {
	p := newParser()
	_, posts, _, _, err := p.parse([]byte(ThreadImageGallery))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}
	if posts[0] == nil {
		t.Fatalf("posts[0] is nil")
	}

	if !strings.HasPrefix(posts[0].Title, "An outfit I made for my Mata Hari cosplay: top and pants, belt, headdress and necklaces") {
		t.Errorf("posts[0] title incorrect: %s", posts[0].Title)
	}

	if posts[0].Author != "fel0ra" {
		t.Errorf("posts[0] author incorrect: %s", posts[0].Author)
	}
	if len(posts[0].MediaMetadata) != 5 {
		t.Errorf("posts[0] len(posts[0].MediaMetadata) incorrect: %d", len(posts[0].MediaMetadata))
	}
	if len(posts[0].GalleryData.Items) != 5 {
		t.Errorf("posts[0] len(posts[0].GalleryData) incorrect: %d", len(posts[0].GalleryData.Items))
	}
	if posts[0].GalleryData.Items[0].ID != 278264474 {
		t.Errorf("posts[0] GalleryData.Items[0].ID incorrect: %d", posts[0].GalleryData.Items[0].ID)
	}
	if posts[0].GalleryData.Items[0].MediaId != "cyxb8fawmd1b1" {
		t.Errorf("posts[0] GalleryData.Items[0].MediaId incorrect: %s", posts[0].GalleryData.Items[0].MediaId)
	}
	if posts[0].MediaMetadata["cyxb8fawmd1b1"].ID != "cyxb8fawmd1b1" {
		t.Errorf("posts[0] MediaMetadata[\"cyxb8fawmd1b1\"].ID incorrect: %s", posts[0].MediaMetadata["cyxb8fawmd1b1"].ID)
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

func TestParseMoreChildren(t *testing.T) {
	comments, mores, err := parseMoreChildren(testdata.MustAsset("more.json"))
	if err != nil {
		t.Fatalf("failed to parse more children: %v", err)
	}

	if len(comments) != 50 {
		t.Fatalf("found unexpected number of comments: %v", len(comments))
	}

	if len(mores) != 11 {
		t.Fatalf("found unexpected number of mores: %v", len(mores))
	}
}
