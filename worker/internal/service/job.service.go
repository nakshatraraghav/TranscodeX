package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/nakshatraraghav/transcodex/worker/config"
)

type ProcessingJobService struct {
	db *sql.DB
}

func NewProcessingJobService(db *sql.DB) *ProcessingJobService {
	return &ProcessingJobService{
		db: db,
	}
}

func (p *ProcessingJobService) ChangeProcessingJobStatus(ctx context.Context, status string) error {
	env := config.GetEnv()

	q := `UPDATE processing_jobs
	SET status = $1
	WHERE upload_id = $2	
	`

	uid, err := uuid.Parse(env.UPLOAD_ID)
	if err != nil {
		return err
	}

	result, err := p.db.ExecContext(ctx, q, status, uid)
	if err != nil {
		return err
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows for upload: %w", err)
	}

	if cnt == 0 {
		return fmt.Errorf("no upload found with id %s", env.UPLOAD_ID)
	}

	return nil

}
