package services

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	s3service "github.com/nakshatraraghav/transcodex/backend/internal/aws/s3"
)

var (
	ErrFailedToCreateS3SignedURL = errors.New("failed to create a presigned url")
	ErrDatabaseNotInitialized    = errors.New("database connection is not initialized")
	ErrS3ServiceNotInitialized   = errors.New("S3 service is not initialized")
)

type MediaService interface {
	CreateUpload(ctx context.Context,
		UserID uuid.UUID,
		ApiKeyID uuid.UUID,
		FileName, FileType string,
		S3Key string) (string, string, error)
}

type mediaService struct {
	s3 *s3service.S3Service
	db *sql.DB
}

func NewMediaService(s3 *s3service.S3Service, db *sql.DB) MediaService {
	return &mediaService{
		s3: s3,
		db: db,
	}
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
