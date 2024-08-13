package routes

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	s3service "github.com/nakshatraraghav/transcodex/backend/internal/aws/s3"
	"github.com/nakshatraraghav/transcodex/backend/internal/controllers"
	"github.com/nakshatraraghav/transcodex/backend/internal/middlewares"
	"github.com/nakshatraraghav/transcodex/backend/internal/schema"
	"github.com/nakshatraraghav/transcodex/backend/internal/services"
)

func MediaRouter(router *chi.Mux, db *sql.DB) {

	subrouter := chi.NewRouter()
	router.Mount("/media", subrouter)

	apikeyService := services.NewApiKeyService(db)

	s3Service, err := s3service.NewS3Service()
	if err != nil {
		slog.Error(err.Error())
	}

	mediaService := services.NewMediaService(s3Service, db)

	controller := controllers.NewMediaController(mediaService)

	subrouter.With(middlewares.ValidateApiKey(apikeyService)).
		Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("REACHED MEDIA ROUTE"))
		})

	subrouter.With(
		middlewares.ValidateRequestBody[schema.MediaUploadRequestBody],
		middlewares.ValidateApiKey(apikeyService),
	).Post("/upload", controller.CreateUploadHandler)
}
