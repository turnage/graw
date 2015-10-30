package monitor

import (
	"github.com/turnage/graw/internal/monitor/internal/handlers"
	"github.com/turnage/graw/internal/operator"
)

// userMonitor monitors a user for new posts and comments.
type userMonitor struct {
	base
}

// UserMonitor returns a user monitor for the requested user.
func UserMonitor(
	op operator.Operator,
	bot handlers.UserHandler,
	user string,
	dir Direction,
) (Monitor, error) {
	u := &userMonitor{
		base: base{
			handlePost:    bot.UserPost,
			handleComment: bot.UserComment,
			tip:           []string{""},
			dir:           dir,
			path:          "/user/" + user,
		},
	}

	if dir == Forward {
		if err := u.sync(op); err != nil {
			return nil, err
		}
	}

	return u, nil
}
