package middlewares

import (
	"net/http"
	"strings"

	"github.com/nakshatraraghav/transcodex/backend/util"
)

func EnsureApiKeyInRequestHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		key := strings.TrimSpace(r.Header.Get("x-api-key"))

		if key == "" {
			util.WriteError(w, http.StatusUnauthorized, "api key missing from request headers")
			return
		}

		next.ServeHTTP(w, r)
	})
}
