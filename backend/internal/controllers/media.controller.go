package controllers

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/nakshatraraghav/transcodex/backend/internal/schema"
	"github.com/nakshatraraghav/transcodex/backend/internal/services"
	"github.com/nakshatraraghav/transcodex/backend/types"
	"github.com/nakshatraraghav/transcodex/backend/util"
)

var kc types.ApiKeyContextKey = "apikey"

type MediaController struct {
	service services.MediaService
}

func NewMediaController(service services.MediaService) *MediaController {
	return &MediaController{
		service: service,
	}
}

func (mc *MediaController) CreateUploadHandler(w http.ResponseWriter, r *http.Request) {
	body, ok := r.Context().Value(types.ContextKey).(schema.MediaUploadRequestBody)
	if !ok {
		util.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	apiKey, ok := r.Context().Value(kc).(*schema.ApiKey)
	if !ok {
		util.WriteError(w, http.StatusBadRequest, "Invalid API key")
		return
	}

	nid, err := gonanoid.New()
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "Failed to generate S3 key")
		return
	}

	key := fmt.Sprintf("input/%v/%v/%v", apiKey.UserID.String(), nid, body.FileName)

	uploadID, presignedURL, err := mc.service.CreateUpload(r.Context(), apiKey.UserID, apiKey.ID, body.FileName, body.FileType, key)
	if err != nil {
		if err == services.ErrFailedToCreateS3SignedURL {
			util.WriteError(w, http.StatusInternalServerError, err.Error())
		} else if err == sql.ErrNoRows {
			util.WriteError(w, http.StatusInternalServerError, "Failed to create processing job, no ID returned")
		} else {
			util.WriteError(w, http.StatusInternalServerError, "Unexpected error occurred")
		}
		return
	}

	response := map[string]string{
		"upload_id":     uploadID,
		"presigned_url": presignedURL,
	}
	util.WriteJSON(w, http.StatusCreated, response)
}

func (mc *MediaController) CreateProcessingJobHandler(w http.ResponseWriter, r *http.Request) {
	apiKey, ok := r.Context().Value(kc).(*schema.ApiKey)
	if !ok {
		util.WriteError(w, http.StatusBadRequest, "Invalid API key")
		return
	}

	body, ok := r.Context().Value(types.ContextKey).(schema.CreateProcessingJobRequestBody)
	if !ok {
		util.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	uploadID, err := uuid.Parse(body.UploadID)
	if err != nil {
		slog.Error(err.Error())
		util.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	key, err := mc.service.GetS3KeyFromUpload(r.Context(), uploadID)
	if err != nil {
		if err == sql.ErrNoRows {
			util.WriteError(w, http.StatusBadRequest, "no such upload exists, invalid upload id")
		} else {
			util.WriteError(w, http.StatusInternalServerError, "unknow internal server error")
		}
		return
	}

	stringOperations, err := util.ConvertOperationsToString(body.Operations)
	if err != nil {
		slog.Error(err.Error())
		util.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	err = mc.service.AddProcessingJobToQueue(body.JobType, key, body.UploadID, stringOperations)
	if err != nil {
		slog.Error(err.Error())
		util.WriteError(w, http.StatusInternalServerError, "failed to write job to queue")
		return
	}

	id, err := mc.service.CreateProcessingJob(r.Context(), apiKey.UserID, uploadID, body.JobType, apiKey.ID)
	if err != nil {
		if err == services.ErrSQSFailedToPushMessage {
			util.WriteError(w, http.StatusInternalServerError, "failed to push operation into queue")
			return
		} else if err == sql.ErrNoRows {
			util.WriteError(w, http.StatusInternalServerError, "Failed to create processing job, no ID returned")
		} else {
			util.WriteError(w, http.StatusInternalServerError, "Unexpected error occurred")
		}
		return
	}

	util.WriteJSON(w, http.StatusCreated, map[string]string{
		"media_id": id,
	})
}
