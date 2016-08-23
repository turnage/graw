package engine

import (
	"container/list"

	"github.com/turnage/graw/botfaces"
	"github.com/turnage/graw/internal/api"
	"github.com/turnage/graw/internal/client"
	"github.com/turnage/graw/internal/monitor"
	"github.com/turnage/redditproto"
)

func New(
	bot interface{},
	cli client.Client,
	subreddits []string,
) (*Engine, error) {
	e := &Engine{
		cli:          cli,
		bot:          bot,
		monitors:     list.New(),
		userMonitors: make(map[string]*list.Element),
		stopSig:      make(chan bool),
	}

	scraper := func(path, tip string, limit int) (
		[]*redditproto.Link,
		[]*redditproto.Comment,
		[]*redditproto.Message,
		error,
	) {
		return api.Scrape(e.cli.Do, path, tip, limit)
	}

	if han, ok := bot.(botfaces.PostHandler); ok && len(subreddits) > 0 {
		mon, err := monitor.PostMonitor(scraper, han.Post, subreddits)
		if err != nil {
			return nil, err
		}
		e.monitors.PushFront(mon)
	}

	if han, ok := bot.(botfaces.MessageHandler); ok {
		mon, err := monitor.MessageMonitor(scraper, han.Message)
		if err != nil {
			return nil, err
		}
		e.monitors.PushFront(mon)
	}

	if han, ok := bot.(botfaces.PostReplyHandler); ok {
		mon, err := monitor.PostReplyMonitor(scraper, han.PostReply)
		if err != nil {
			return nil, err
		}
		e.monitors.PushFront(mon)
	}

	if han, ok := bot.(botfaces.CommentReplyHandler); ok {
		mon, err := monitor.CommentReplyMonitor(scraper, han.CommentReply)
		if err != nil {
			return nil, err
		}
		e.monitors.PushFront(mon)
	}

	if han, ok := bot.(botfaces.MentionHandler); ok {
		mon, err := monitor.MentionMonitor(scraper, han.Mention)
		if err != nil {
			return nil, err
		}
		e.monitors.PushFront(mon)
	}

	return e, nil
}
