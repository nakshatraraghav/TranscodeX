package server

import (
	"log/slog"
	"net/http"

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

	slog.Info("starting server on localhost:3000")

	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return s.server.ListenAndServe()
}
