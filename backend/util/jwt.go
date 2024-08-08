package util

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nakshatraraghav/transcodex/backend/internal/config"
	"github.com/nakshatraraghav/transcodex/backend/types"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Claims struct {
	Exp float64
	SID uuid.UUID
	ID  uuid.UUID
}

func NewJwtClaims(r *http.Request) (Claims, error) {
	claims := r.Context().Value(types.AuthContextKey).(jwt.MapClaims)

	id, ok := claims["uid"]
	sid, sok := claims["sid"]
	exp, eok := claims["exp"]

	if !ok || !sok || !eok {
		return Claims{}, errors.New("failed to parse the claims")
	}

	userIdString := id.(string)
	sessionIdString := sid.(string)
	expirationTime := exp.(float64)

	userUUID, err := uuid.Parse(userIdString)
	if err != nil {
		return Claims{}, err
	}

	sessionUUID, err := uuid.Parse(sessionIdString)
	if err != nil {
		return Claims{}, err
	}

	return Claims{
		Exp: expirationTime,
		ID:  userUUID,
		SID: sessionUUID,
	}, nil
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

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return TokenValidationResult{Valid: false, Expired: false, Error: errors.New("invalid token")}
	}

	return TokenValidationResult{Valid: true, Expired: false, Error: nil, Claims: claims}
}
