package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/nakshatraraghav/transcodex/backend/internal/schema"
)

var (
	ErrApiKeyNotFound = errors.New("api key not found")
)

// in FINDAPIKEY we wlready have the key

type ApiKeyService interface {
	CreateApiKey(context.Context, string, uuid.UUID) (*schema.ApiKey, error)
	FindApiKey(context.Context, string) (*schema.ApiKey, error)
	FindValidApiKey(context.Context, uuid.UUID) (*schema.ApiKey, error)
	RevokeApiKey(context.Context, string) error
}

type apiKeyService struct {
	db *sql.DB
}

func NewApiKeyService(db *sql.DB) ApiKeyService {
	return &apiKeyService{
		db: db,
	}
}

func (aks *apiKeyService) CreateApiKey(ctx context.Context, key string, id uuid.UUID) (*schema.ApiKey, error) {
	var apikey schema.ApiKey

	q := `INSERT INTO api_keys (key, user_id)
	VALUES ($1, $2)
	RETURNING id, enabled, key, user_id, created_at, updated_at`

	err := aks.db.QueryRowContext(ctx, q, key, id).Scan(
		&apikey.ID, &apikey.Enabled, &apikey.Key, &apikey.UserID, &apikey.CreatedAt, &apikey.UpdatedAt)

	if err != nil {
		return &apikey, err
	}

	return &apikey, nil
}

func (aks *apiKeyService) FindApiKey(ctx context.Context, key string) (*schema.ApiKey, error) {
	var apikey schema.ApiKey

	q := `SELECT 
	id, enabled, key, user_id, created_at, updated_at FROM api_keys
	WHERE key = $1`

	err := aks.db.QueryRowContext(ctx, q, key).Scan(
		&apikey.ID, &apikey.Enabled, &apikey.Key, &apikey.UserID, &apikey.CreatedAt, &apikey.UpdatedAt)

	if err != nil {
		return &apikey, err
	}

	return &apikey, err
}

func (aks *apiKeyService) FindValidApiKey(ctx context.Context, id uuid.UUID) (*schema.ApiKey, error) {

	var apikey schema.ApiKey

	q := `SELECT
	id, enabled, key, user_id, created_at, updated_at
	FROM api_keys
	WHERE user_id = $1 AND enabled = true`

	err := aks.db.QueryRowContext(ctx, q, id).Scan(
		&apikey.ID, &apikey.Enabled, &apikey.Key, &apikey.UserID, &apikey.CreatedAt, &apikey.UpdatedAt)

	if err != nil {
		return &apikey, err
	}

	return &apikey, nil

}

func (aks *apiKeyService) RevokeApiKey(ctx context.Context, key string) error {

	q := `UPDATE api_keys
	SET enabled = false
	WHERE key = $1`

	res, err := aks.db.ExecContext(ctx, q, key)
	if err != nil {
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return errors.New("failed to get the affected rows")
	}

	if cnt == 0 {
		return ErrApiKeyNotFound
	}

	return nil
}
