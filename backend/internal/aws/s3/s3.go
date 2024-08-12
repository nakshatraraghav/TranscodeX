package s3service

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	cg "github.com/nakshatraraghav/transcodex/backend/internal/config"
)

type S3Service struct {
	client     *s3.PresignClient
	bucketName string
}

func NewS3Service() (*S3Service, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	c := s3.NewFromConfig(cfg)
	client := s3.NewPresignClient(c)

	return &S3Service{
		client:     client,
		bucketName: cg.GetEnv().BUCKET_NAME,
	}, nil
}

func (s *S3Service) GeneratePresignedDownloadURL(key string) (string, error) {
	url, err := s.client.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	},
		s3.WithPresignExpires(10*time.Minute),
	)

	if err != nil {
		return "", err
	}

	return url.URL, nil
}

func (s *S3Service) GeneratePresignedUploadURL(key string) (string, error) {
	url, err := s.client.PresignPutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return "", err
	}

	return url.URL, nil
}
