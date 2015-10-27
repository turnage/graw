package monitor

import (
	"github.com/turnage/graw/api"
	"github.com/turnage/graw/internal/monitor/internal/scanner"
	"github.com/turnage/graw/internal/operator"
)

// userMonitor monitors a user for new posts and comments.
type userMonitor struct {
	// userScanner is the scanner userMonitor uses to get updates from
	// the user it monitors.
	userScanner scanner.Scanner
	// userHandler is the handler UserMonitor will send new posts and
	// comments by the watched user to.
	userHandler api.UserHandler
}

// UserMonitor returns a user monitor for the requested user.
func UserMonitor(
	op operator.Operator,
	bot api.UserHandler,
	user string,
) Monitor {
	return &userMonitor{
		userScanner: scanner.NewUserScanner(user, op),
		userHandler: bot,
	}
}

// Update polls for new user activity and sends the content to Bot when it is
// found.
func (p *userMonitor) Update() error {
	posts, comments, err := p.userScanner.Scan()
	if err != nil {
		return err
	}

	for _, post := range posts {
		go p.userHandler.UserPost(post)
	}

	for _, comment := range comments {
		go p.userHandler.UserComment(comment)
	}

	return nil
}
