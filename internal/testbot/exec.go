package main

import (
	"log"
)

func (b *bot) exec(cmd string, args []string) error {
	switch cmd {
	case "reply":
		// args[0]: name, args[1]: text
		return b.Reply(args[0], args[1])
	case "postself":
		// args[0]: subreddit, args[1]: title, args[2]:text
		return b.PostSelf(args[0], args[1], args[2])
	case "postlink":
		// args[0]: subreddit, args[1]: title, args[2]:url
		return b.PostLink(args[0], args[1], args[2])
	case "send":
		// args[0]: user, args[1]: subject, args[2]: text
		return b.SendMessage(args[0], args[1], args[2])
	case "stop":
		{
			b.Stop()
			return nil
		}
	}

	log.Printf("Did not recognize command: %s, args: %v.\n", cmd, args)
	return nil
}
