package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/nakshatraraghav/transcodex/backend/types"
	"github.com/nakshatraraghav/transcodex/backend/util"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract tokens from the single cookie
		cookie, err := r.Cookie("authorization")
		if err != nil {
			util.WriteError(w, http.StatusUnauthorized, "authentication cookies missing")
			return
		}

		// Decode the URL-encoded JSON
		c, err := url.QueryUnescape(cookie.Value)
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, "invalid cookie format")
			return
		}

		var tokens util.Tokens
		if err := json.Unmarshal([]byte(c), &tokens); err != nil {
			util.WriteError(w, http.StatusUnauthorized, "invalid token format")
			return
		}

		// Validate access token
		at := util.ValidateToken(tokens.AccessToken)
		if !at.Valid {
			// Access token is invalid (tampered or expired)
			util.WriteError(w, http.StatusUnauthorized, "invalid access token")
			return
		}

		// Check if access token is expired
		if at.Expired {
			// Only try to refresh if access token is expired, not if it's tampered
			rt := util.ValidateToken(tokens.RefreshToken)
			if !rt.Valid {
				// Refresh token is also invalid
				util.WriteError(w, http.StatusUnauthorized, "invalid refresh token")
				return
			}

			// Refresh token is valid, generate new tokens
			ntok, err := util.CreateTokens(rt.Claims)
			if err != nil {
				util.WriteError(w, http.StatusInternalServerError, "failed to create new tokens")
				return
			}

			// Set new cookie
			ec, err := json.Marshal(ntok)
			if err != nil {
				util.WriteError(w, http.StatusInternalServerError, "failed to encode tokens")
				return
			}
			encoded := url.QueryEscape(string(ec))
			http.SetCookie(w, &http.Cookie{
				Name:     "authorization",
				Value:    encoded,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})

			// Use the new access token for this request
			at = util.ValidateToken(ntok.AccessToken)
		}

		// Token is valid, proceed with the request
		ctx := context.WithValue(r.Context(), types.AuthContextKey, at.Claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
