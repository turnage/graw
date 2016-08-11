package engine

import (
	"container/list"

	"github.com/turnage/graw/internal/botfaces"
	"github.com/turnage/graw/internal/monitor"
	"github.com/turnage/graw/internal/operator"
)

func RealTime(
	bot interface{},
	op operator.Operator,
	subreddits []string,
) (*Engine, error) {
	return baseFrom(bot, op, subreddits, monitor.Forward)
}

func BackTime(
	bot interface{},
	op operator.Operator,
	subreddits []string,
) (*Engine, error) {
	return baseFrom(bot, op, subreddits, monitor.Backward)
}

func baseFrom(
	bot interface{},
	op operator.Operator,
	subreddits []string,
	dir monitor.Direction,
) (*Engine, error) {
	e := &Engine{
		op:           op,
		bot:          bot,
		dir:          dir,
		monitors:     list.New(),
		userMonitors: make(map[string]*list.Element),
		stopSig:      make(chan bool),
	}

	if han, ok := bot.(botfaces.PostHandler); ok && len(subreddits) > 0 {
		mon, err := monitor.PostMonitor(op, han.Post, subreddits, dir)
		if err != nil {
			return nil, err
		}
		e.monitors.PushFront(mon)
	}

	if han, ok := bot.(botfaces.MessageHandler); ok {
		mon, err := monitor.MessageMonitor(op, han.Message, dir)
		if err != nil {
			return nil, err
		}
		e.monitors.PushFront(mon)
	}

	if han, ok := bot.(botfaces.PostReplyHandler); ok {
		mon, err := monitor.PostReplyMonitor(op, han.PostReply, dir)
		if err != nil {
			return nil, err
		}
		e.monitors.PushFront(mon)
	}

	if han, ok := bot.(botfaces.CommentReplyHandler); ok {
		mon, err := monitor.CommentReplyMonitor(
			op,
			han.CommentReply,
			dir,
		)
		if err != nil {
			return nil, err
		}
		e.monitors.PushFront(mon)
	}

	if han, ok := bot.(botfaces.MentionHandler); ok {
		mon, err := monitor.MentionMonitor(op, han.Mention, dir)
		if err != nil {
			return nil, err
		}
		e.monitors.PushFront(mon)
	}

	return e, nil
}
