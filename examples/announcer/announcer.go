package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/turnage/graw"
)

const usage = "Usage of announcer: ./announcer [subreddit] [subreddit] ...\n"

const agent = "graw:announcer-demo-bot:1.0.0 by /u/roxven"

var rate = flag.Duration("rate", time.Second, "Interval between updates.")

func init() {
	flag.Usage = func() { fmt.Fprintf(os.Stderr, usage) }
}

type announcer struct{}

func (a *announcer) Post(p *graw.Post) error {
	fmt.Printf("[%s]: \"%s\"\n", p.Subreddit, p.Title)
	return nil
}

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
	}

	fmt.Printf(
		"Error: %v\n", graw.Run(
			graw.Config{
				Agent:      agent,
				Subreddits: flag.Args(),
				Rate:       *rate,
			},
			&announcer{},
		),
	)
}
