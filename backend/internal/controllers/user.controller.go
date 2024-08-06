package controllers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/nakshatraraghav/transcodex/backend/internal/schema"
	"github.com/nakshatraraghav/transcodex/backend/internal/services"
	"github.com/nakshatraraghav/transcodex/backend/types"
	"github.com/nakshatraraghav/transcodex/backend/util"
)

var (
	errCreateUser error = errors.New("create user operation failed, internal server error")
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

func (us *UserController) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	body := r.Context().Value(types.ContextKey).(schema.CreateUserSchema)

	exists := us.service.UserExists(r.Context(), body.Email)
	if exists {
		util.WriteError(w, http.StatusConflict, errCreateUser)
		return
	}

	hash, err := argon2id.CreateHash(body.Password, argon2id.DefaultParams)
	if err != nil {
		slog.Error(err.Error())
		util.WriteError(w, http.StatusInternalServerError, errCreateUser)
		return
	}

	body.Password = hash

	user, err := us.service.CreateUser(r.Context(), body)
	if err != nil {
		slog.Error(err.Error())
		util.WriteError(w, http.StatusInternalServerError, "create user operation failed, internal server error")
		return
	}

	util.WriteJSON(w, http.StatusCreated, user)
}
