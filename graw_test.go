package graw

import (
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"github.com/turnage/graw/botfaces"
	"github.com/turnage/redditproto"
)

const (
	agent1 = "res/one"
	agent2 = "res/two"
)

var (
	user1 = flag.String("user1", "", "username of the first bot")
	user2 = flag.String("user2", "", "username of the second bot")
	sub   = flag.String("subreddit", "", "subreddit to test in")
)

var (
	one = newTestBot()
	two = newTestBot()
)

type testBot struct {
	commentReplies chan *redditproto.Comment
	postReplies    chan *redditproto.Comment
	userPosts      chan *redditproto.Link
	userComments   chan *redditproto.Comment
	posts          chan *redditproto.Link
	messages       chan *redditproto.Message
	comments       chan *redditproto.Comment
	mentions       chan *redditproto.Comment
	setupCalls     chan bool
	teardownCalls  chan bool
	failures       chan error
	eng            Engine
}

func newTestBot() *testBot {
	return &testBot{
		commentReplies: make(chan *redditproto.Comment),
		postReplies:    make(chan *redditproto.Comment),
		userPosts:      make(chan *redditproto.Link),
		userComments:   make(chan *redditproto.Comment),
		posts:          make(chan *redditproto.Link),
		messages:       make(chan *redditproto.Message),
		comments:       make(chan *redditproto.Comment),
		mentions:       make(chan *redditproto.Comment),
		setupCalls:     make(chan bool),
		teardownCalls:  make(chan bool),
		failures:       make(chan error),
	}
}

func (t *testBot) SetUp() error {
	t.setupCalls <- true
	t.eng = GetEngine(t)
	return nil
}

func (t *testBot) TearDown() {
	t.teardownCalls <- true
}

func (t *testBot) Fail(err error) bool {
	t.failures <- err
	return true
}

func (t *testBot) BlockTime() time.Duration {
	return 2 * time.Second
}

func (t *testBot) CommentReply(reply *redditproto.Comment) {
	t.commentReplies <- reply
}

func (t *testBot) PostReply(reply *redditproto.Comment) {
	t.postReplies <- reply
}

func (t *testBot) Mention(mention *redditproto.Comment) {
	t.mentions <- mention
}

func (t *testBot) Message(msg *redditproto.Message) {
	t.messages <- msg
}

func (t *testBot) Post(post *redditproto.Link) {
	t.posts <- post
}

func (t *testBot) UserPost(post *redditproto.Link) {
	t.userPosts <- post
}

func (t *testBot) UserComment(comment *redditproto.Comment) {
	t.userComments <- comment
}

func TestSelfPost(t *testing.T) {
	t.Parallel()
	if err := one.eng.SelfPost(*sub, "Test Self Post", "body"); err != nil {
		t.Fatalf("/u/%s failed to make a self post in /r/%s: %v", *user1, *sub, err)
	} else {
		t.Logf("/u/%s made a self post in /r/%s.\n", *user1, *sub)
	}
}

func TestReceivePostFromWatchedUser(t *testing.T) {
	t.Parallel()
	select {
	case <-two.userPosts:
	case <-time.After(2 * time.Minute):
	}
}

func TestReceivePostFromWathedSubreddit(t *testing.T) {
	t.Parallel()
	select {
	case <-two.posts:
	case <-time.After(2 * time.Minute):
	}
}

func TestSendMessage(t *testing.T) {
	t.Parallel()
	if err := one.eng.SendMessage(*user2, "test", "different"); err != nil {
		t.Fatalf("/u/%s failed to send a message to /u/%s: %v", *user1, *user2, err)
	} else {
		t.Logf("/u/%s failed to send a message to /u/%s.\n", *user1, *user2)
	}
}

func TestReceiveMessage(t *testing.T) {
	t.Parallel()
	select {
	case <-two.messages:
	case <-time.After(2 * time.Minute):
	}
}

func TestMain(m *testing.M) {
	flag.Parse()

	var bot interface{} = one
	if _, ok := bot.(botfaces.Loader); !ok {
		log.Panic("Test bot does not implement Loader.")
	} else if _, ok := bot.(botfaces.Tearer); !ok {
		log.Panic("Test bot does not implement Tearer.")
	} else if _, ok := bot.(botfaces.PostHandler); !ok {
		log.Panic("Test bot does not implement PostHandler.")
	} else if _, ok := bot.(botfaces.CommentReplyHandler); !ok {
		log.Panic("Test bot does not implement CommentReplyHandler.")
	} else if _, ok := bot.(botfaces.PostReplyHandler); !ok {
		log.Panic("Test bot does not implement PostReplyHandler.")
	} else if _, ok := bot.(botfaces.MessageHandler); !ok {
		log.Panic("Test bot does not implement MessageHandler.")
	} else if _, ok := bot.(botfaces.BlockTimer); !ok {
		log.Panic("Test bot does not implement BlockTimer.")
	} else if _, ok := bot.(botfaces.MentionHandler); !ok {
		log.Panic("Test bot does not implement MentionHandler.")
	} else if _, ok := bot.(botfaces.UserHandler); !ok {
		log.Panic("Test bot does not implement UserHandler.")
	}

	if *sub == "" || *user1 == "" || *user2 == "" {
		os.Exit(0)
	}

	runner := func(agent string, bot interface{}, sub string) {
		if err := Run(agent, bot, sub); err != nil {
			log.Printf("Bot %s encountered an error: %v\n", agent, err)
			os.Exit(-1)
		}
	}

	go runner(agent1, one, *sub)
	go runner(agent2, two, *sub)

	<-one.setupCalls
	<-two.setupCalls

	if err := one.eng.WatchUser(*user2); err != nil {
		log.Panicf("/u/%s failed to watch /u/%s.\n", *user1, *user2)
	}

	if err := two.eng.WatchUser(*user1); err != nil {
		log.Panicf("/u/%s failed to watch /u/%s.\n", *user2, *user1)
	}

	os.Exit(m.Run())
}
