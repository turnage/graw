// Package feed is an example grawbot that announces the feed of a given
// subreddit.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mix/graw"

	"github.com/mix/graw/reddit"
)

var (
	feed           = kingpin.New("feed", "A cli tool for announcing Reddit feeds.")
	agent          = feed.Flag("agent", "Filename of the agent file to use.").String()
	rate           = feed.Flag("rate", "Update interval.").Duration()
	subreddits     = feed.Flag("subreddits", "Subreddits to announce.").Strings()
	comments       = feed.Flag("comments", "Subreddits to announce comments in.").Strings()
	users          = feed.Flag("users", "Users to announce activity from.").Strings()
	postreplies    = feed.Flag("postreplies", "Announce replies to bot's posts.").Bool()
	commentreplies = feed.Flag("commentreplies", "Announce replies to bot's comments.").Bool()
	mentions       = feed.Flag("mentions", "Announce mentions of the bot's username.").Bool()
	messages       = feed.Flag("messages", "Announce messages sent to the bot.").Bool()
)

type announcer struct{}

func (a *announcer) Post(p *reddit.Post) error {
	fmt.Printf(
		"[New Post in %s][by %s]: %s\n",
		p.Subreddit, p.Author, p.Title,
	)
	return nil
}

func (a *announcer) Comment(c *reddit.Comment) error {
	fmt.Printf(
		"[New Comment in %s][by %s]: %s\n\n",
		c.Subreddit, c.Author, c.Body,
	)
	return nil
}

func (a *announcer) UserPost(p *reddit.Post) error {
	fmt.Printf(
		"[Watched user %s][posted in %s]: %s\n",
		p.Author, p.Subreddit, p.Title,
	)
	return nil
}

func (a *announcer) UserComment(c *reddit.Comment) error {
	fmt.Printf(
		"[Watched user %s][commented in %s]: %s\n",
		c.Author, c.Subreddit, c.LinkTitle,
	)
	return nil
}

func (a *announcer) PostReply(m *reddit.Message) error {
	fmt.Printf(
		"[Post reply to %s][by %s]: %s\n",
		m.LinkTitle, m.Author, m.Body,
	)
	return nil
}

func (a *announcer) CommentReply(m *reddit.Message) error {
	fmt.Printf(
		"[Comment reply in thread %s][by %s]: %s\n",
		m.LinkTitle, m.Author, m.Body,
	)
	return nil
}

func (a *announcer) Mention(m *reddit.Message) error {
	fmt.Printf(
		"[Username mention in thread %s][by %s]: %s\n",
		m.LinkTitle, m.Author, m.Body,
	)
	return nil
}

func (a *announcer) Message(m *reddit.Message) error {
	fmt.Printf(
		"[Message from %s][%s]: %s\n",
		m.Author, m.Subject, m.Body,
	)
	return nil
}

func bot(agentfile string) reddit.Bot {
	bot, err := reddit.NewBotFromAgentFile(agentfile, 0)
	if err != nil {
		log.Fatalf("Failed to create api handle: %v\n", err)
	}

	return bot
}

func main() {
	kingpin.MustParse(feed.Parse(os.Args[1:]))

	if *agent == "" && (*postreplies || *commentreplies || *mentions || *messages) {
		fmt.Printf("You must provide an agent file to subscribe to the inbox.")
		os.Exit(-1)
	}

	cfg := graw.Config{
		Subreddits:        *subreddits,
		SubredditComments: *comments,
		Users:             *users,
		PostReplies:       *postreplies,
		CommentReplies:    *commentreplies,
		Messages:          *messages,
		Mentions:          *mentions,
		Logger:            log.New(os.Stderr, "", log.LstdFlags),
	}

	var err error
	var wait func() error
	if *agent != "" {
		if _, wait, err = graw.Run(&announcer{}, bot(*agent), cfg); err != nil {
			log.Fatalf("Failed to launch graw run: %v\n", err)
		}
	} else {
		if script, err := reddit.NewScript(
			"graw:feed demo bot:0.5.1 by /u/roxven",
			*rate,
		); err != nil {
			log.Fatalf("Failed to create reddit script: %v\n", err)
		} else if _, wait, err = graw.Scan(&announcer{}, script, cfg); err != nil {
			log.Fatalf("graw launch failed: %v\n", err)
		}
	}

	if err := wait(); err != nil {
		log.Fatalf("graw run failed: %v\n", err)
	}

}
