package services

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/nakshatraraghav/transcodex/backend/internal/schema"
)

type SessionService interface {
	CreateSession(context.Context, uuid.UUID, string, string) (*schema.Session, error)
	GetSessionByID(context.Context, uuid.UUID) (*schema.Session, error)
	GetAllActiveSessions(context.Context, uuid.UUID) ([]schema.Session, error)
	InvalidateSession(context.Context, uuid.UUID) error
	InvalidateAllSessions(context.Context, uuid.UUID) error
}

type sessionService struct {
	db *sql.DB
}

func NewSessionService(db *sql.DB) SessionService {
	return &sessionService{
		db: db,
	}
}

func (ss *sessionService) CreateSession(ctx context.Context, id uuid.UUID, ip, useragent string) (*schema.Session, error) {

	var session schema.Session

	q := `INSERT INTO sessions (user_agent, ip, user_id)
	VALUES ($1, $2, $3)
	RETURNING id, valid, user_agent, ip, user_id, created_at, updated_at`

	row := ss.db.QueryRowContext(ctx, q, useragent, ip, id)
	err := populateSession(row, &session)

	if err != nil {
		return nil, err
	}

	return &session, err
}

func (ss *sessionService) GetSessionByID(ctx context.Context, id uuid.UUID) (*schema.Session, error) {

	var session schema.Session

	q := `SELECT * FROM sessions WHERE id = $1`

	row := ss.db.QueryRowContext(ctx, q, id)
	err := populateSession(row, &session)

	if err != nil {
		return nil, err
	}

	return &session, err
}

func (ss *sessionService) GetAllActiveSessions(ctx context.Context, id uuid.UUID) ([]schema.Session, error) {
	return nil, nil
}

func (ss *sessionService) InvalidateSession(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (ss *sessionService) InvalidateAllSessions(ctx context.Context, id uuid.UUID) error {
	return nil
}

func populateSession(row *sql.Row, session *schema.Session) error {
	err := row.Scan(
		&session.ID,
		&session.Valid,
		&session.UserAgent,
		&session.Ip,
		&session.UserID,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	return err
}
