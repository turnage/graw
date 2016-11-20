package reddit

// Account defines behaviors only an account can perform on Reddit.
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
	r reaper
}

// newAccount returns a new Account using the given reaper to make requests
// to Reddit.
func newAccount(r reaper) Account {
	return &account{
		r: r,
	}
}

func (a *account) Reply(parentName, text string) error {
	return a.r.sow(
		"/api/comment", map[string]string{
			"thing_id": parentName,
			"text":     text,
		},
	)
}

func (a *account) SendMessage(user, subject, text string) error {
	return a.r.sow(
		"/api/compose", map[string]string{
			"to":      user,
			"subject": subject,
			"text":    text,
		},
	)
}

func (a *account) PostSelf(subreddit, title, text string) error {
	return a.r.sow(
		"/api/submit", map[string]string{
			"sr":    subreddit,
			"kind":  "self",
			"title": title,
			"text":  text,
		},
	)
}

func (a *account) PostLink(subreddit, title, url string) error {
	return a.r.sow(
		"/api/submit", map[string]string{
			"sr":    subreddit,
			"kind":  "link",
			"title": title,
			"url":   url,
		},
	)
}
