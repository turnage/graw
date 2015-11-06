package graw

import (
	"fmt"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/turnage/graw/internal/operator"
	"github.com/turnage/redditproto"
)

// End to end test configuration.
const (
	agentFile        = "test.agent"
	testSubredditKey = "GRAW_TEST_SUBREDDIT"
	testRedditKey    = "GRAW_TEST_REDDIT"
)

const (
	// eventTimeout is the amount of time that a test should wait for an
	// event to register with a bot.
	eventTimeout = time.Second * 20
)

// Globals for tests.
var (
	bot           = &testBot{}
	testSubreddit = ""
	agent         = &redditproto.UserAgent{}
)

type testBot struct {
	eng            Engine
	posts          uint64
	userPosts      uint64
	userComments   uint64
	messages       uint64
	firstMessage   *redditproto.Message
	postReplies    uint64
	commentReplies uint64
	mentions       uint64
	events         uint64
}

func (t *testBot) SetUp() error {
	t.eng = GetEngine(t)
	return nil
}

func (t *testBot) TearDown() {}

func (t *testBot) Post(post *redditproto.Link) {
	atomic.AddUint64(&t.posts, 1)
	atomic.AddUint64(&t.events, 1)
}

func (t *testBot) UserPost(post *redditproto.Link) {
	atomic.AddUint64(&t.userPosts, 1)
	atomic.AddUint64(&t.events, 1)
}

func (t *testBot) UserComment(comment *redditproto.Comment) {
	atomic.AddUint64(&t.userComments, 1)
	atomic.AddUint64(&t.events, 1)
}

func (t *testBot) Message(message *redditproto.Message) {
	if atomic.LoadUint64(&t.messages) == 0 {
		t.firstMessage = message
	}
	atomic.AddUint64(&t.messages, 1)
	atomic.AddUint64(&t.events, 1)
}

func (t *testBot) PostReply(comment *redditproto.Comment) {
	atomic.AddUint64(&t.postReplies, 1)
	atomic.AddUint64(&t.events, 1)
}

func (t *testBot) CommentReply(comment *redditproto.Comment) {
	atomic.AddUint64(&t.commentReplies, 1)
	atomic.AddUint64(&t.events, 1)
}

func (t *testBot) Mention(comment *redditproto.Comment) {
	atomic.AddUint64(&t.mentions, 1)
	atomic.AddUint64(&t.events, 1)
}

// e2eParams returns whether the environment is configured for end to end tests,
// and the parameters for them.
func e2eParams() (bool, string, string) {
	reddit := os.Getenv(testRedditKey)
	subreddit := os.Getenv(testSubredditKey)
	if subreddit == "" || reddit == "" {
		return false, "", ""
	}

	return true, reddit, subreddit
}

func launchBot(errors chan<- error, agent string) *testBot {
	bot := &testBot{}
	go func() {
		if err := Run(agent, bot, testSubreddit); err != nil {
			errors <- err
		}
	}()
	return bot
}

// wait returns true if a condition evaluates to true within the timeout.
func wait(condition func() bool, timeout time.Duration) bool {
	kill := make(chan bool)
	event := make(chan bool)
	go func() {
		stop := false
		for !stop {
			select {
			case <-kill:
				stop = true
			case <-time.After(10 * time.Millisecond):
				if condition() {
					event <- true
				}
			}
		}
	}()
	defer func() { kill <- true }()

	select {
	case <-event:
		return true
	case <-time.After(timeout):
		return false
	}

	return false
}

func TestPostStreamAndSubmissions(t *testing.T) {
	if err := bot.eng.SelfPost(testSubreddit, "title", "content"); err != nil {
		t.Fatal(err)
	}
	if !wait(
		func() bool { return atomic.LoadUint64(&bot.posts) == 1 },
		eventTimeout,
	) {
		t.Errorf("the new post did not register with the bot")
	}
}

func TestInboxAndReply(t *testing.T) {
	if err := bot.eng.SendMessage(agent.GetUsername(), "subject", "text"); err != nil {
		t.Fatal(err)
	}
	if !wait(
		func() bool { return atomic.LoadUint64(&bot.messages) == 1 },
		eventTimeout,
	) {
		t.Fatalf("the new message did not register with the bot")
	}
	if err := bot.eng.Reply(bot.firstMessage.GetName(), "text"); err != nil {
		t.Fatal(err)
	}
	if !wait(
		func() bool { return atomic.LoadUint64(&bot.messages) == 2 },
		eventTimeout,
	) {
		t.Errorf("the new message did not register with the bot")
	}
}

func TestUserWatch(t *testing.T) {
	if err := bot.eng.WatchUser(agent.GetUsername()); err != nil {
		t.Fatal(err)
	}
	if err := bot.eng.SelfPost(testSubreddit, "title", "content"); err != nil {
		t.Fatal(err)
	}
	if !wait(
		func() bool { return atomic.LoadUint64(&bot.userPosts) == 1 },
		eventTimeout,
	) {
		t.Fatalf("the new watched user post did not register with the bot")
	}
	if err := bot.eng.SelfPost(testSubreddit, "title", "content"); err != nil {
		t.Fatal(err)
	}
	if err := bot.eng.UnwatchUser(agent.GetUsername()); err != nil {
		t.Fatal(err)
	}
	if wait(
		func() bool { return atomic.LoadUint64(&bot.userPosts) == 2 },
		eventTimeout,
	) {
		t.Errorf("an event was generated from an unwatched user")
	}
}

func TestMain(m *testing.M) {
	var configured bool
	var domain string
	configured, domain, testSubreddit = e2eParams()
	if !configured {
		fmt.Printf("End to end tests not configured; skipping.\n")
		os.Exit(0)
	}

	operator.SetTestDomain(domain)
	errors := make(chan error)
	go func() {
		err := <-errors
		fmt.Printf("Error: %v\n", err)
		os.Exit(-1)
	}()

	var err error
	agent, err = redditproto.Load(agentFile)
	if err != nil {
		fmt.Printf("Failed to load test user agent file.\n")
		os.Exit(-1)
	}

	bot = launchBot(errors, agentFile)
	if !wait(func() bool { return bot.eng != nil }, eventTimeout) {
		fmt.Printf("the bot did not receive the engine in time")
		os.Exit(-1)
	}

	defer bot.eng.Stop()
	os.Exit(m.Run())
}
