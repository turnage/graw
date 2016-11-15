package main

import (
	"log"

	"github.com/turnage/graw"
)

func (b *bot) UserPost(p *graw.Post) error {
	log.Printf("Stalked user %s posted in %s.\n", p.Author, p.Subreddit)
	return nil
}

func (b *bot) UserComment(c *graw.Comment) error {
	log.Printf("Stalked user %s commented in %s.\n", c.Author, c.Subreddit)
	return nil
}
