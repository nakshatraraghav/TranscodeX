package schema

import (
	"time"

	"github.com/google/uuid"
)

type ApiKey struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Enabled   bool      `db:"enabled" json:"enabled"`
	Key       string    `db:"key" json:"key"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
