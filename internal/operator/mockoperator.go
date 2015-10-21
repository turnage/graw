package operator

import (
	"github.com/turnage/redditproto"
)

// MockOperator mocks Operator; it returns canned responses.
type MockOperator struct {
	// PostsErr is returned in the error field of Posts.
	PostsErr error
	// PostsReturn is returned by Posts.
	PostsReturn []*redditproto.Link
	// UserContentErr is returned in the error field of UserContent.
	UserContentErr error
	// UserContentLinksReturn is returned by UserContent.
	UserContentLinksReturn []*redditproto.Link
	// UserContentCommentsReturn is returned by UserContent.
	UserContentCommentsReturn []*redditproto.Comment
	// IsThereThingErr is returned in the error field of IsThereThing.
	IsThereThingErr error
	// IsThereThingReturn is returned by IsThereThing.
	IsThereThingReturn bool
	// ThreadErr is returned in the error field of Thread.
	ThreadErr error
	// ThreadReturn is returned by Thread.
	ThreadReturn *redditproto.Link
	// InboxErr is returned in the error field of Inbox.
	InboxErr error
	// InboxReturn is returned by Inbox.
	InboxReturn []*redditproto.Message
	// MarkAsReadErr is returned in the error field of MarkAsRead.
	MarkAsReadErr error
	// ReplyErr is returned in the error field of Reply.
	ReplyErr error
	// SubmitErr is returned in the error field of Submit.
	SubmitErr error
	// ComposeErr is returned in the error field of Compose.
	ComposeErr error
}

func (m *MockOperator) Posts(
	subreddit,
	after,
	before string,
	limit uint,
) ([]*redditproto.Link, error) {
	return m.PostsReturn, m.PostsErr
}

func (m *MockOperator) UserContent(
	user,
	after,
	before string,
	limit uint,
) ([]*redditproto.Link, []*redditproto.Comment, error) {
	return m.UserContentLinksReturn,
		m.UserContentCommentsReturn,
		m.UserContentErr
}

func (m *MockOperator) IsThereThing(id string) (bool, error) {
	return m.IsThereThingReturn, m.IsThereThingErr
}

func (m *MockOperator) Thread(permalink string) (*redditproto.Link, error) {
	return m.ThreadReturn, m.ThreadErr
}

func (m *MockOperator) Inbox() ([]*redditproto.Message, error) {
	return m.InboxReturn, m.InboxErr
}

func (m *MockOperator) MarkAsRead(fullnames ...string) error {
	return m.MarkAsReadErr
}

func (m *MockOperator) Reply(parent, content string) error {
	return m.ReplyErr
}

func (m *MockOperator) Submit(subreddit, kind, title, content string) error {
	return m.SubmitErr
}

func (m *MockOperator) Compose(user, subject, content string) error {
	return m.ComposeErr
}
