package streams

import (
	"strings"

	"github.com/turnage/graw/reddit"

	"github.com/turnage/graw/streams/internal/monitor"
	"github.com/turnage/graw/streams/internal/rsort"
)

func Subreddits(
	scanner reddit.Scanner,
	kill <-chan bool,
	errs chan<- error,
	subreddits ...string,
) (
	<-chan *reddit.Post,
	error,
) {
	path := "/r/" + strings.Join(subreddits, "+")
	posts, _, _, err := streamFromPath(scanner, kill, errs, path)
	return posts, err
}

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

func Messages(
	bot reddit.Bot,
	kill <-chan bool,
	errs chan<- error,
) (
	<-chan *reddit.Message,
	error,
) {
	return inboxStream(bot, kill, errs, "messages")
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
