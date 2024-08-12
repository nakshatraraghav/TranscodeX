package s3

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	cfg "github.com/nakshatraraghav/transcodex/worker/config"
)

type S3Service struct {
	client *s3.Client
}

func NewS3Service() (*S3Service, error) {
	service := &S3Service{}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return service, err
	}

	service.client = s3.NewFromConfig(cfg)

	return service, nil
}

func (s *S3Service) Download(ctx context.Context) error {
	bucketName := cfg.GetEnv().BUCKET_NAME
	objectKey := cfg.GetEnv().OBJECT_KEY

	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		return err
	}

	defer result.Body.Close()

	fpath := filepath.Join("assets", filepath.Base(objectKey))

	file, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer file.Close()

	buf, err := io.ReadAll(result.Body)
	if err != nil {
		return err
	}

	_, err = file.Write(buf)
	if err != nil {
		return err
	}

	return nil

}

// s3.go

func (s *S3Service) Upload(ctx context.Context, filePath string) error {
	bucketName := cfg.GetEnv().BUCKET_NAME
	objectKey := GetNewObjectKey(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		return err
	}

	return nil
}

func GetNewObjectKey(filePath string) string {
	objectKey := strings.Replace(cfg.GetEnv().OBJECT_KEY, "input", "output", 1)

	base := filepath.Base(cfg.GetEnv().OBJECT_KEY)
	newBase := filepath.Base(filePath)
	objectKey = strings.Replace(objectKey, base, newBase, 1)

	return objectKey

}
