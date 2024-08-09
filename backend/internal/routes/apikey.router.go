package routes

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/nakshatraraghav/transcodex/backend/internal/controllers"
	"github.com/nakshatraraghav/transcodex/backend/internal/middlewares"
	"github.com/nakshatraraghav/transcodex/backend/internal/services"
)

func ApiKeyRouter(router *chi.Mux, db *sql.DB) {

	subrouter := chi.NewRouter()
	router.Mount("/apikeys", subrouter)

	service := services.NewApiKeyService(db)
	controller := controllers.NewApiKeyController(service)

	subrouter.With(
		middlewares.AuthMiddleware,
	).Post("/", controller.CreateApiKeyHandler)

	subrouter.With(
		middlewares.AuthMiddleware,
	).Get("/", controller.GetActiveApiKeyHandler)

	subrouter.With(
		middlewares.AuthMiddleware,
		middlewares.EnsureApiKeyInRequestHeaders,
	).Delete("/", controller.RevokeApiKeyController)
}
