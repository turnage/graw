// Package act is a utility to operating a Reddit account from cli using the
// reddit package.
package main

import (
	"log"
	"os"

	"github.com/turnage/graw/reddit"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app   = kingpin.New("act", "A cli tool for acting as a Reddit account.")
	agent = app.Flag("agent", "Filename of the agent file to use.").Required().String()

	postlink      = app.Command("postlink", "Post a link to a subreddit.")
	linkSubreddit = postlink.Arg("subreddit", "Subredit to post in.").Required().String()
	linkTitle     = postlink.Arg("title", "Title of the link post.").Required().String()
	linkURL       = postlink.Arg("url", "URL to post.").Required().String()

	postself      = app.Command("postself", "Post a text post to a subreddit.")
	selfSubreddit = postself.Arg("subreddit", "Subredit to post in.").Required().String()
	selfTitle     = postself.Arg("title", "Title of the self post.").Required().String()
	selfText      = postself.Arg("text", "Text to post.").Required().String()

	reply     = app.Command("reply", "Post a reply to a reddit element.")
	replyName = reply.Arg("name", "Full name of the origin element in the form t#_xxxxx").Required().String()
	replyText = reply.Arg("text", "Text to post.").Required().String()

	send     = app.Command("send", "Send a message to a user.")
	sendTo   = send.Arg("recipient", "User to send a message to.").Required().String()
	sendSubj = send.Arg("subject", "Subject of the message.").Required().String()
	sendBody = send.Arg("body", "Body of the message.").Required().String()
)

func bot(agentfile string) reddit.Bot {
	bot, err := reddit.NewBotFromAgentFile(agentfile, 0)
	if err != nil {
		log.Fatalf("Failed to create api handle: %v\n", err)
	}

	return bot
}

func maybeFail(err error) {
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case postlink.FullCommand():
		b := bot(*agent)
		maybeFail(b.PostLink(*linkSubreddit, *linkTitle, *linkURL))
	case postself.FullCommand():
		b := bot(*agent)
		maybeFail(b.PostSelf(*selfSubreddit, *selfTitle, *selfText))
	case reply.FullCommand():
		b := bot(*agent)
		maybeFail(b.Reply(*replyName, *replyText))
	case send.FullCommand():
		b := bot(*agent)
		maybeFail(b.SendMessage(*sendTo, *sendSubj, *sendBody))
	}
}
