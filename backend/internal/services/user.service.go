package services

import (
	"context"
	"database/sql"

	"github.com/nakshatraraghav/transcodex/backend/internal/schema"
)

type UserService interface {
	UserExists(context.Context, string) bool
	GetUserByEmail(context.Context, string) (*schema.User, error)
	CreateUser(context.Context, schema.CreateUserSchema) (*schema.User, error)
}

type userService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) UserService {
	return &userService{
		db: db,
	}
}

func (us *userService) UserExists(ctx context.Context, email string) bool {
	q := "SELECT id FROM users WHERE email = $1"

	_, err := us.db.ExecContext(ctx, q, email)
	return err != nil
}

func (us *userService) GetUserByEmail(ctx context.Context, email string) (*schema.User, error) {
	var user schema.User

	q := `SELECT  
	id, name, username, email, password, created_at, updated_at
	FROM users WHERE email = $1`

	row := us.db.QueryRowContext(ctx, q, email)

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (us *userService) CreateUser(ctx context.Context, body schema.CreateUserSchema) (*schema.User, error) {
	var user schema.User

	q := `INSERT INTO users (name, username, email, password)
	VALUES ($1, $2, $3, $4)
	RETURNING id, name, username, email, password, created_at, updated_at
	`

	err := us.db.QueryRowContext(ctx, q, body.Name, body.Username, body.Email, body.Password).Scan(
		&user.ID, &user.Name, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	user.Password = "<REDACTED>"

	return &user, nil
}
