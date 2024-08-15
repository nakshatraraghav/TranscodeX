package lib

import (
	"github.com/go-playground/validator/v10"
)

var vd *validator.Validate

func init() {
	vd = validator.New()
}

func GetValidator() *validator.Validate {
	return vd
}
