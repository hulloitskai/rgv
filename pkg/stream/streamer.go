package stream

import (
	"bytes"
	"net/http"

	"go.uber.org/zap"

	ws "github.com/gorilla/websocket"
)

// Streamer responds to http requests by streaming activity from a particular
// subreddit.
type Streamer struct {
	buf      *bytes.Buffer
	upgrader *ws.Upgrader
	l        *zap.SugaredLogger
}

// NewStreamer creates a new Streamer that logs to logger (which may be nil).
func NewStreamer(logger *zap.SugaredLogger) *Streamer {
	return &Streamer{
		buf:      new(bytes.Buffer),
		upgrader: new(ws.Upgrader), // use default upgrader options
		l:        logger,
	}
}

func (s *Streamer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Guard against bad origin.
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed.", http.StatusForbidden)
		return
	}

	// Upgrade connection to TLS protocol.
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.errResp(w, "Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// TODO: Actually stream data into the socket.

	if err = conn.Close(); err != nil {
		s.errOrPanic("Failed to close websocket connection: %v", err)
	}
}
