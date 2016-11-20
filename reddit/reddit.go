// Package reddit is a small wrapper around the Reddit API.
//
// It provides two handles for the Reddit API. Logged out users can claim a
// Script handle.
//
//   rate := 5 * time.Second
//   script, _ := NewScript("graw:doc_script:0.3.1 by /u/yourusername", rate)
//   post, _ := script.Thread("r/programming/comments/5du93939")
//   fmt.Printf("%s posted \"%s\"!", post.Author, post.Title)
//
// Logged in users can claim a Bot handle with a superset of the former's
// features.
//
//   cfg := BotConfig{
//     Agent: "graw:doc_demo_bot:0.3.1 by /u/yourusername"
//     // Your registered app info from following:
//     // https://github.com/reddit/reddit/wiki/OAuth2
//     App: App{
//       ID:     "sdf09ofnsdf",
//       Secret: "skldjnfksjdnf",
//       Username: "yourbotusername",
//       Password: "yourbotspassword",
//     }
//   }
//   bot, _ := NewBot(cfg)
//   bot.SendMessage("roxven", "Thanks for making this Reddit API!", "It's ok.")
//
// Requests made by this API are rate limited with no bursting. All interfaces
// exported by this package have goroutine safe implementations, but when shared
// by many goroutines some calls may block for multiples of the rate limit
// interval.
//
// This API for accessing feeds from Reddit is low level, built specifically for
// graw. If you are interested in a simple high level event feed, see graw.
package reddit
