package streams

import (
	"testing"

	"github.com/turnage/graw/internal/data"
)

type mockSubredditHandler struct{}

func (m *mockSubredditHandler) Post(_ *data.Post) error {
	return nil
}

type mockUserHandler struct{}

func (m *mockUserHandler) UserPost(_ *data.Post) error {
	return nil
}

func (m *mockUserHandler) UserComment(_ *data.Comment) error {
	return nil
}

type mockInboxHandler struct{}

func (m *mockInboxHandler) PostReply(_ *data.Message) error {
	return nil
}

func (m *mockInboxHandler) CommentReply(_ *data.Message) error {
	return nil
}

func (m *mockInboxHandler) Mention(_ *data.Message) error {
	return nil
}

func (m *mockInboxHandler) Message(_ *data.Message) error {
	return nil
}

func TestSubreddits(t *testing.T) {
	if _, err := subreddits([]string{"all"}, nil); err == nil {
		t.Errorf("wanted error for missing handler")
	}

	m := &mockSubredditHandler{}
	if cfg, err := subreddits([]string{"all"}, m); err != nil {
		t.Errorf("error creating subreddits config: %v", err)
	} else if cfg.path != "/r/all/new" || cfg.ph == nil {
		t.Errorf("incorrect config; got %v", cfg)
	}
}

func TestUsers(t *testing.T) {
	if _, err := users([]string{"user"}, nil); err == nil {
		t.Errorf("wanted error for missing handler")
	}

	m := &mockUserHandler{}
	if cfgs, err := users([]string{"user"}, m); err != nil {
		t.Errorf("error creating users config: %v", err)
	} else if len(cfgs) != 1 {
		t.Errorf("wanted 1 config; got %v", cfgs)
	} else if cfgs[0].path != "/u/user" ||
		cfgs[0].ph == nil ||
		cfgs[0].ch == nil {
		t.Errorf("incorrect config; got %v", cfgs)
	}
}

func TestInbox(t *testing.T) {
	if _, err := inbox(true, nil); err == nil {
		t.Errorf("wanted error for missing handler")
	}

	m := &mockInboxHandler{}
	if _, err := inbox(false, m); err == nil {
		t.Errorf("wanted error for being logged out")
	}

	if cfg, err := inbox(true, m); err != nil {
		t.Errorf("error creating inbox config: %v", err)
	} else if cfg.path != "/message/inbox" || cfg.mh == nil {
		t.Errorf("incorrect config; got %v", cfg)
	}
}
