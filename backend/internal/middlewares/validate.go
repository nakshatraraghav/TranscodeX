package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/nakshatraraghav/transcodex/backend/lib"
	"github.com/nakshatraraghav/transcodex/backend/types"
	"github.com/nakshatraraghav/transcodex/backend/util"
)

func ValidateRequestBody[schema any](next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var body schema

		if err := util.Decode(r, &body); err != nil {
			util.WriteError(w, http.StatusBadRequest, "internal server error")
			return
		}

		vd := lib.GetValidator()

		if err := vd.Struct(body); err != nil {
			verr, ok := err.(validator.ValidationErrors)
			if !ok {
				util.WriteError(w, http.StatusBadRequest, "request validation errors")
				return
			}

			fields := make(map[string]string)

			for _, e := range verr {
				fields[e.Field()] = getErrorMessage(e)
			}

			util.WriteError(w, http.StatusBadRequest, map[string]any{
				"fields": fields,
			})

			return
		}

		var ckey types.RequestBodyContextKey = "body"

		ctx := context.WithValue(r.Context(), ckey, body)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "invalid email format, please fix it"
	case "min":
		return fmt.Sprintf("should be at least %s characters long", e.Param())
	case "max":
		return fmt.Sprintf("should be at most %s characters long", e.Param())
	default:
		return "invalid value"
	}
}
