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

func UserRouter(router *chi.Mux, db *sql.DB) {

	subrouter := chi.NewRouter()
	router.Mount("/users", subrouter)

	service := services.NewUserService(db)
	controller := controllers.NewUserController(service)

	subrouter.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	subrouter.With(middlewares.ValidateRequestBody[schema.CreateUserSchema]).Post("/", controller.CreateUserHandler)
}
