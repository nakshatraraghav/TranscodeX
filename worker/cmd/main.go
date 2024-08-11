package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/nakshatraraghav/transcodex/worker/config"
	"github.com/nakshatraraghav/transcodex/worker/internal/s3"
)

func main() {
	if err := config.LoadEnv(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	service, err := s3.NewS3Service()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	err = service.Upload(context.Background())
	if err != nil {
		slog.Error(err.Error())
	}

}
