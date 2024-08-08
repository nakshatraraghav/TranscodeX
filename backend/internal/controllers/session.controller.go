package controllers

import (
	"database/sql"
	"errors"
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

	user, err := sc.userService.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, "user does not exist")
		return
	}

	match, err := argon2id.ComparePasswordAndHash(body.Password, user.Password)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if !match {
		util.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	agent := r.Header.Get("User-Agent")

	session, err := sc.sessionService.CreateSession(r.Context(), user.ID, r.RemoteAddr, agent)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to create the user session")
		return
	}

	tokens, err := util.CreateTokens(jwt.MapClaims{"uid": user.ID, "sid": session.ID})
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	tstr, err := util.Encode(tokens)
	if err != nil {
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
	claims, err := util.NewJwtClaims(r)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	session, err := sc.sessionService.GetSessionByID(r.Context(), claims.SID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			util.WriteError(w, http.StatusNotFound, "session not found")
		} else {
			util.WriteError(w, http.StatusInternalServerError, err.Error())
		}
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
			util.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, sessions)
}

func (sc *SessionController) InvalidateCurrentSessionHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := util.NewJwtClaims(r)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = sc.sessionService.InvalidateSession(r.Context(), claims.SID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			util.WriteError(w, http.StatusNotFound, "session not found")
		} else {
			util.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	http.SetCookie(w, util.ClearCookie("authorization"))

	util.WriteJSON(w, http.StatusOK, "successfully logged out of this session")
}

func (sc *SessionController) InvalidateAllSessionsHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := util.NewJwtClaims(r)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = sc.sessionService.InvalidateAllSessions(r.Context(), claims.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			util.WriteError(w, http.StatusNotFound, "no active sessions found")
		} else {
			util.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	http.SetCookie(w, util.ClearCookie("authorization"))

	util.WriteJSON(w, http.StatusOK, "successfully logged out of all sessions for this user")
}
