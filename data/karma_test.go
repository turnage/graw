package data

import (
	"encoding/json"
	"testing"
)

func TestKarmaList(t *testing.T) {
	karma := &KarmaList{}
	err := json.Unmarshal([]byte(`{
		"data": [
			{
				"sr": "self",
				"comment_karma": 80,
				"link_karma": 60
			}
		]
	}`), karma)
	if err != nil {
		t.Fatal("failed to unmarshal KarmaList")
	}

	subreddits := karma.GetData()

	if len(subreddits) != 1 {
		t.Errorf("should have 1 subreddit; found %d", len(subreddits))
	}

	if subreddits[0].GetSr() != "self" {
		t.Errorf("should have self; found %s", subreddits[0].GetSr())
	}

	if subreddits[0].GetCommentKarma() != 80 {
		t.Errorf(
			"should have 80 comment karma; found %d",
			subreddits[0].GetCommentKarma())
	}

	if subreddits[0].GetLinkKarma() != 60 {
		t.Errorf(
			"should have 60 link karma; found %d",
			subreddits[0].GetLinkKarma())
	}
}
