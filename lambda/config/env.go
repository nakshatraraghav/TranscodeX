package config

import (
	"os"
	"strings"

	"github.com/nakshatraraghav/transcodex/lambda/lib"
)

type env struct {
	REGION_STRING                string   `validate:"required"`
	BUCKET_NAME                  string   `validate:"required"`
	ECS_CLUSTER_NAME             string   `validate:"required"`
	ECS_TASK_DEFINITION          string   `validate:"required"`
	RDS_DATABASE_USERNAME        string   `validate:"required"`
	RDS_DATABASE_PASSWORD        string   `validate:"required"`
	DATABASE_INSTANCE_IDENTIFIER string   `validate:"required"`
	SUBNET_IDS                   []string `validate:"required"`
	SECURITY_GROUP_ID            string   `validate:"required"`
	CONNECTION_STRING            string
}

var ev env

func LoadEnv() error {
	var e env

	e.REGION_STRING = os.Getenv("REGION_STRING")
	e.BUCKET_NAME = os.Getenv("BUCKET_NAME")
	e.ECS_CLUSTER_NAME = os.Getenv("ECS_CLUSTER_NAME")
	e.ECS_TASK_DEFINITION = os.Getenv("ECS_TASK_DEFINITION")
	e.RDS_DATABASE_USERNAME = os.Getenv("RDS_DATABASE_USERNAME")
	e.RDS_DATABASE_PASSWORD = os.Getenv("RDS_DATABASE_PASSWORD")
	e.DATABASE_INSTANCE_IDENTIFIER = os.Getenv("DATABASE_INSTANCE_IDENTIFIER")
	e.SUBNET_IDS = strings.Split(os.Getenv("SUBNET_IDS"), ",")
	e.SECURITY_GROUP_ID = os.Getenv("SECURITY_GROUP_ID")
	e.CONNECTION_STRING = os.Getenv("CONNECTION_STRING")

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
