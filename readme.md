graw
--------------------------------------------------------------------------------

![Build Status](https://travis-ci.org/turnage/graw.svg?branch=master)
![Version: 1.2.0](https://img.shields.io/badge/version-1.2.0-brightgreen.svg)
[![GoDoc](https://godoc.org/github.com/turnage/graw?status.svg)](https://godoc.org/github.com/turnage/graw)

    go get github.com/turnage/graw

graw is a library for building Reddit bots that takes care of everything you
don't want to. [Read the tutorial book](https://turnage.gitbooks.io/graw/content/)!

As of major version 1, the API promise is: no breaking changes, ever. Details
below. This applies to all (library) subpackages of graw.

### Usage

The design of graw is that your bot is a handler for events, Reddit is a source
of events, and graw connects the two. If you want to announce all the new posts
in a given subreddit, this is your bot:

````go
type announcer struct {}

func (a *announcer) Post(post *reddit.Post) error {
        fmt.Printf(`%s posted "%s"\n`, post.Author, post.Title)
        return nil
}
````

Give this to graw with an
[api handle from the reddit package](https://godoc.org/github.com/turnage/graw/reddit)
and a tell it what events you want to subscribe to; graw will take care of the
rest. See the [godoc](https://godoc.org/github.com/turnage/graw) and
[tutorial book](https://turnage.gitbooks.io/graw/content/) for more information.

### Features

The primary feature of graw is robust event streams. graw supports many exciting
event streams:

* New posts in subreddits.
* New comments in subreddits.
* New posts or comments by users.
* Private messages sent to the bot.
* Replies to the bot's posts.
* Replies to the bot's comments.
* Mentions of the bot's username.

Processing all of these events is as as simple as implementing a method to
receive them!

graw also provides two lower level packages for developers to tackle other
interactions with Reddit like one-shot scripts and bot actions. See
subdirectories in the godoc.

### API Promise

As of version 1.0.0, the graw API is stable. I will not make any backwards
incompatible changes, ever. The only exceptions are:

* I may add methods to an interface. This will only break you if you embed it
  and implement a method with the same name as the one I add.
* I may add fields to the Config struct. This will only break you if you embed
  it and add a field with the same name as the one I add, or initialize it
  positionally.

I don't foresee anyone having a reason to do either of these things. 
