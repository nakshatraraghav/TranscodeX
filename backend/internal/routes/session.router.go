package routes

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SessionRouter(router *chi.Mux, db *sql.DB) {

	subrouter := chi.NewRouter()
	router.Mount("/sessions", subrouter)

	subrouter.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}
