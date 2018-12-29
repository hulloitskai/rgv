package stream

import (
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"

	ws "github.com/gorilla/websocket"
)

// Streamer responds to http requests by streaming activity from a particular
// subreddit.
type Streamer struct {
	*botMan

	upgrader *ws.Upgrader
	l        *zap.SugaredLogger
}

// NewStreamer creates a new Streamer that logs to logger (which may be nil).
func NewStreamer(logger *zap.SugaredLogger) (*Streamer, error) {
	// Derive logger for botMan.
	var bml *zap.SugaredLogger
	if logger != nil {
		bml = logger.Named("botMan")
	}

	bm, err := newBotMan(bml)
	if err != nil {
		return nil, err
	}

	// Create upgrader based on GO_ENV.
	upgrader := new(ws.Upgrader)
	if os.Getenv("GO_ENV") == "development" {
		upgrader.CheckOrigin = func(*http.Request) bool { return true }
	}

	return &Streamer{
		botMan:   bm,
		upgrader: upgrader,
		l:        logger,
	}, nil
}

func (s *Streamer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection to TLS protocol.
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.errResp(w, "Failed to upgrade connection: %v", err)
		return
	}

	// Log request origin.
	if s.l != nil {
		s.l.Debugf("Got a connection from: %s", r.Header.Get("Origin"))
	}

	var confmsg struct {
		Subreddit string `json:"subreddit"`
	}
	var jsonErr struct {
		Error string `json:"error"`
	}
	if err = conn.ReadJSON(&confmsg); err != nil {
		if s.l != nil {
			s.l.Errorf("Error reading initial config message as JSON: %v", err)
		}
		jsonErr.Error = "bad initial config message: " + err.Error()

		// Send error to conn.
		if err = conn.WriteJSON(&jsonErr); (s.l != nil) && (err != nil) {
			s.l.Errorf("Error while reporting config error: %v", err)
		}

		if err = conn.Close(); (s.l != nil) && (err != nil) {
			s.l.Errorf("Error closing connection: %v", err)
		}
		return
	}

	if err = s.botMan.Subscribe(conn, confmsg.Subreddit); err != nil {
		if s.l != nil {
			s.l.Errorf("Error while subscribing client to subreddit '%s': %v",
				confmsg.Subreddit, err)
		}
		jsonErr.Error = fmt.Sprintf("failed to subscribe client: %v", err)

		if err = conn.WriteJSON(&jsonErr); (s.l != nil) && (err != nil) {
			s.l.Errorf("Error while reporting subscription error: %v", err)
		}

		if err = conn.Close(); (s.l != nil) && (err != nil) {
			s.l.Errorf("Error closing connection: %v", err)
		}
		return
	}
}

// errResp writes an error to both w and Streamer's internal logger.
//
// It sets the response status code (for w) to 500 (Internal Server Error).
func (s *Streamer) errResp(w http.ResponseWriter, format string,
	a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	http.Error(w, msg, http.StatusInternalServerError)
	if s.l != nil {
		s.l.Error(msg)
	}
}
