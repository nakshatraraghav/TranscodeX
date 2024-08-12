package routes

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nakshatraraghav/transcodex/backend/internal/middlewares"
	"github.com/nakshatraraghav/transcodex/backend/internal/services"
)

func MediaRouter(router *chi.Mux, db *sql.DB) {

	subrouter := chi.NewRouter()
	router.Mount("/media", subrouter)

	apikeyService := services.NewApiKeyService(db)

	subrouter.With(middlewares.ValidateApiKey(apikeyService)).
		Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("REACHED MEDIA ROUTE"))
		})
}
