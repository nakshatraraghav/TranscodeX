package controllers

import (
	"database/sql"
	"net/http"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/nakshatraraghav/transcodex/backend/internal/services"
	"github.com/nakshatraraghav/transcodex/backend/util"
)

type ApiKeyController struct {
	service services.ApiKeyService
}

func NewApiKeyController(service services.ApiKeyService) ApiKeyController {
	return ApiKeyController{
		service: service,
	}
}

func (akc *ApiKeyController) CreateApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := util.NewJwtClaims(r)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid cookies")
		return
	}

	// Step 1: Check for a valid existing api key
	_, err = akc.service.FindValidApiKey(r.Context(), claims.ID)
	if err == nil {
		util.WriteError(w, http.StatusConflict, "user already has a valid API key")
		return
	}

	// Step 2: Generate a api key for this user
	key, err := gonanoid.New()
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to generate a apikey for this user")
		return
	}

	// Step 3: Store the api key in the database
	apikey, err := akc.service.CreateApiKey(r.Context(), key, claims.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			util.WriteError(w, http.StatusInternalServerError, "failed to store the api key, please try again")
		} else {
			util.WriteError(w, http.StatusInternalServerError, "unexpected error, while storing the api  key")
		}

		return
	}

	r.Header.Set("x-api-key", apikey.Key)

	util.WriteJSON(w, http.StatusOK, apikey)

}

func (akc *ApiKeyController) GetActiveApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := util.NewJwtClaims(r)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "invalid tokens")
		return
	}

	key, err := akc.service.FindValidApiKey(r.Context(), claims.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			util.WriteError(w, http.StatusNotFound, "no active API key found")
		} else {
			util.WriteError(w, http.StatusInternalServerError, "failed to retrieve API key")
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, key)
}

func (akc *ApiKeyController) RevokeApiKeyController(w http.ResponseWriter, r *http.Request) {

	claims, err := util.NewJwtClaims(r)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid tokens")
		return
	}

	key := r.Header.Get("x-api-key")

	// Step 1: check if this api key belongs to the logged in user
	valid, err := akc.service.FindValidApiKey(r.Context(), claims.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			util.WriteError(w, http.StatusNotFound, "no active API key found")
		} else {
			util.WriteError(w, http.StatusInternalServerError, "failed to retrieve API key")
		}
		return
	}

	if valid.Key != key {
		util.WriteError(w, http.StatusUnauthorized, "unauthorized, you cannot delete this api key")
		return
	}

	err = akc.service.RevokeApiKey(r.Context(), valid.Key)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to revoke the api key")
		return
	}

	util.WriteJSON(w, http.StatusOK, "successfully revoked the api key")

}
