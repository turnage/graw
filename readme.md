# graw

_A Reddit API wrapper, for Go._

![Build Status](https://travis-ci.org/turnage/graw.svg?branch=master)
![Version: 1.2.0](https://img.shields.io/badge/version-1.2.0-brightgreen.svg)
[![GoDoc](https://godoc.org/github.com/turnage/graw?status.svg)](https://godoc.org/github.com/turnage/graw)

> This is a personally maintained fork of
> [`turnage/graw`](https://github.com/turnage/graw).
>
> All attributions to this project should really be going to
> [`turnage](https://github.com/turnage).

`graw` is a library for building Reddit bots that takes care of everything you
don't want to. [Read the tutorial book][gitbook]!

```bash
go get -u github.com/turnage/graw
```

## Usage

The design of `graw` is that your bot is a handler for events, Reddit is a
source of events, and `graw` connects the two. If you want to announce all the
new posts in a given subreddit, this is your bot:

```go
type Announcer struct {}

func (a *Announcer) Post(post *reddit.Post) error {
	fmt.Printf(`%s posted "%s"\n`, post.Author, post.Title
	return nil
}
```

Give this to `graw` with an
[API handle from the `reddit` package](https://godoc.org/github.com/turnage/graw/reddit)
and tell it what events you want to subscribe to; `graw` will take care of the
rest. See the [godoc][godoc] and [tutorial book][gitbook] for more information.

## Features

The primary feature of `graw` is robust event streams. `graw` supports many
exciting event streams:

- New posts in subreddits.
- New comments in subreddits.
- New posts or comments by users.
- Private messages sent to the bot.
- Replies to the bot's posts.
- Replies to the bot's comments.
- Mentions of the bot's username.

Processing all of these events is as as simple as implementing a method to
receive them!

`graw` also provides two lower level packages for developers to tackle other
interactions with Reddit like one-shot scripts and bot actions. See
subdirectories in the [godoc][godoc].

[godoc]: https://godoc.org/github.com/stevenxie/graw
[gitbook]: https://turnage.gitbooks.io/graw/content/
