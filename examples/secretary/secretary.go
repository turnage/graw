package main

import (
	"fmt"

	"github.com/turnage/graw"
)

const agent = "graw:secretary-demo-bot:1.0.0 by /u/roxven"

type secretary struct{}

func (s *secretary) CommentReply(m *graw.Message) error {
	fmt.Printf("Received a comment reply from %s: %s\n", m.Author, m.Body)
	return nil
}

func (s *secretary) PostReply(m *graw.Message) error {
	fmt.Printf("Received a post reply from %s: %\ns", m.Author, m.Body)
	return nil
}

func (s *secretary) Mention(m *graw.Message) error {
	fmt.Printf("Received a mention from %s: %s\n", m.Author, m.Body)
	return nil
}

func (s *secretary) Message(m *graw.Message) error {
	fmt.Printf("Received a message from %s: %s\n", m.Author, m.Body)
	return nil
}

func main() {
	fmt.Printf(
		"Error: %v\n", graw.Run(
			graw.Config{
				Agent: agent,
				Inbox: true,
				// USE `git add -p` when staging this file! Do
				// not check these values in!!
				App: &graw.App{
					ID:       "",
					Secret:   "",
					Username: "",
					Password: "",
				},
			},
			&secretary{},
		),
	)
}
