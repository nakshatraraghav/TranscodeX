package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nakshatraraghav/transcodex/backend/internal/config"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type TokenValidationResult struct {
	Valid   bool
	Expired bool
	Error   error
	Claims  jwt.MapClaims
}

func CreateTokens(payload jwt.MapClaims) (*Tokens, error) {
	var tokens Tokens

	env := config.GetEnv()

	attl, _ := time.ParseDuration(env.ACCESS_TOKEN_TTL)
	rttl, _ := time.ParseDuration(env.REFRESH_TOKEN_TTL)

	tok, err := create(payload, attl, env.JWT_PRIVATE_KEY)
	if err != nil {
		return nil, err
	}

	tokens.AccessToken = tok

	tok, err = create(payload, rttl, env.JWT_PRIVATE_KEY)
	if err != nil {
		return nil, err
	}

	tokens.RefreshToken = tok

	return &tokens, nil
}

func create(payload jwt.MapClaims, ttl time.Duration, key string) (string, error) {

	payload["exp"] = time.Now().Add(ttl).Unix()

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return t.SignedString([]byte(key))
}

func ValidateToken(tokenString string) TokenValidationResult {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.GetEnv().JWT_PRIVATE_KEY), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return TokenValidationResult{Valid: false, Expired: true, Error: err}
		}
		return TokenValidationResult{Valid: false, Expired: false, Error: err}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return TokenValidationResult{Valid: true, Expired: false, Error: nil, Claims: claims}
	}

	return TokenValidationResult{Valid: false, Expired: false, Error: errors.New("invalid token")}
}
