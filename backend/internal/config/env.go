package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/nakshatraraghav/transcodex/backend/lib"
)

type env struct {
	Addr              string `validate:"required"`
	CONNECTION_STRING string `validate:"required"`
}

var ev env

func LoadEnv() error {
	if err := godotenv.Load(".env.local"); err != nil {
		slog.Error("failed to find the .env.local file, please check for it again")
		return err
	}

	e := env{}

	e.Addr = os.Getenv("PORT")
	e.CONNECTION_STRING = os.Getenv("CONNECTION_STRING")

	vd := lib.GetValidator()

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
