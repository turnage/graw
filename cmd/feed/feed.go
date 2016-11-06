// Feed provides the feed of a subreddit.
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/turnage/graw/internal/api"
	"github.com/turnage/graw/internal/client"
	"github.com/turnage/graw/internal/data"
	"github.com/turnage/graw/internal/monitor"
	"github.com/turnage/graw/internal/reap"
	"github.com/turnage/graw/internal/rsort"
)

const userAgent = "cli:graw_feed_tool:0.1.0 (by /u/roxven)"

var subreddit = flag.String("subreddit", "", "Subreddit to subscribe to.")

func main() {
	flag.Parse()

	if *subreddit == "" {
		fmt.Printf("You must provide a subreddit! Use -h.")
		return
	}

	cli, err := client.New(client.Config{Agent: userAgent})
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}

	m, err := monitor.New(
		monitor.Config{
			Path: "/r/" + *subreddit + "/new.json",
			Lurker: api.NewLurker(
				reap.New(
					reap.Config{
						Client:   cli,
						Parser:   data.NewParser(),
						Hostname: "www.reddit.com",
						TLS:      true,
					},
				),
			),
			Sorter: rsort.New(),
		},
	)
	if err != nil {
		log.Fatalf("Failed to create monitor: %v\n", err)
	}

	for _ = range time.Tick(time.Second) {
		harvest, err := m.Update()
		if err != nil {
			log.Fatalf("Monitoring error: %v\n", err)
		}

		for _, p := range harvest.Posts {
			log.Printf("[%s]: %s\n", p.Author, p.Title)
		}
	}
}
