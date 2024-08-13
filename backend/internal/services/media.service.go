package services

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	s3service "github.com/nakshatraraghav/transcodex/backend/internal/aws/s3"
	sqsservice "github.com/nakshatraraghav/transcodex/backend/internal/aws/sqs"
)

var (
	ErrFailedToCreateS3SignedURL = errors.New("failed to create a presigned url")
	ErrDatabaseNotInitialized    = errors.New("database connection is not initialized")
	ErrS3ServiceNotInitialized   = errors.New("S3 service is not initialized")
	ErrSQSFailedToPushMessage    = errors.New("failed to push message into queue")
)

type MediaService interface {
	GetS3KeyFromUpload(ctx context.Context, UploadID uuid.UUID) (string, error)

	AddProcessingJobToQueue(mediaType, key, transformations string) error

	CreateUpload(ctx context.Context,
		UserID uuid.UUID,
		ApiKeyID uuid.UUID,
		FileName, FileType string,
		S3Key string) (string, string, error)

	CreateProcessingJob(ctx context.Context,
		UserID uuid.UUID,
		UploadID uuid.UUID,
		JobType string,
		ApiKeyID uuid.UUID) (string, error)
}

type mediaService struct {
	s3  *s3service.S3Service
	sqs *sqsservice.SQSService
	db  *sql.DB
}

func NewMediaService(s3 *s3service.S3Service,
	sqs *sqsservice.SQSService,
	db *sql.DB) MediaService {
	return &mediaService{
		s3:  s3,
		sqs: sqs,
		db:  db,
	}
}

func (ms *mediaService) GetS3KeyFromUpload(ctx context.Context, UploadID uuid.UUID) (string, error) {
	var key string
	q := "SELECT s3_url FROM uploads WHERE id = $1"

	row := ms.db.QueryRowContext(ctx, q, UploadID)
	err := row.Scan(&key)
	if err != nil {
		return "", err
	}

	return key, nil
}

func (ms *mediaService) CreateUpload(ctx context.Context,
	UserID uuid.UUID,
	ApiKeyID uuid.UUID,
	FileName, FileType string,
	S3Key string) (string, string, error) {

	url, err := ms.s3.GeneratePresignedUploadURL(S3Key)
	if err != nil {
		slog.Error("Failed to generate presigned URL", "error", err.Error())
		return "", "", ErrFailedToCreateS3SignedURL
	}
	var uploadID string
	q := `INSERT INTO uploads 
	(user_id, file_name, file_type, s3_url, status, apikey_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id`

	row := ms.db.QueryRowContext(ctx, q, UserID, FileName, FileType, S3Key, "URLGENERATED", ApiKeyID)
	err = row.Scan(&uploadID)
	if err != nil {
		slog.Error("Failed to insert upload record:", "error", err.Error())
		return "", "", err
	}

	return uploadID, url, nil
}

func (ms *mediaService) CreateProcessingJob(
	ctx context.Context,
	UserID uuid.UUID,
	UploadID uuid.UUID,
	JobType string,
	ApiKeyID uuid.UUID) (string, error) {
	var id string

	q := `INSERT INTO processing_jobs
	(user_id, upload_id, job_type, apikey_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	row := ms.db.QueryRowContext(ctx, q, UserID, UploadID, JobType, ApiKeyID)
	err := row.Scan(&id)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (ms *mediaService) AddProcessingJobToQueue(mediaType, key, transformations string) error {
	return ms.sqs.Enqueue(mediaType, key, transformations)
}
