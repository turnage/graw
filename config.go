package graw

import (
	"log"
)

type Config struct {
	Subreddits     []string
	Users          []string
	PostReplies    bool
	CommentReplies bool
	Mentions       bool
	Messages       bool
	Logger         *log.Logger
}
