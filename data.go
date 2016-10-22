// Data types defined here all derive from Reddit's definitions. See the Reddit
// documentation for more context: https://github.com/reddit/reddit/wiki/JSON
package graw

// Comment represents a comment on Reddit (Reddit type t1_).
type Comment struct {
	ID        string `mapstructure:"id"`
	Name      string `mapstructure:"name"`
	Permalink string `mapstructure:"permalink"`

	CreatedUTC float64 `mapstructure:"created_utc"`
	Deleted    bool    `mapstructure:"deleted"`

	Ups   int32 `mapstructure:"ups"`
	Downs int32 `mapstructure:"downs"`
	Likes bool  `mapstructure:"likes"`

	Author              string `mapstructure:"author"`
	AuthorFlairCssClass string `mapstructure:"author_flair_css_class"`
	AuthorFlairText     string `mapstructure:"author_flair_text"`

	LinkAuthor string `mapstructure:"link_author"`
	LinkURL    string `mapstructure:"link_url"`

	Subreddit   string `mapstructure:"subreddit"`
	SubredditID string `mapstructure:"subreddit_id"`

	Body     string `mapstructure:"body"`
	BodyHtml string `mapstructure:"body_html"`

	ParentID string     `mapstructure:"parent_id"`
	Replies  []*Comment `mapstructure:"reply_tree"`

	Gilded        int32  `mapstructure:"gilded"`
	Distinguished string `mapstructure:"distinguished"`

	// This is present if the comment appears in a user's inbox.
	Subject string `mapstructure:"subject"`
}

// Post represents posts on Reddit (Reddit type t3_).
type Post struct {
	ID        string `mapstructure:"id"`
	Name      string `mapstructure:"name"`
	Permalink string `mapstructure:"permalink"`

	CreatedUTC float64 `mapstructure:"created_utc"`
	Deleted    bool    `mapstructure:"deleted"`

	Ups   int32 `mapstructure:"ups"`
	Downs int32 `mapstructure:"downs"`
	Likes bool  `mapstructure:"likes"`

	Author              string `mapstructure:"author"`
	AuthorFlairCssClass string `mapstructure:"author_flair_css_class"`
	AuthorFlairText     string `mapstructure:"author_flair_text"`

	Title  string `mapstructure:"title"`
	Score  int32  `mapstructure:"score"`
	URL    string `mapstructure:"url"`
	Domain string `mapstructure:"domain"`
	NSFW   bool   `mapstructure:"over_18"`

	Subreddit   string `mapstructure:"subreddit"`
	SubredditID string `mapstructure:"subreddit_id"`

	IsSelf       bool   `mapstructure:"is_self"`
	SelfText     string `mapstructure:"selftext"`
	SelfTextHtml string `mapstructure:"selftext_html"`

	Replies []*Comment `mapstructure:"reply_tree"`

	Hidden            bool   `mapstructure:"hidden"`
	LinkFlairCssClass string `mapstructure:"link_flair_css_class"`
	LinkFlairText     string `mapstructure:"link_flair_text"`

	NumComments int32  `mapstructure:"num_comments"`
	Locked      bool   `mapstructure:"locked"`
	Thumbnail   string `mapstructure:"thumbnail"`

	Gilded        int32  `mapstructure:"gilded"`
	Distinguished string `mapstructure:"distinguished"`
	Stickied      bool   `mapstructure:"stickied"`
}

type Message struct {
	Author   string `mapstructure:"author"`
	Body     string `mapstructure:"body"`
	BodyHtml string `mapstructure:"body_html"`
	Context  string `mapstructure:"context"`

	FirstMessageName string `mapstructure:"first_message_name"`
	Likes            bool   `mapstructure:"likes"`
	LinkTitle        string `mapstructure:"link_title"`

	New      bool   `mapstructure:"new"`
	ParentID string `mapstructure:"parent_id"`
	Subject  string `mapstructure:"subject"`

	Subreddit  string `mapstructure:"subreddit"`
	WasComment bool   `mapstructure:"was_comment"`

	Messages []*Message `mapstructure:"messages"`
}
