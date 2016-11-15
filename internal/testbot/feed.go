package main

import (
	"log"

	"github.com/turnage/graw"
)

func (b *bot) Post(p *graw.Post) error {
	log.Printf("%s posted \"%s\" in %s.\n", p.Author, p.Title, p.Subreddit)
	return nil
}
