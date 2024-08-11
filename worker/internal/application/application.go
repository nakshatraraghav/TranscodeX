package application

import (
	"context"

	"github.com/nakshatraraghav/transcodex/worker/config"
	"github.com/nakshatraraghav/transcodex/worker/internal/processors/image"
	"github.com/nakshatraraghav/transcodex/worker/internal/processors/video"
	"github.com/nakshatraraghav/transcodex/worker/internal/s3"
)

type MediaProcessor interface {
	LoadData() error
	ApplyTransformations(map[string]string) []error
}

type Application struct {
	processor *MediaProcessor
	service   *s3.S3Service
}

func NewApp() (*Application, error) {

	err := config.LoadEnv()
	if err != nil {
		return nil, err
	}

	env := config.GetEnv()

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
		service:   s,
		processor: &p,
	}, nil
}

func (a *Application) Run() error {

	// Step 1: Download the media file from S3
	err := a.service.Download(context.Background())
	if err != nil {
		return err
	}

	// Step 2: run transformations on it

	// Step 3: upload the s3
	err = a.service.Upload(context.Background())
	if err != nil {
		return err
	}

	return nil
}
