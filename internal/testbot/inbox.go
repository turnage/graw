package main

import (
	"log"
	"strings"

	"github.com/turnage/graw"
)

func (b *bot) Mention(m *graw.Message) error {
	log.Printf("Mentioned by %s in %s.\n", m.Author, m.Subreddit)
	return nil
}

func (b *bot) PostReply(m *graw.Message) error {
	log.Printf("Received reply from %s on %s.\n", m.Author, m.LinkTitle)
	return nil
}

func (b *bot) CommentReply(m *graw.Message) error {
	log.Printf("Received reply from %s on one of our comments.\n", m.Author)
	return nil
}

func (b *bot) Message(m *graw.Message) error {
	log.Printf("Received message from %s.\n", m.Author)

	if strings.HasPrefix(m.Subject, "cmd:") {
		return b.exec(
			strings.TrimPrefix(m.Subject, "cmd:"),
			strings.Split(m.Body, ":"),
		)
	}

	return nil
}
