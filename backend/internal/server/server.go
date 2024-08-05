package server

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nakshatraraghav/transcodex/backend/db"
	"github.com/nakshatraraghav/transcodex/backend/internal/config"
	"github.com/nakshatraraghav/transcodex/backend/util"
)

type Server struct {
	addr   string
	db     *sql.DB
	router *chi.Mux
	server *http.Server
}

func init() {
	if err := config.LoadEnv(); err != nil {
		panic(err)
	}
	art := util.AsciiArt()
	slog.Info(art)
}

func New() *Server {

	addr := config.GetEnv().Addr

	router := chi.NewRouter()
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	db, err := db.NewPostgresConnection()
	if err != nil {
		panic(err)
	}

	return &Server{
		addr:   addr,
		db:     db,
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
		slog.Info("starting server on localhost" + s.addr)
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

	slog.Info("attempting to close database connection")
	err := s.db.Close()
	if err != nil {
		return err
	}

	slog.Info("database connection successfully closed")

	slog.Info("server shutdown procedure complete, graceful shutdown successful")
	return nil
}
