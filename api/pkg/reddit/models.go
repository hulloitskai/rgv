package reddit

// Post represents a Reddit post, based off of graw/reddit.Post.
//
// It has been simplified to remove extraneous fields, and has struct tags to
// make it JSON marshallable.
type Post struct {
	ID         string `json:"id"`
	CreatedUTC uint64 `json:"created_utc"`
	Permalink  string `json:"permalink"`
	URL        string `json:"url"`
	Domain     string `json:"domain"`

	Author      string `json:"author"`
	Subreddit   string `json:"subreddit"`
	SubredditID string `json:"subreddit_id"`

	Title    string `json:"title"`
	IsSelf   bool   `json:"is_self"`
	SelfText string `json:"selftext"`

	Deleted bool `json:"deleted"`
	Hidden  bool `json:"hidden"`
	NSFW    bool `json:"nsfw"`
}

// Comment represents a Reddit comment, based off of graw/reddit.Comment.
//
// It has been simplified to remove extraneous fields, and has struct tags to
// make it JSON marshallable.
type Comment struct {
	ID         string `json:"id"`
	CreatedUTC uint64 `json:"created_utc"`
	Permalink  string `json:"permalink"`

	Author      string `json:"author"`
	Subreddit   string ` json:"subreddit"`
	SubredditID string `json:"subreddit_id"`

	Body     string `json:"body"`
	ParentID string `json:"parent_id"`

	LinkAuthor string `json:"link_author"`
	LinkURL    string `json:"link_url"`
	LinkTitle  string `json:"linke_title"`

	Deleted bool `json:"deleted"`
}
