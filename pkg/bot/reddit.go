package bot

import (
	red "github.com/stevenxie/graw/reddit"
	"github.com/stevenxie/rgv/pkg/reddit"
)

// Receiver receives posts and comments from Reddit.
type Receiver interface {
	ReceivePost(p *reddit.Post) error
	ReceiveComment(p *reddit.Comment) error
}

// Post implements graw/botfaces.PostHandler.
func (b *Bot) Post(p *red.Post) error {
	return b.rec.ReceivePost(reddit.ConvPost(p))
}

// Comment implements graw/botfaces.CommentHandler.
func (b *Bot) Comment(c *red.Comment) error {
	return b.rec.ReceiveComment(reddit.ConvComment(c))
}
