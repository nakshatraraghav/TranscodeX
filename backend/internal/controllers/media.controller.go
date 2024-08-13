package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

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
			util.WriteError(w, http.StatusBadRequest, "No matching records found")
		} else {
			util.WriteError(w, http.StatusInternalServerError, "Unexpected error occurred")
		}
		return
	}

	// Respond with the upload ID and the presigned URL
	response := map[string]string{
		"upload_id":     uploadID,
		"presigned_url": presignedURL,
	}
	util.WriteJSON(w, http.StatusCreated, response)
}
