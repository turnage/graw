graw
--------------------------------------------------------------------------------

Status: Pilot Program

![Build Status](https://travis-ci.org/turnage/graw.svg?branch=master)
[![Coverage Status](https://coveralls.io/repos/turnage/graw/badge.svg?branch=master&service=github)](https://coveralls.io/github/turnage/graw?branch=master)

    go get github.com/turnage/graw

graw
* is for writing Reddit bots that run forever and all time.
* is for writing Reddit bots fast without worrying about things like "loops".
* is for writing Reddit bots *in Go*!

Choose what events on Reddit to listen for (e.g. private messages, or new posts 
in certain subreddits) and graw will feed them to your bot. Here is a simple
graw bot that announces new posts:

    type AnnouncerBot struct {}
    
    func (a *AnnouncerBot) Post(post *redditproto.Link) {
        fmt.Printf("New post by %s: %s\n", post.GetAuthor(), post.GetTitle())
    }

graw provides all data from Reddit in the form of
[Protocol Buffers](https://developers.google.com/protocol-buffers/).
See graw's [proto definitions](https://github.com/turnage/redditproto/blob/master/reddit.proto).

See the [getting started](https://github.com/turnage/graw/wiki/Getting-Started)
page to see if graw is for you, and how to start using it.
