package api

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/stevenxie/rgv/pkg/stream"
	ess "github.com/unixpickle/essentials"
)

// Server is the API server for rgv.
type Server struct {
	*http.Server
	Streamer *stream.Streamer

	l *zap.SugaredLogger
}

// NewServer returns a new Server.
func NewServer(logger *zap.SugaredLogger) (*Server, error) {
	var sl *zap.SugaredLogger
	if logger != nil {
		sl = logger.Named("Streamer")
	}

	streamer, err := stream.NewStreamer(sl)
	if err != nil {
		return nil, err
	}

	return &Server{
		Server:   &http.Server{Handler: streamer},
		Streamer: streamer,
		l:        logger,
	}, nil
}

// Shutdown gracefully shuts down the Server by shutting down its internal
// http.Server and its stream.Streamer.
func (s *Server) Shutdown(ctx context.Context) error {
	ch := make(chan error)
	go func() { ch <- s.Server.Shutdown(ctx) }()
	go func() { ch <- s.Streamer.Shutdown(ctx) }()

	var err error
	select {
	case <-ctx.Done():
		<-ch
		<-ch
		err = ctx.Err()
	case err = <-ch:
	}

	if serr := s.l.Sync(); err == nil {
		err = ess.AddCtx("api: final logger sync", serr)
	}
	return err
}
