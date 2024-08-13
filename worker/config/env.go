package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/nakshatraraghav/transcodex/worker/lib"
)

type env struct {
	MEDIA_TYPE        string            `validate:"required"`
	OBJECT_KEY        string            `validate:"required"`
	BUCKET_NAME       string            `validate:"required"`
	TRANSFORMATIONS   map[string]string `validate:"required"`
	CONNECTION_STRING string            `validate:"required"`
	UPLOAD_ID         string            `validate:"required"`
}

var ev env

func LoadEnv() error {
	err := godotenv.Load(".env.local")
	if err != nil {
		return err
	}

	var e env

	e.MEDIA_TYPE = os.Getenv("MEDIA_TYPE")
	e.BUCKET_NAME = os.Getenv("BUCKET_NAME")
	e.OBJECT_KEY = os.Getenv("OBJECT_KEY")
	e.TRANSFORMATIONS = parseTransformations(os.Getenv("TRANSFORMATIONS"))
	e.CONNECTION_STRING = os.Getenv("CONNECTION_STRING")
	e.UPLOAD_ID = os.Getenv("UPLOAD_ID")

	vd := lib.GetValidator()

	err = vd.Struct(e)
	if err != nil {
		return err
	}

	ev = e

	return nil
}

func GetEnv() env {
	return ev
}

// resize:1980x1090|greyscale

func parseTransformations(transformations string) map[string]string {
	tmap := make(map[string]string)

	if transformations == "" {
		return nil
	}

	transforms := strings.Split(transformations, ",")

	for _, t := range transforms {
		parts := strings.Split(t, ":")

		key := parts[0]
		value := "NO_PARAMETERS"

		if len(parts) > 1 {
			value = parts[1]
		}

		tmap[key] = value
	}

	return tmap

}
