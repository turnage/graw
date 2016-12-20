// Package streams provides robust event streams from Reddit.
//
// This package is not abstract. If you are looking for a simpler, high level
// interface, see graw.
//
// The streams provided by this package will not be deterred like naive
// implementations by a post getting caught in the spam filter, removed by mods,
// the author being shadowbanned, or an author deleting their post. These
// streams are in it to win it.
//
// All of the streams provisioned by this package depend on an api handle from
// the reddit package, and two control channels: one kill signal and one error
// feed.
//
// The kill channel can be shared by multiple streams as long as you signal kill
// by close()ing the channel. Sending data over it will kill an arbitrary one of
// the streams sharing the channel but not all of them.
//
// The error channel will return issues which may be intermittent. They are not
// wrapped, so you can check them against the definitions in the reddit package
// and choose to wait when Reddit is busy or the connection faults, instead of
// failing.
//
// If there is a problem setting up the stream, such as the endpoint being
// invalid, that will be caught in the initial construction of the stream; you
// don't need to worry about that on the error channel.
//
// These streams will consume "intervals" of the Reddit handle given to them.
// Since the reddit handlers are rate limited and do not allow bursts, there is
// essentially a schedule on which they execute requests, and the executions
// will be divided roughly evenly between the goroutines sharing the handle.
// E.g. if you create two user streams which depend on a handle with a rate
// limit of 5 seconds, each of them will be unblocked once every 10 seconds
// (ish), since they each consume one interval, and the interval is 5 seconds.
package streams

import (
	"strings"

	"github.com/turnage/graw/reddit"

	"github.com/turnage/graw/streams/internal/monitor"
	"github.com/turnage/graw/streams/internal/rsort"
)

// Subreddits returns a stream of new posts from the requested subreddits. This
// stream monitors the combination listing of all subreddits using Reddit's "+"
// feature e.g. /r/golang+rust. This will consume one interval of the handle per
// call, so it is best to gather all the subreddits needed and invoke this
// function once.
//
// Be aware that these posts are new and will not have comments. If you are
// interested in comment trees, save their permalinks and fetch them later.
func Subreddits(
	scanner reddit.Scanner,
	kill <-chan bool,
	errs chan<- error,
	subreddits ...string,
) (
	<-chan *reddit.Post,
	error,
) {
	path := "/r/" + strings.Join(subreddits, "+") + "/new"
	posts, _, _, err := streamFromPath(scanner, kill, errs, path)
	return posts, err
}

// SubredditComments returns a stream of new comments from the requested
// subreddits. This stream monitors the combination listing of all subreddits
// using Reddit's "+" feature e.g. /r/golang+rust. This will consume one
// interval of the handle per call, so it is best to gather all the subreddits
// needed and invoke this function once.
//
// Be aware that these comments are new, and will not have reply trees. If you
// are interested in comment trees, save the permalinks of their parent posts
// and fetch them later once they may have had activity.
func SubredditComments(
	scanner reddit.Scanner,
	kill <-chan bool,
	errs chan<- error,
	subreddits ...string,
) (
	<-chan *reddit.Comment,
	error,
) {
	path := "/r/" + strings.Join(subreddits, "+") + "/comments"
	_, comments, _, err := streamFromPath(scanner, kill, errs, path)
	return comments, err
}

// User returns a stream of new posts and comments made by a user. Each user
// stream consumes one interval of the handle.
func User(
	scanner reddit.Scanner,
	kill <-chan bool,
	errs chan<- error,
	user string,
) (
	<-chan *reddit.Post,
	<-chan *reddit.Comment,
	error,
) {
	path := "/u/" + user
	posts, comments, _, err := streamFromPath(scanner, kill, errs, path)
	return posts, comments, err
}

// PostReplies returns a stream of top level replies to posts made by the bot's
// account. This stream consumes one interval of the handle.
func PostReplies(
	bot reddit.Bot,
	kill <-chan bool,
	errs chan<- error,
) (
	<-chan *reddit.Message,
	error,
) {
	return inboxStream(bot, kill, errs, "selfreply")
}

// CommentReplies returns a stream of replies to comments made by the bot's
// account. This stream consumes one interval of the handle.
func CommentReplies(
	bot reddit.Bot,
	kill <-chan bool,
	errs chan<- error,
) (
	<-chan *reddit.Message,
	error,
) {
	return inboxStream(bot, kill, errs, "comments")
}

// Mentions returns a stream of mentions of the bot's username anywhere on
// Reddit. It consumes one interval of the handle. Note, that a username mention
// which can reach the inbox in any other way (as a pm, or a reply), will not
// come through the mention stream because Reddit labels it differently.
func Mentions(
	bot reddit.Bot,
	kill <-chan bool,
	errs chan<- error,
) (
	<-chan *reddit.Message,
	error,
) {
	return inboxStream(bot, kill, errs, "mentions")
}

// Messages returns a stream of messages sent to the bot's inbox. It consumes
// one interval of the handle.
func Messages(
	bot reddit.Bot,
	kill <-chan bool,
	errs chan<- error,
) (
	<-chan *reddit.Message,
	error,
) {
	onlyMessages := make(chan *reddit.Message)

	messages, err := inboxStream(bot, kill, errs, "inbox")
	go func() {
		for m := range messages {
			if !m.WasComment {
				onlyMessages <- m
			}
		}
	}()

	return onlyMessages, err
}

func inboxStream(
	scanner reddit.Scanner,
	kill <-chan bool,
	errs chan<- error,
	subpath string,
) (
	<-chan *reddit.Message,
	error,
) {
	path := "/message/" + subpath
	_, _, messages, err := streamFromPath(scanner, kill, errs, path)
	return messages, err
}

func streamFromPath(
	scanner reddit.Scanner,
	kill <-chan bool,
	errs chan<- error,
	path string,
) (
	<-chan *reddit.Post,
	<-chan *reddit.Comment,
	<-chan *reddit.Message,
	error,
) {
	mon, err := monitorFromPath(path, scanner)
	if err != nil {
		return nil, nil, nil, err
	}

	posts, comments, messages := stream(mon, kill, errs)
	return posts, comments, messages, nil
}

func monitorFromPath(path string, sc reddit.Scanner) (monitor.Monitor, error) {
	return monitor.New(
		monitor.Config{
			Path:    path,
			Scanner: sc,
			Sorter:  rsort.New(),
		},
	)
}

func stream(
	mon monitor.Monitor,
	kill <-chan bool,
	errs chan<- error,
) (
	<-chan *reddit.Post,
	<-chan *reddit.Comment,
	<-chan *reddit.Message,
) {
	posts := make(chan *reddit.Post)
	comments := make(chan *reddit.Comment)
	messages := make(chan *reddit.Message)

	go flow(mon, kill, errs, posts, comments, messages)

	return posts, comments, messages
}

func flow(
	mon monitor.Monitor,
	kill <-chan bool,
	errs chan<- error,
	posts chan<- *reddit.Post,
	comments chan<- *reddit.Comment,
	messages chan<- *reddit.Message,
) {
	for {
		select {
		// if the errors channel is closed, the master goroutine is
		// shutting us down.
		case <-kill:
			close(posts)
			close(comments)
			close(messages)
			return
		default:
			if h, err := mon.Update(); err != nil {
				errs <- err
			} else {
				// lol no generics
				for _, p := range h.Posts {
					posts <- p
				}
				for _, c := range h.Comments {
					comments <- c
				}
				for _, m := range h.Messages {
					messages <- m
				}
			}
		}
	}
}
