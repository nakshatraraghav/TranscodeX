package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/nakshatraraghav/transcodex/backend/lib"
)

type env struct {
	Addr              string `validate:"required"`
	CONNECTION_STRING string `validate:"required"`
	JWT_PRIVATE_KEY   string `validate:"required"`
	ACCESS_TOKEN_TTL  string `validate:"required,duration"`
	REFRESH_TOKEN_TTL string `validate:"required,duration"`
	BUCKET_NAME       string `validate:"required"`
}

var ev env

func ValidateDuration(fl validator.FieldLevel) bool {
	_, err := time.ParseDuration(fl.Field().String())
	return err == nil
}

func LoadEnv() error {
	if err := godotenv.Load(".env.local"); err != nil {
		slog.Error("failed to find the .env.local file, please check for it again")
		return err
	}

	e := env{}

	e.Addr = os.Getenv("PORT")
	e.CONNECTION_STRING = os.Getenv("CONNECTION_STRING")
	e.JWT_PRIVATE_KEY = os.Getenv("JWT_PRIVATE_KEY")
	e.ACCESS_TOKEN_TTL = os.Getenv("ACCESS_TOKEN_TTL")
	e.REFRESH_TOKEN_TTL = os.Getenv("REFRESH_TOKEN_TTL")
	e.BUCKET_NAME = os.Getenv("BUCKET_NAME")

	vd := lib.GetValidator()
	vd.RegisterValidation("duration", ValidateDuration)

	if err := vd.Struct(e); err != nil {
		slog.Error("failed to validate the environment variables file please check it again")
		return err
	}

	ev = e

	return nil
}

func GetEnv() env {
	return ev
}
