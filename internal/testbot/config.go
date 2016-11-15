package main

import (
	"log"
	"os"
	"time"

	"github.com/turnage/graw"
)

var config = graw.Config{
	Agent:      "graw:graw-internal-testbot:0.1.0 (by /u/roxven)",
	Subreddits: []string{},
	Users:      []string{},
	Inbox:      true,
	Rate:       5 * time.Second,
	App: &graw.App{
		ID:       "",
		Secret:   "",
		Username: "",
		Password: "",
	},
	Logger: log.New(os.Stderr, "", log.LstdFlags),
}
