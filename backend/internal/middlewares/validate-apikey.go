package middlewares

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/nakshatraraghav/transcodex/backend/internal/services"
	"github.com/nakshatraraghav/transcodex/backend/types"
	"github.com/nakshatraraghav/transcodex/backend/util"
)

func ValidateApiKey(service services.ApiKeyService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			key := r.Header.Get("x-api-key")
			if key == "" {
				util.WriteError(w, http.StatusUnauthorized, "api key missing from headers")
				return
			}

			k, err := service.FindApiKey(r.Context(), key)
			if err != nil {
				if err == sql.ErrNoRows {
					util.WriteError(w, http.StatusUnauthorized, "unauthorized, invalid api key")
				} else {
					util.WriteError(w, http.StatusInternalServerError, "failed to get the apikey")
				}
				return
			}

			var kc types.ApiKeyContextKey = "apikey"

			ctx := context.WithValue(r.Context(), kc, k)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
