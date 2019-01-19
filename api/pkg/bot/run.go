package bot

import (
	"time"

	"github.com/stevenxie/graw"
	"github.com/stevenxie/graw/reddit"
	ess "github.com/unixpickle/essentials"
)

// Stop stops the Bot, if it is running.
func (b *Bot) Stop() {
	if b.stop != nil {
		b.stop()
		b.stop = nil
	}
}

// Wait blocks until the Bot has finished running, if it is running.
func (b *Bot) Wait() error {
	if b.wait != nil {
		err := b.wait()
		b.wait = nil
		return err
	}
	return nil
}

// IsActive returns true if the Bot is currently running.
func (b *Bot) IsActive() bool {
	return b.stop != nil
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
	if b.creds.IsEmpty() {
		// Validate rate.
		if rate < (2 * time.Second) {
			rate = 2 * time.Second
		}

		client, err = reddit.NewScript(b.UserAgent(), rate)
	} else {
		// Validate rate.
		if rate < time.Second {
			rate = time.Second
		}

		cfg := reddit.BotConfig{
			Agent: b.UserAgent(),
			App: reddit.App{
				ID:       b.creds.ClientID,
				Secret:   b.creds.Secret,
				Username: b.creds.User,
				Password: b.creds.Pass,
			},
			Rate: rate,
		}

		client, err = reddit.NewBot(cfg)
	}

	// Check for client-creation error.
	if err != nil {
		return ess.AddCtx("bot: creating Reddit client (script)", err)
	}

	// Configure and run client using graw, with b as the handler.
	cfg := graw.Config{
		Subreddits:        []string{subreddit},
		SubredditComments: []string{subreddit},
	}

	b.stop, b.wait, err = graw.Scan(b, client, cfg)
	return ess.AddCtx("bot: performing graw scan", err)
}
