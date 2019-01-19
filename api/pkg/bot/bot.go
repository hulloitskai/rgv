package bot

import (
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/stevenxie/rgv/api/internal/info"
	ess "github.com/unixpickle/essentials"
)

// A Bot is capable of receiving events from Reddit (as a logged-in account, or
// a logged-out script), using package github.com/stevenxie/graw.
type Bot struct {
	rec             Receiver
	creds           *creds // credentials for authentication
	version, uagent string
	l               *zap.SugaredLogger

	stop func()       // stops the Bot, only set if it is running
	wait func() error // block until the Bot stops, only set if it is running
}

// New returns a Bot that saves posts and comments from Reddit into the provided
// Receiver.
//
// The Bot will use credentials read from the environment; if these variables
// are not set, New will return an error (with code InvalidConfig).
func New(r Receiver, logger *zap.SugaredLogger) (*Bot, error) {
	// Validate arguments.
	if r == nil {
		return nil, errors.New("bot: cannot make Bot with a nil Receiver")
	}
	if logger == nil {
		logger = zap.NewNop().Sugar()
	}

	// Configure version string.
	version := "unset"
	if info.Version != "" {
		version = strings.TrimPrefix(info.Version, "v") // trim 'v' prefix
	}

	// Read creds from environment.
	c, err := readCreds()
	if err != nil {
		return nil, ess.AddCtx("bot", err)
	}

	// Only accept valid auth configurations.
	if err = c.Validate(); err != nil {
		return nil, ess.AddCtx("bot: validating credentials", err)
	}

	return &Bot{
		rec:     r,
		creds:   c,
		version: version,
		l:       logger,
	}, err
}

// UserAgent returns the Bot's user agent string.
func (b *Bot) UserAgent() string {
	if b.uagent == "" { // generate user agent string
		b.uagent = fmt.Sprintf("%s:%s:%s (by %s)", platform, appid, b.version,
			author)
	}
	return b.uagent
}

// Version describes the Bot's version number.
func (b *Bot) Version() string { return b.version }
