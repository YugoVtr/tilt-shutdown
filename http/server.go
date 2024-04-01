package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	DefaultTerminationGracePeriod = 30 * time.Second
)

// Server is something that listens for connections and do HTTP things.
type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// ServerOption provides ways of configuring Server.
type ServerOption func(server *GracefulServer)

// GracefulServer wraps around a Server and will gracefully shutdown when it receives a shutdown signal.
type GracefulServer struct {
	delegate Server
	timeout  time.Duration
	shutdown <-chan os.Signal
	logger   *slog.Logger
}

// WithLogger allows logs in the server.
func WithLogger(logger *slog.Logger) ServerOption {
	return func(server *GracefulServer) {
		server.logger = logger
	}
}

// NewServer returns a Server that can gracefully shutdown on shutdown signals.
func NewServer(server Server, options ...ServerOption) *GracefulServer {
	gracefulServer := &GracefulServer{
		delegate: server,
		timeout:  DefaultTerminationGracePeriod,
		shutdown: newInterruptSignalChannel(),
	}
	for _, option := range options {
		option(gracefulServer)
	}
	return gracefulServer
}

// ListenAndServe will call the ListenAndServe function of the delegate Server.
// On a signal being sent to the shutdown signal provided in the constructor, it will call the server's Shutdown method to attempt to gracefully shutdown.
func (s *GracefulServer) ListenAndServe(ctx context.Context) error {
	select {
	case err := <-s.delegateListenAndServe():
		return err
	case <-ctx.Done():
		return s.shutdownDelegate(ctx)
	case <-s.shutdown:
		return s.shutdownDelegate(ctx)
	}
}

func (s *GracefulServer) delegateListenAndServe() chan error {
	s.log("server running...")
	listenErr := make(chan error)
	go func() {
		if err := s.delegate.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			listenErr <- err
		}
	}()
	return listenErr
}

func (s *GracefulServer) shutdownDelegate(ctx context.Context) error {
	s.log("shutting down server...")
	defer s.log("shutting down server...(done)")
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	if err := s.delegate.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server shutdown failed: %w", err)
	}
	return ctx.Err()
}

func (s *GracefulServer) log(msg string) {
	if s.logger != nil {
		s.logger.Info(msg)
	}
}
