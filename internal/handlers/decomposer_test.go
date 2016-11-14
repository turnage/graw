package handlers

import (
	"testing"

	"github.com/turnage/graw/internal/data"
)

type mockSubredditHandler struct {
	called bool
}

func (m *mockSubredditHandler) Post(_ *data.Post) error {
	m.called = true
	return nil
}

type mockUserHandler struct {
	postCalled    bool
	commentCalled bool
}

func (m *mockUserHandler) UserPost(_ *data.Post) error {
	m.postCalled = true
	return nil
}

func (m *mockUserHandler) UserComment(_ *data.Comment) error {
	m.commentCalled = true
	return nil
}

type mockInboxHandler struct {
	postReplyCalled    bool
	commentReplyCalled bool
	mentionCalled      bool
	messageCalled      bool
}

func (m *mockInboxHandler) PostReply(_ *data.Message) error {
	m.postReplyCalled = true
	return nil
}

func (m *mockInboxHandler) CommentReply(_ *data.Message) error {
	m.commentReplyCalled = true
	return nil
}

func (m *mockInboxHandler) Mention(_ *data.Message) error {
	m.mentionCalled = true
	return nil
}

func (m *mockInboxHandler) Message(_ *data.Message) error {
	m.messageCalled = true
	return nil
}

func TestDecomposeSubredditHandler(t *testing.T) {
	m := &mockSubredditHandler{}
	if DecomposeSubredditHandler(m).HandlePost(nil); !m.called {
		t.Errorf("Subreddit handler's function was not called.")
	}
}

func TestDecomposeUserHandler(t *testing.T) {
	m := &mockUserHandler{}
	ph, ch := DecomposeUserHandler(m)
	if ph.HandlePost(nil); !m.postCalled {
		t.Errorf("User handler's Post function was not called.")
	}
	if ch.HandleComment(nil); !m.commentCalled {
		t.Errorf("User handler's Post function was not called.")
	}
}

func TestDecomposeInboxHandler(t *testing.T) {
	m := &mockInboxHandler{}
	mh := DecomposeInboxHandler(m)
	for _, test := range []struct {
		input   *data.Message
		flipped *bool
	}{
		{
			&data.Message{
				WasComment: true,
				Subject:    "comment reply",
			},
			&m.commentReplyCalled,
		},
		{
			&data.Message{
				WasComment: true,
				Subject:    "post reply",
			},
			&m.postReplyCalled,
		},
		{
			&data.Message{
				WasComment: true,
				Subject:    "username mention",
			},
			&m.mentionCalled,
		},
		{
			&data.Message{
				Subject: "comment reply",
			},
			&m.messageCalled,
		},
	} {
		if mh.HandleMessage(test.input); !*test.flipped {
			t.Errorf("Did not get proper route for %v", test.input)
		}
	}
}
