package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/turnage/graw"
)

const usage = "Usage of stalker: ./stalker [user] [user] ...\n"

const agent = "graw:stalker-demo-bot:1.0.0 by /u/roxven"

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
		os.Exit(0)
	}
}

type stalker struct{}

func (a *stalker) UserPost(p *graw.Post) error {
	fmt.Printf(
		"[%s in %s]: Posted \"%s\"\n",
		p.Author, p.Subreddit, p.Title,
	)
	return nil
}

func (a *stalker) UserComment(c *graw.Comment) error {
	fmt.Printf(
		"[%s in %s]: Commented on \"%s\"\n",
		c.Author, c.Subreddit, c.LinkTitle,
	)
	return nil
}

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(-1)
	}

	fmt.Printf(
		"Error: %v\n", graw.Run(
			graw.Config{
				Agent: agent,
				Users: flag.Args(),
			},
			&stalker{},
		),
	)
}
