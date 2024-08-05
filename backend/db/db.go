package db

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
	"github.com/nakshatraraghav/transcodex/backend/internal/config"
)

func NewPostgresConnection() (*sql.DB, error) {
	env := config.GetEnv()

	db, err := sql.Open("postgres", env.CONNECTION_STRING)
	if err != nil {
		slog.Info("failed to connect to the pg database")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		slog.Info("failed to ping database")
		return nil, err
	}

	slog.Info("connected to the database successfully")
	return db, nil
}
