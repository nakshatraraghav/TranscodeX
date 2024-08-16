package schema

import (
	"encoding/json"
	"time"

	"github.com/go-playground/validator/v10"
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

type CreateProcessingJobRequestBody struct {
	JobType    string          `validate:"required,allowed_jobtype" json:"job_type"`
	UploadID   string          `validate:"required" json:"upload_id"`
	Operations json.RawMessage `validate:"required,allowed_operations" json:"operations"`
}

type GetProcessingJobStatusRequestBody struct {
	ProcessingJobID string `validate:"required"`
}

func ValidateJobTypeField(fl validator.FieldLevel) bool {
	jtype, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	if jtype != "VIDEO" && jtype != "IMAGE" {
		return false
	}

	return true
}

func ValidateOperationsField(fl validator.FieldLevel) bool {
	allowed := map[string]bool{
		"RESIZE":               true,
		"FORCE-RESIZE":         true,
		"ROTATE":               true,
		"CONVERT-FORMAT":       true,
		"WATERMARK":            true,
		"GENERATE-THUMBNAIL":   true,
		"TRANSCODE":            true,
		"TRANSCODE-RESOLUTION": true,
	}

	operations, ok := fl.Field().Interface().(json.RawMessage)
	if !ok {
		return false
	}

	var ops map[string]string
	err := json.Unmarshal(operations, &ops)
	if err != nil {
		return false
	}

	for key, _ := range ops {
		if _, exists := allowed[key]; !exists {
			return false
		}
	}

	return true

}
