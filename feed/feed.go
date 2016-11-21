package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

var (
	subreddit = flag.String("subreddit", "all", "subreddit to watch")
	rate      = flag.Duration("rate", time.Second, "interval between updates")
)

type announcer struct{}

func (a *announcer) Post(p *reddit.Post) error {
	fmt.Printf("[%s]: %s\n", p.Author, p.Title)
	return nil
}

func main() {
	flag.Parse()

	cfg := graw.Config{
		Subreddits: []string{*subreddit},
		Logger:     log.New(os.Stderr, "", log.LstdFlags),
	}
	if script, err := reddit.NewScript(
		"graw:feed demo bot:0.5.1 by /u/roxven",
		*rate,
	); err != nil {
		fmt.Printf("Failed to create reddit script: %v\n", err)
	} else if _, wait, err := graw.Scan(&announcer{}, script, cfg); err != nil {
		fmt.Printf("graw launch failed: %v\n", err)
	} else if err := wait(); err != nil {
		fmt.Printf("graw run failed: %v\n", err)
	}

}
