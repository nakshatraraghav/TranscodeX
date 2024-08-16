package routes

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	s3service "github.com/nakshatraraghav/transcodex/backend/internal/aws/s3"
	sqsservice "github.com/nakshatraraghav/transcodex/backend/internal/aws/sqs"
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
		panic(err)
	}

	sqsService, err := sqsservice.NewSQSService()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	mediaService := services.NewMediaService(s3Service, sqsService, db)

	controller := controllers.NewMediaController(mediaService)

	subrouter.With(
		middlewares.ValidateApiKey(apikeyService),
	).Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	subrouter.With(
		middlewares.ValidateApiKey(apikeyService),
		middlewares.ValidateRequestBody[schema.MediaUploadRequestBody],
	).Post("/upload", controller.CreateUploadHandler)

	subrouter.With(
		middlewares.ValidateApiKey(apikeyService),
		middlewares.ValidateRequestBody[schema.CreateProcessingJobRequestBody],
	).Post("/process", controller.CreateProcessingJobHandler)

	subrouter.With(
		middlewares.ValidateApiKey(apikeyService),
	).Get("/status/{job_id}", controller.GetProcessingJobStatus)

	subrouter.With(
		middlewares.ValidateApiKey(apikeyService),
	).Get("/download/{job_id}", controller.DownloadProcessedMediaHandler)
}
