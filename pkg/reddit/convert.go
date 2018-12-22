package reddit

import red "github.com/stevenxie/graw/reddit"

// ConvPost converts a graw/reddit.Post into a Post.
func ConvPost(p *red.Post) *Post {
	return &Post{
		ID:         p.ID,
		CreatedUTC: p.CreatedUTC,
		Permalink:  p.Permalink,
		URL:        p.URL,
		Domain:     p.Domain,

		Author:      p.Author,
		Subreddit:   p.Subreddit,
		SubredditID: p.SubredditID,

		Title:    p.Title,
		IsSelf:   p.IsSelf,
		SelfText: p.SelfText,

		Deleted: p.Deleted,
		Hidden:  p.Hidden,
		NSFW:    p.NSFW,
	}
}

// ConvComment converts a graw/reddit.Comment into a Comment.
func ConvComment(c *red.Comment) *Comment {
	return &Comment{
		ID:         c.ID,
		CreatedUTC: c.CreatedUTC,
		Permalink:  c.Permalink,

		Author:      c.Author,
		Subreddit:   c.Subreddit,
		SubredditID: c.SubredditID,

		Body:     c.Body,
		ParentID: c.ParentID,

		LinkAuthor: c.LinkAuthor,
		LinkURL:    c.LinkURL,
		LinkTitle:  c.LinkTitle,

		Deleted: c.Deleted,
	}
}
