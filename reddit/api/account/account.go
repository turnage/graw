// Package account implements a representation of user accounts on Reddit, for
// interacting with Reddit.
package account

import (
	"time"

	"github.com/turnage/graw/internal/api"
	"github.com/turnage/graw/internal/reap"
)

type Account interface {
	// Reply posts a reply to something on reddit. The behavior depends on
	// what is being replied to. For
	//
	//   messages, this sends a private message reply.
	//   posts, this posts a top level comment.
	//   comments, this posts a comment reply.
	//
	// Use .Name on the parent post, message, or comment to find its
	// name.
	Reply(parentName, text string) error

	// SendMessage sends a private message to a user.
	SendMessage(user, subject, text string) error

	// PostSelf makes a text (self) post to a subreddit.
	PostSelf(subreddit, title, text string) error

	// PostLink makes a link post to a subreddit.
	PostLink(subreddit, title, url string) error
}

type account struct {
	// r is used to execute requests to Reddit.
	r reap.Reaper
	// grant is a channel over which Account receives grants for issuing
	// requests.
	grant <-chan time.Time
}

// New returns a new Account using the given reaper to make requests Reddit.
func New(r reap.Reaper, grant <-chan time.Time) Account {
	return &account{
		r:     r,
		grant: grant,
	}
}

func (a *account) Reply(parentName, text string) error {
	<-a.grant
	return a.r.Sow(
		"/api/comment", api.WithDefaults(
			map[string]string{
				"thing_id": parentName,
				"text":     text,
			},
		),
	)
}

func (a *account) SendMessage(user, subject, text string) error {
	<-a.grant
	return a.r.Sow(
		"/api/compose", api.WithDefaults(
			map[string]string{
				"to":      user,
				"subject": subject,
				"text":    text,
			},
		),
	)
}

func (a *account) PostSelf(subreddit, title, text string) error {
	return a.post(subreddit, title, text, "self")
}

func (a *account) PostLink(subreddit, title, url string) error {
	return a.post(subreddit, title, url, "link")
}

func (a *account) post(subreddit, title, content, kind string) error {
	<-a.grant
	return a.r.Sow(
		"/api/submit", api.WithDefaults(
			map[string]string{
				"sr":    subreddit,
				"kind":  kind,
				"title": title,
				"url":   content,
				"text":  content,
			},
		),
	)
}
