package http_server

import (
	"context"
	"net/http"
	"time"
)

const (
	defaultReadTimeout     = 10 * time.Second
	defaultWriteTimeout    = 10 * time.Second
	defalutIdelTimeout     = 30 * time.Second
	defaultShutdownTimeout = 3 * time.Second
	defaultAddress         = "localhost:8080"
)

// Server - http server
type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

// New - create new server
func New(handler http.Handler, options ...Option) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defalutIdelTimeout,
		Addr:         defaultAddress,
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range options {
		opt(s)
	}

	s.start()

	return s
}

func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

// Notify - notify
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown - shutdown
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
