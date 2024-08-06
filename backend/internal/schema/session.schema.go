package schema

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Valid     bool      `db:"valid" json:"valid"`
	UserAgent string    `db:"user_agent" json:"user_agent"`
	Ip        string    `db:"ip" json:"ip"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateSessionSchema struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
