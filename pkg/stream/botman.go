package stream

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"

	ws "github.com/gorilla/websocket"
	ess "github.com/unixpickle/essentials"
	"go.uber.org/zap"
)

const (
	// DefaultBotRate is the default Reddit-polling rate for a bot created by
	// botMan.
	DefaultBotRate = 10 * time.Second

	// DefaultPruneInterval is the default empty-bot-pruning interval.
	DefaultPruneInterval = 10 * time.Minute
)

type botManConfig struct {
	Rate          time.Duration
	PruneInterval time.Duration `split_words:"true"`
}

// botMan manages a set of Reddit bots. It is concurrent-safe.
type botMan struct {
	Config botManConfig

	bots        *botMap
	pruneTicker *time.Ticker
	pruneDone   chan empty
	l           *zap.SugaredLogger
}

// newBotMan returns a new botMan.
func newBotMan(logger *zap.SugaredLogger) (*botMan, error) {
	if logger == nil {
		logger = zap.NewNop().Sugar()
	}

	var cfg botManConfig
	if err := envconfig.Process(Namespace, &cfg); err != nil {
		return nil, err
	}

	bm := &botMan{
		Config: cfg,
		bots:   new(botMap),
		l:      logger,
	}
	go bm.pruneBots()
	return bm, nil
}

// Subscribes subscribes a websocket connection to a bot.Bot monitoring a
// particular subreddit. It is concurrent-safe.
func (bm *botMan) Subscribe(c *ws.Conn, subreddit string) error {
	// If a bot exists for the target subreddit, add c as a listener.
	if bot, ok := bm.bots.Load(subreddit); ok {
		bm.l.Debugf("Added a listener (%p) to the bot for '%s'.", c, subreddit)
		bot.AddListener(c)
		return nil
	}

	// No bot exists for the specified subreddit, so create one.
	bm.l.Infof("Creating bot for subreddit '%s'.", subreddit)

	// Ensure that the request subreddit exists before creating a bot to
	// monitor it.
	exists, err := validateSubreddit(subreddit)
	if err != nil {
		return ess.AddCtx("validating subreddit", err)
	}
	if !exists {
		bm.l.Debugf("Request subreddit '%s' does not exist. Reporting error...",
			subreddit)
		jerr := jsonError{fmt.Sprintf("subreddit '%s' does not exist", subreddit)}

		if err = c.WriteJSON(&jerr); err != nil {
			bm.l.Errorf("Error while reporting invalid subreddit: %v", err)
		}
		if err = c.Close(); err != nil {
			bm.l.Errorf("Error while closing connection: %v", err)
		}
		return nil
	}

	// Create and run bot.
	bot, err := newStreamBot(bm.l.Named(fmt.Sprintf(`streamBot["%s"]`, subreddit)))
	if err != nil {
		return ess.AddCtx("creating streamBot", err)
	}
	if err := bot.Run(subreddit, bm.Config.Rate); err != nil {
		return ess.AddCtx("running streamBot", err)
	}
	go bm.monitorBot(bot, subreddit)

	// Save bot to bm.bots, make sure that no other bot for this subreddit has
	// been concurrently added during the creation of this bot.
	if other, ok := bm.bots.Load(subreddit); ok {
		bot = other
		bm.l.Infof("Found other bot for subreddit '%s', aborting bot creation.",
			subreddit)
	} else {
		bm.bots.Store(subreddit, bot)
		bm.l.Infof("Bot for subreddit '%s' successfully created.", subreddit)
	}
	bot.AddListener(c)
	return nil
}

// monitorBot waits for sb to finish running, disconnects sb's listeners, and
// removes sb from bm.bots.
func (bm *botMan) monitorBot(sb *streamBot, subreddit string) {
	if err := sb.Wait(); err != nil {
		bm.l.Errorf("streamBot monitoring subreddit '%s' exited with an error: %v",
			subreddit, err)
	}

	if err := sb.DisconnectAll(); err != nil {
		bm.l.Errorf("Error while disconnecting listeners from streamBot for "+
			"subreddit '%s': %v", subreddit, err)
	}

	bm.l.Infof("Bot for subreddit '%s' has finished running. Removing from "+
		"botMan...", subreddit)
	bm.bots.Delete(subreddit)
}

// pruneBots checks for and removes bots with no listeners.
//
// It is a blocking function, and should be run as a goroutine.
func (bm *botMan) pruneBots() {
	bm.pruneTicker = time.NewTicker(bm.Config.PruneInterval)
	bm.pruneDone = make(chan empty)

	for {
		select {
		case <-bm.pruneTicker.C:
			bm.bots.Range(func(subreddit string, bot *streamBot) bool {
				if bot.Listeners.Len() == 0 {
					bm.l.Infof("Bot for subreddit '%s' is empty, stopping...", subreddit)
					bot.Stop()
				}
				return true
			})
		case <-bm.pruneDone:
			break
		}
	}
}

// Shutdown shuts down botMan gracefully, by stopping all bots (which
// disconnects their clients), and stopping the bot-pruning ticker.
func (bm *botMan) Shutdown(ctx context.Context) error {
	bm.pruneTicker.Stop()
	bm.pruneDone <- empty{}

	var err error
	bm.bots.Range(func(_ string, bot *streamBot) bool {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return false
		default:
			if err := bot.DisconnectAll(); err != nil {
				bm.l.Errorf("Error while disconnecting all listeners: %v", err)
			}
			return true
		}
	})
	return err
}

////////////
// botMap
///////////

type botMap struct {
	sm sync.Map
}

func (bm *botMap) Delete(key string) {
	bm.sm.Delete(key)
}

func (bm *botMap) Load(key string) (value *streamBot, ok bool) {
	var val interface{}
	val, ok = bm.sm.Load(key)
	if !ok {
		return nil, ok
	}
	return val.(*streamBot), ok
}

func (bm *botMap) LoadOrStore(key string, value *streamBot) (
	actual *streamBot, loaded bool) {
	var val interface{}
	val, loaded = bm.sm.LoadOrStore(key, value)
	return val.(*streamBot), loaded
}

func (bm *botMap) Range(f func(key string, value *streamBot) bool) {
	bm.sm.Range(func(k, v interface{}) bool {
		return f(k.(string), v.(*streamBot))
	})
}

func (bm *botMap) Store(key string, value *streamBot) {
	bm.sm.Store(key, value)
}
