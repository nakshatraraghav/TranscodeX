package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
		return nil, fmt.Errorf("error creating session: %w", err)
	}

	return &session, nil
}

func (ss *sessionService) GetSessionByID(ctx context.Context, id uuid.UUID) (*schema.Session, error) {
	var session schema.Session

	q := `SELECT * FROM sessions WHERE id = $1`

	row := ss.db.QueryRowContext(ctx, q, id)
	err := populateSession(row, &session)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("session with id %s not found", id)
		}
		return nil, fmt.Errorf("error retrieving session: %w", err)
	}

	return &session, nil
}

func (ss *sessionService) GetAllActiveSessions(ctx context.Context, id uuid.UUID) ([]schema.Session, error) {
	var sessions []schema.Session

	q := `SELECT id, valid, user_agent, ip, user_id, created_at, updated_at
	FROM sessions
	WHERE user_id = $1 AND valid = true`

	rows, err := ss.db.QueryContext(ctx, q, id)
	if err != nil {
		return nil, fmt.Errorf("error querying active sessions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var session schema.Session

		err := rows.Scan(
			&session.ID,
			&session.Valid,
			&session.UserAgent,
			&session.Ip,
			&session.UserID,
			&session.CreatedAt,
			&session.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning session row: %w", err)
		}

		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over session rows: %w", err)
	}

	return sessions, nil
}

func (ss *sessionService) InvalidateSession(ctx context.Context, id uuid.UUID) error {
	q := `UPDATE sessions
	SET valid = false
	WHERE id = $1`

	result, err := ss.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("error invalidating session: %w", err)
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows for session: %w", err)
	}

	if cnt == 0 {
		return fmt.Errorf("no session found with id %s", id)
	}

	return nil
}

func (ss *sessionService) InvalidateAllSessions(ctx context.Context, id uuid.UUID) error {
	q := `UPDATE sessions
	SET valid = false
	WHERE user_id = $1`

	result, err := ss.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("error invalidating all sessions: %w", err)
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows for user: %w", err)
	}

	if cnt == 0 {
		return fmt.Errorf("no sessions found for user with id %s", id)
	}

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
