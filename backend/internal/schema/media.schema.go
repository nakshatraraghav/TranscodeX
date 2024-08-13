package schema

import (
	"time"

	"github.com/google/uuid"
)

type Upload struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	FileName  string    `db:"file_name" json:"file_name"`
	FileType  string    `db:"file_type" json:"file_type"`
	ApiKeyID  uuid.UUID `db:"apikey_id" json:"apikey_id"`
	S3Url     string    `db:"s3_url" json:"s3_url"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type ProcessingJob struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	UploadID  uuid.UUID `db:"upload_id" json:"upload_id"`
	ApiKeyID  uuid.UUID `db:"apikey_id" json:"apikey_id"`
	JobType   string    `db:"job_type" json:"job_type"`
	Status    string    `db:"status" json:"status"`
	ResultUrl string    `db:"result_url" json:"result_url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type MediaUploadRequestBody struct {
	FileName string `validate:"required" json:"file_name"`
	FileType string `validate:"required" json:"file_type"`
	MimeType string `validate:"required" json:"mime_type"`
}
