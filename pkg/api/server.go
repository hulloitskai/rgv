package api

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/stevenxie/rgv/pkg/stream"
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
	}, nil
}

// Shutdown gracefully shuts down the Server by shutting down its internal
// http.Server and its stream.Streamer.
func (s *Server) Shutdown(ctx context.Context) error {
	ch := make(chan error)
	go func() { ch <- s.Server.Shutdown(ctx) }()
	go func() { ch <- s.Streamer.Shutdown(ctx) }()

	select {
	case <-ctx.Done():
		<-ch
		<-ch
		return ctx.Err()
	case err := <-ch:
		return err
	}
}
