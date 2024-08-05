package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	addr   string
	router *chi.Mux
	server *http.Server
}

func New() *Server {

	addr := ":3000"

	router := chi.NewRouter()
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return &Server{
		addr:   addr,
		router: router,
		server: server,
	}
}

func (s *Server) Start() error {
	s.middlewares()
	s.routes()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	echan := make(chan error, 1)

	go func() {
		slog.Info("starting server on localhost:3000")
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			echan <- err
		}
		close(echan)
	}()

	select {
	case err := <-echan:
		return err
	case <-ctx.Done():
		stop()
	}

	slog.Info("shutting down the server gracefully, press Ctrl+C again to force")

	tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.server.Shutdown(tctx); err != nil && err != context.Canceled {
		return err
	}

	slog.Info("server shutdown procedure complete, graceful shutdown successful")
	return nil
}
