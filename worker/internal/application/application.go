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
		p = image.NewImageProcessor()
	case "VIDEO":
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

	// Step 0: Create the assets and output folder
	op := filepath.Join("assets", "output")
	err := os.MkdirAll(op, os.ModePerm)
	if err != nil {
		fmt.Printf("error creating output directory: %v\n", err)
	}

	service := service.NewProcessingJobService(a.db)

	// Step 1: Download the media file from S3
	slog.Info("download starting")

	err = service.ChangeProcessingJobStatus(ctx, "WORKER:DOWNLOAD_STARTING")
	if err != nil {
		slog.Error(err.Error())
	}
	if err := a.service.Download(ctx); err != nil {
		return fmt.Errorf("error downloading media file: %w", err)
	}

	// Step 2: Apply transformations
	transformations := config.GetEnv().TRANSFORMATIONS
	if len(transformations) == 0 {
		return fmt.Errorf("no transformations provided")
	}

	slog.Info("starting the transformations")

	err = service.ChangeProcessingJobStatus(ctx, "WORKER:APPLYING_TRANSFORMATIONS")
	if err != nil {
		slog.Error(err.Error())
	}
	errors := a.processor.ApplyTransformations(transformations)
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Printf("Transformation error: %v\n", err)
		}
	}

	// Step 3: Find all generated files to upload
	files, err := getTransformedFiles()
	if err != nil {
		return fmt.Errorf("error retrieving transformed files: %w", err)
	}

	slog.Info("upload starting")
	err = service.ChangeProcessingJobStatus(ctx, "WORKER:STARTING_UPLOADS")
	if err != nil {
		slog.Error(err.Error())
	}

	for _, file := range files {
		if err := a.service.Upload(ctx, file); err != nil {
			return fmt.Errorf("error uploading file %s: %w", file, err)
		}
		slog.Info("succesfully uploaded :" + file)
	}

	err = service.ChangeProcessingJobStatus(ctx, "WORKER:UPLOADS_FINISHED_EXITING")
	if err != nil {
		slog.Error(err.Error())
	}

	slog.Info("upload finished, exiting")

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
