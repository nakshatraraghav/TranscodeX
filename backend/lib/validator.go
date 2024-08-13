package lib

import (
	"github.com/go-playground/validator/v10"
	"github.com/nakshatraraghav/transcodex/backend/internal/schema"
)

var vd *validator.Validate

func init() {
	vd = validator.New()
	vd.RegisterValidation("allowed_operations", schema.ValidateOperationsField)
	vd.RegisterValidation("allowed_jobtype", schema.ValidateJobTypeField)
}

func GetValidator() *validator.Validate {
	return vd
}
