package services

import (
	"context"

	"github.com/nakshatraraghav/transcodex/backend/internal/schema"
)

type UserService interface {
	GetUserByEmail(context.Context, string) (*schema.User, error)
	CreateUser(context.Context, schema.CreateUserSchema) (*schema.User, error)
}
