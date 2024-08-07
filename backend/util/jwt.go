package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nakshatraraghav/transcodex/backend/internal/config"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
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
