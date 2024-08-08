package routes

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nakshatraraghav/transcodex/backend/internal/controllers"
	"github.com/nakshatraraghav/transcodex/backend/internal/middlewares"
	"github.com/nakshatraraghav/transcodex/backend/internal/schema"
	"github.com/nakshatraraghav/transcodex/backend/internal/services"
)

func SessionRouter(router *chi.Mux, db *sql.DB) {

	subrouter := chi.NewRouter()
	router.Mount("/sessions", subrouter)

	subrouter.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	us := services.NewUserService(db)
	ss := services.NewSessionService(db)

	controller := controllers.NewSessionController(us, ss)

	subrouter.With(
		middlewares.ValidateRequestBody[schema.CreateSessionSchema],
	).Post("/", controller.CreateSessionHandler)

	subrouter.With(
		middlewares.AuthMiddleware,
	).Get("/", controller.GetCurrentSessionHandler)
}
