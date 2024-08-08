package routes

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nakshatraraghav/transcodex/backend/internal/middlewares"
)

func ApiKeyRouter(router *chi.Mux, db *sql.DB) {

	subrouter := chi.NewRouter()
	router.Mount("/apikeys", subrouter)

	subrouter.With(
		middlewares.AuthMiddleware,
	).Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("protected api key route"))
	})
}
