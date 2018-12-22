package stream

import (
	"fmt"
	"net/http"
)

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

// errOrPanic writes an error log to the Streamer's internal logger if it is
// non-nil, or else it panics with the error.
func (s *Streamer) errOrPanic(format string, a ...interface{}) {
	if s.l == nil {
		panic(fmt.Sprintf(format, a...))
	}
	s.l.Errorf(format, a...)
}

// err writes an error log to the Streamer's internal logger, if it is non-nil.
func (s *Streamer) err(format string, a ...interface{}) {
	if s.l != nil {
		s.l.Errorf(format, a...)
	}
}
