graw
--------------------------------------------------------------------------------

Version: 0.2.2

Before depending on graw, please consult the [API Promise](promise.md).

![Build Status](https://travis-ci.org/turnage/graw.svg?branch=master)
[![Coverage Status](https://coveralls.io/repos/turnage/graw/badge.svg?branch=master&service=github)](https://coveralls.io/github/turnage/graw?branch=master)
[![GoDoc](https://godoc.org/github.com/turnage/graw?status.svg)](https://godoc.org/github.com/turnage/graw)

    go get github.com/turnage/graw

graw is for writing Reddit bots
* that run forever and all the time.
* quickly without worrying about things like "loops".
* *in Go!*

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

See the [wiki](https://github.com/turnage/graw/wiki) for a quick start.

Here is an [example grawbot](https://gist.github.com/turnage/468f981f3b1e85bb19f2#file-announcer-go) that announces all of the new posts in /r/all.

Here is an [example grawbot](https://gist.github.com/turnage/468f981f3b1e85bb19f2#file-replier-go) that
automatically replies to private messages.
