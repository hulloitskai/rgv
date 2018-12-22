package bot

import (
	"time"

	"github.com/stevenxie/graw"
	"github.com/stevenxie/graw/reddit"
	ess "github.com/unixpickle/essentials"
)

// IsActive returns true if the Bot is currently running.
func (b *Bot) IsActive() bool {
	return b.Stop != nil
}

// Run runs the Bot on the specified subreddit. It scans for new activity on
// the subreddit at an interval of `rate`.
func (b *Bot) Run(subreddit string, rate time.Duration) error {
	var (
		client reddit.Script
		err    error
	)

	// Create a bot-style client if auth information is available and valid, and
	// a script-style client with reduced polling capabilities otherwise.
	if b.creds.IsValid() {
		// Validate rate.
		if rate < time.Second {
			rate = time.Second
		}

		cfg := reddit.BotConfig{
			Agent: b.UserAgent(),
			App: reddit.App{
				ID:     b.creds.ClientID,
				Secret: b.creds.Secret,
			},
			Rate: rate,
		}
		client, err = reddit.NewBot(cfg)
	} else {
		// Validate rate.
		if rate < (2 * time.Second) {
			rate = 2 * time.Second
		}

		client, err = reddit.NewScript(b.UserAgent(), rate)
	}

	// Check for client-creation error.
	if err != nil {
		return ess.AddCtx("bot: creating Reddit client (script): %v", err)
	}

	// Configure and run client using graw, with b as the handler.
	cfg := graw.Config{
		Subreddits:        []string{subreddit},
		SubredditComments: []string{subreddit},
	}
	b.Stop, b.Wait, err = graw.Scan(b, client, cfg)
	ess.AddCtxTo("bot: performing graw scan", &err)
	return err
}
