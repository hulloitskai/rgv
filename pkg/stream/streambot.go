package stream

import (
	"fmt"
	"strings"
	"sync"
	"syscall"

	"go.uber.org/zap"

	ws "github.com/gorilla/websocket"
	"github.com/stevenxie/rgv/pkg/bot"
	"github.com/stevenxie/rgv/pkg/reddit"
	ess "github.com/unixpickle/essentials"
)

// A streamBot is a Bot which streams data from Reddit to a set of websocket
// listeners.
//
// It is concurrent-safe: listeners may be added / removed concurrently.
type streamBot struct {
	*bot.Bot
	Listeners *socketSet // concurrent-safe
}

// newStreamBot returns a new streamBot.
func newStreamBot(logger *zap.SugaredLogger) (*streamBot, error) {
	sb := &streamBot{Listeners: new(socketSet)}
	var err error
	sb.Bot, err = bot.New(sb, logger)
	return sb, err
}

// ReceivePost implements bot.Receiver for streamBot.
func (sb *streamBot) ReceivePost(p *reddit.Post) error {
	return sb.BroadcastJSON(p, "Reddit post")
}

// ReceiveComment implements bot.Receiver for streamBot.
func (sb *streamBot) ReceiveComment(c *reddit.Comment) error {
	return sb.BroadcastJSON(c, "Reddit comment")
}

// BroadcastJSON broadcasts an object as JSON to all of streamBot's listeners.
//
// If an error occurs, desc is used to describe the broadcasted object in the
// error description.
func (sb *streamBot) BroadcastJSON(v interface{}, desc string) error {
	var err error
	sb.Listeners.Range(func(sc *syncConn) bool {
		if cerr := sc.WriteJSON(v); cerr != nil {
			if ws.IsCloseError(cerr) ||
				strings.Contains(cerr.Error(), syscall.EPIPE.Error()) {
				sb.Listeners.Delete(sc)
				return true
			}

			err = ess.AddCtx(fmt.Sprintf("stream: writing %s to conn", desc), cerr)
			return false
		}
		return true
	})
	return err
}

// DisconnectAll disconnects all listeners from the streamBot.
func (sb *streamBot) DisconnectAll() error {
	var err error
	sb.Listeners.Range(func(sc *syncConn) bool {
		if err == nil {
			err = sc.C.Close()
		}
		sb.Listeners.Delete(sc)
		return true
	})
	return err
}

// AddListener adds a websocket connection to streamBot.
func (sb *streamBot) AddListener(c *ws.Conn) {
	sb.Listeners.Store(&syncConn{C: c})
}

//////////////
// syncConn
//////////////

// A syncConn wraps a websocket connection with special methods that synchronize
// concurrent writes to the websocket.
type syncConn struct {
	C   *ws.Conn
	mux sync.Mutex
}

// WriteJSON calls sc.C.WriteJSON. It is concurrent-safe.
func (sc *syncConn) WriteJSON(v interface{}) error {
	sc.mux.Lock()
	defer sc.mux.Unlock()
	return sc.C.WriteJSON(v)
}

///////////////
// socketSet
///////////////

type empty struct{}

// A socketSet is a concurrent-safe set of websockets.
type socketSet struct {
	sm sync.Map
}

func (ss *socketSet) Delete(sc *syncConn) {
	ss.sm.Delete(sc)
}

func (ss *socketSet) Range(f func(*syncConn) bool) {
	ss.sm.Range(func(key, _ interface{}) bool {
		return f(key.(*syncConn))
	})
}

func (ss *socketSet) Store(sc *syncConn) {
	ss.sm.Store(sc, empty{})
}

// Len returns the number of elements in the set. It has an O(n) time
// complexity (relatively inefficient).
func (ss *socketSet) Len() int {
	var count int
	ss.sm.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}
