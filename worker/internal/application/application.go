package application

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/nakshatraraghav/transcodex/worker/config"
	"github.com/nakshatraraghav/transcodex/worker/db"
	"github.com/nakshatraraghav/transcodex/worker/internal/processors/image"
	"github.com/nakshatraraghav/transcodex/worker/internal/processors/video"
	"github.com/nakshatraraghav/transcodex/worker/internal/s3"
	"github.com/nakshatraraghav/transcodex/worker/internal/service"
)

type MediaProcessor interface {
	ApplyTransformations(map[string]string) []error
}

type Application struct {
	db        *sql.DB
	processor MediaProcessor
	service   *s3.S3Service
}

func NewApp() (*Application, error) {

	err := config.LoadEnv()
	if err != nil {
		return nil, err
	}

	env := config.GetEnv()

	db, err := db.NewPostgresConnection()
	if err != nil {
		return nil, err
	}

	s, err := s3.NewS3Service()
	if err != nil {
		return nil, err
	}

	var p MediaProcessor

	switch env.MEDIA_TYPE {
	case "IMAGE":
		slog.Info("IMAGE Processing Selected")
		p = image.NewImageProcessor()
	case "VIDEO":
		slog.Info("VIDEO Processing Selected")
		p = video.NewVideoProcessor()
	default:
		p = nil
	}

	return &Application{
		db:        db,
		service:   s,
		processor: p,
	}, nil
}

func (a *Application) Run() error {
	ctx := context.Background()

	// Create output directory if it doesn't exist
	outputDir := filepath.Join("assets", "output")
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		slog.Error("Failed to create output directory", "error", err)
		return fmt.Errorf("error creating output directory: %w", err)
	}

	service := service.NewProcessingJobService(a.db)

	// Update processing job status to indicate that download is starting
	slog.Info("Starting downloads from S3")
	if err := service.ChangeProcessingJobStatus(ctx, "WORKER:DOWNLOAD_STARTING"); err != nil {
		slog.Error("Failed to update job status", "status", "WORKER:DOWNLOAD_STARTING", "error", err)
		return err
	}

	// Download files from S3
	if err := a.service.Download(ctx); err != nil {
		slog.Error("Failed to download files from S3", "error", err)
		return fmt.Errorf("error downloading media file: %w", err)
	}

	// Apply transformations
	transformations := config.GetEnv().TRANSFORMATIONS
	if len(transformations) == 0 {
		return fmt.Errorf("no transformations provided")
	}

	slog.Info("Applying transformations")
	if err := service.ChangeProcessingJobStatus(ctx, "WORKER:APPLYING_TRANSFORMATIONS"); err != nil {
		slog.Error("Failed to update job status", "status", "WORKER:APPLYING_TRANSFORMATIONS", "error", err)
		return err
	}

	errors := a.processor.ApplyTransformations(transformations)
	if len(errors) > 0 {
		for _, err := range errors {
			slog.Error("Transformation error", "error", err)
		}
		return fmt.Errorf("errors occurred during transformations")
	}
	slog.Info("Transformations complete")

	// Retrieve transformed files
	files, err := getTransformedFiles()
	if err != nil {
		slog.Error("Failed to retrieve transformed files", "error", err)
		return fmt.Errorf("error retrieving transformed files: %w", err)
	}

	// Upload transformed files to S3
	slog.Info("Starting upload of transformed files")
	if err := service.ChangeProcessingJobStatus(ctx, "WORKER:STARTING_UPLOADS"); err != nil {
		slog.Error("Failed to update job status", "status", "WORKER:STARTING_UPLOADS", "error", err)
		return err
	}

	for _, file := range files {
		url, err := a.service.Upload(ctx, file)
		if err != nil {
			slog.Error("Failed to upload file", "file", file, "error", err)
			return fmt.Errorf("error uploading file %s: %w", file, err)
		}

		err = service.AddResultURL(ctx, url)
		if err != nil {
			slog.Error("Failed to update result_url", "status", "WORKER:UPLOADS_STATUS_CHANGING_ERROR", "error", err)
			return err
		}

		slog.Info("Successfully uploaded file", "file", file)
	}

	// Update job status to indicate uploads are finished
	slog.Info("Upload finished, exiting")
	if err := service.ChangeProcessingJobStatus(ctx, "WORKER:UPLOADS_FINISHED_EXITING"); err != nil {
		slog.Error("Failed to update job status", "status", "WORKER:UPLOADS_FINISHED_EXITING", "error", err)
		return err
	}

	return nil
}

// getTransformedFiles returns a list of file paths for all transformed files.
func getTransformedFiles() ([]string, error) {
	// This implementation assumes transformed files are stored in "assets" directory.
	// Adjust the pattern as needed to match the generated files.
	files, err := filepath.Glob(filepath.Join("assets", "output", "*"))
	if err != nil {
		return nil, err
	}
	return files, nil
}
