// Package graw is a high level, easy to use, Reddit bot library.
//
// graw will take a low level handle from the graw/reddit package and manage
// everything for you. You just specify in a config what events you want to
// listen to on Reddit, and graw will take care of maintaining the event stream
// and calling the handler methods of your bot whenever new events occur.
//
// Announcing new posts in /r/self is as simple as:
//
//   type announcer struct{}
//
//   func (a *announcer) Post(post *reddit.Post) error {
//     fmt.Printf(`%s posted "%s"\n`, post.Author, post.Title)
//     return nil
//   }
//
//   .....
//
//   // Get an api handle to reddit for a logged out (script) program,
//   // which forwards this user agent on all requests and issues a request at
//   // most every 5 seconds.
//   apiHandle := reddit.NewScript("your user agent", 5 * time.Second)
//
//   // Create a configuration specifying what event sources on Reddit graw
//   // should connect to the bot.
//   cfg := graw.Config{Subreddits: []string{"self"}}
//
//   // launch a graw scan in a goroutine using the bot, handle, and config. The
//   // returned "stop" and "wait" are functions. "stop" will stop the graw run
//   // at any time, and "wait" will block until it finishes.
//   stop, wait, err := graw.Scan(&announcer{}, apiHandle, cfg)
//
//   // This time, let's block so the bot will announce (ideally) forever.
//   if err := wait(); err != nil {
//     fmt.Printf("graw run encountered an error: %v\n", err)
//   }
//
// graw can handle many event sources on Reddit; see the Config struct for the
// complete set of offerings.
//
// graw has one other function that behaves like Scan(), which is Run(). Scan()
// is for logged-out bots (what Reddit calls "scripts"). Run() handles logged in
// bots, which can subscribe to logged-in event sources in the bot's account
// inbox like mentions and private messages.
package graw
