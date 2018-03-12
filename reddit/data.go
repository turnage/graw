package reddit

import "strings"

// Comment represents a comment on Reddit (Reddit type t1_).
// https://github.com/reddit/reddit/wiki/JSON#comment-implements-votable--created
type Comment struct {
	ID        string `mapstructure:"id"`
	Name      string `mapstructure:"name"`
	Permalink string `mapstructure:"permalink"`

	CreatedUTC uint64 `mapstructure:"created_utc"`
	Deleted    bool   `mapstructure:"deleted"`

	Ups   int32 `mapstructure:"ups"`
	Downs int32 `mapstructure:"downs"`
	Likes bool  `mapstructure:"likes"`

	Author              string `mapstructure:"author"`
	AuthorFlairCSSClass string `mapstructure:"author_flair_css_class"`
	AuthorFlairText     string `mapstructure:"author_flair_text"`

	LinkAuthor string `mapstructure:"link_author"`
	LinkURL    string `mapstructure:"link_url"`
	LinkTitle  string `mapstructure:"link_title"`

	Subreddit   string `mapstructure:"subreddit"`
	SubredditID string `mapstructure:"subreddit_id"`

	Body     string `mapstructure:"body"`
	BodyHTML string `mapstructure:"body_html"`

	ParentID string     `mapstructure:"parent_id"`
	Replies  []*Comment `mapstructure:"reply_tree"`

	Gilded        int32  `mapstructure:"gilded"`
	Distinguished string `mapstructure:"distinguished"`
}

// IsTopLevel is true when the comment is a top level comment.
func (c *Comment) IsTopLevel() bool {
	parentType := strings.Split(c.ParentID, "_")[0]
	return parentType == postKind
}

// Media represents a subfield in the response about posts
type Media struct {
	Type   string `mapstructure:"type"`
	OEmbed struct {
		ProviderURL     string `mapstructure:"provider_url"`
		Description     string `mapstructure:"description"`
		Title           string `mapstructure:"title"`
		ThumbnailWidth  int    `mapstructure:"thumbnail_width"`
		Height          int    `mapstructure:"height"`
		Width           int    `mapstructure:"width"`
		HTML            string `mapstructure:"html"`
		Version         string `mapstructure:"version"`
		ProviderName    string `mapstructure:"provider_name"`
		ThumbnailURL    string `mapstructure:"thumbnail_url"`
		Type            string `mapstructure:"type"`
		ThumbnailHeight int    `mapstructure:"thumbnail_height"`
	} `mapstructure:"oembed"`
	RedditVideo struct {
		FallbackURL       string `mapstructure:"fallback_url"`
		Height            int    `mapstructure:"height"`
		Width             int    `mapstructure:"width"`
		ScrubberMediaURL  string `mapstructure:"scrubber_media_url"`
		DashURL           string `mapstructure:"dash_url"`
		Duration          int    `mapstructure:"duration"`
		HLSURL            string `mapstructure:"hls_url"`
		IsGIF             bool   `mapstructure:"is_gif"`
		TranscodingStatus string `mapstructure:"transcoding_status"`
	} `mapstructure:"reddit_video"`
}

// Post represents posts on Reddit (Reddit type t3_).
// https://github.com/reddit/reddit/wiki/JSON#link-implements-votable--created
type Post struct {
	ID        string `mapstructure:"id"`
	Name      string `mapstructure:"name"`
	Permalink string `mapstructure:"permalink"`

	CreatedUTC uint64 `mapstructure:"created_utc"`
	Deleted    bool   `mapstructure:"deleted"`

	Ups   int32 `mapstructure:"ups"`
	Downs int32 `mapstructure:"downs"`
	Likes bool  `mapstructure:"likes"`

	Author              string `mapstructure:"author"`
	AuthorFlairCSSClass string `mapstructure:"author_flair_css_class"`
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
	SelfTextHTML string `mapstructure:"selftext_html"`

	Replies []*Comment `mapstructure:"reply_tree"`

	Hidden            bool   `mapstructure:"hidden"`
	LinkFlairCSSClass string `mapstructure:"link_flair_css_class"`
	LinkFlairText     string `mapstructure:"link_flair_text"`

	NumComments int32  `mapstructure:"num_comments"`
	Locked      bool   `mapstructure:"locked"`
	Thumbnail   string `mapstructure:"thumbnail"`

	Gilded        int32  `mapstructure:"gilded"`
	Distinguished string `mapstructure:"distinguished"`
	Stickied      bool   `mapstructure:"stickied"`

	IsRedditMediaDomain bool  `mapstructure:"is_reddit_media_domain"`
	Media               Media `mapstructure:"media"`
	SecureMedia         Media `mapstructure:"secure_media"`
}

// Message represents messages on Reddit (Reddit type t4_).
// https://github.com/reddit/reddit/wiki/JSON#message-implements-created
type Message struct {
	ID   string `mapstructure:"id"`
	Name string `mapstructure:"name"`

	CreatedUTC uint64 `mapstructure:"created_utc"`

	Author   string `mapstructure:"author"`
	Subject  string `mapstructure:"subject"`
	Body     string `mapstructure:"body"`
	BodyHTML string `mapstructure:"body_html"`

	Context          string `mapstructure:"context"`
	FirstMessageName string `mapstructure:"first_message_name"`
	Likes            bool   `mapstructure:"likes"`
	LinkTitle        string `mapstructure:"link_title"`

	New      bool   `mapstructure:"new"`
	ParentID string `mapstructure:"parent_id"`

	Subreddit  string `mapstructure:"subreddit"`
	WasComment bool   `mapstructure:"was_comment"`
}

// Harvest is a set of all possible elements that Reddit could return in a
// listing.
type Harvest struct {
	Comments []*Comment
	Posts    []*Post
	Messages []*Message
}
