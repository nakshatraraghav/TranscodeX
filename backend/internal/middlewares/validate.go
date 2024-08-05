package middlewares

import (
	"context"
	"net/http"

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
			util.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		var ckey types.RequestBodyContextKey = "body"

		ctx := context.WithValue(r.Context(), ckey, body)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
