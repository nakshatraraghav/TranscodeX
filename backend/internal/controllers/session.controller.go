package controllers

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nakshatraraghav/transcodex/backend/internal/schema"
	"github.com/nakshatraraghav/transcodex/backend/internal/services"
	"github.com/nakshatraraghav/transcodex/backend/types"
	"github.com/nakshatraraghav/transcodex/backend/util"
)

type SessionController struct {
	userService    services.UserService
	sessionService services.SessionService
}

func NewSessionController(us services.UserService, ss services.SessionService) *SessionController {
	return &SessionController{
		userService:    us,
		sessionService: ss,
	}
}

func (sc *SessionController) CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	body := r.Context().Value(types.ContextKey).(schema.CreateSessionSchema)

	// Step 1: Check if the user exists or not
	user, err := sc.userService.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, "user does not exist")
		return
	}

	// Step 2: check the password of this user
	match, err := argon2id.ComparePasswordAndHash(body.Password, user.Password)
	if err != nil {
		slog.Error(err.Error())
		util.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if !match {
		util.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// Step 3: create session for the user
	agent := r.Header.Get("User-Agent")

	session, err := sc.sessionService.CreateSession(r.Context(), user.ID, r.RemoteAddr, agent)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to create the user session, internal server error")
		return
	}

	// Step 4: Create the json tokens
	tokens, err := util.CreateTokens(jwt.MapClaims{"uid": user.ID, "sid": session.ID})
	if err != nil {
		slog.Error(err.Error())
		util.WriteError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	tstr, err := util.Encode(tokens)
	if err != nil {
		slog.Error(err.Error())
		util.WriteError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	encoded := url.QueryEscape(string(tstr))

	cookie := http.Cookie{
		Name:     "authorization",
		Path:     "/",
		Value:    encoded,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)

	util.WriteJSON(w, http.StatusCreated, "session created, cookies set")
}

func (sc *SessionController) GetCurrentSessionHandler(w http.ResponseWriter, r *http.Request) {
	// authenticated route
	claims, err := util.NewJwtClaims(r)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	session, err := sc.sessionService.GetSessionByID(r.Context(), claims.SID)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, session)
}

func (sc *SessionController) GetAllActiveSessionsHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := util.NewJwtClaims(r)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	sessions, err := sc.sessionService.GetAllActiveSessions(r.Context(), claims.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			util.WriteError(w, http.StatusNotFound, "No active sessions found")
		} else {
			util.WriteError(w, http.StatusInternalServerError, "Failed to retrieve sessions")
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, sessions)
}
