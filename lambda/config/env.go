package config

import (
	"os"

	"github.com/nakshatraraghav/transcodex/lambda/lib"
)

type env struct {
	AWS_REGION                   string `validate:"required"`
	DATABASE_SECRET_ID           string `validate:"required"`
	DATABASE_INSTANCE_IDENTIFIER string `validate:"required"`
	CONNECTION_STRING_SECRET_ID  string `validate:"required"`
}

var ev env

func LoadEnv() error {
	var e env

	e.AWS_REGION = os.Getenv("AWS_REGION")
	e.DATABASE_SECRET_ID = os.Getenv("DATABASE_SECRET_ID")
	e.DATABASE_INSTANCE_IDENTIFIER = os.Getenv("DATABASE_INSTANCE_IDENTIFIER")
	e.CONNECTION_STRING_SECRET_ID = os.Getenv("CONNECTION_STRING_SECRET_ID")

	vd := lib.GetValidator()

	err := vd.Struct(e)
	if err != nil {
		return err
	}

	ev = e

	return nil
}

func Getenv() env {
	return ev
}
