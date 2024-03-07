package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
)

const (
	SessionKey = "sessionKey"
)

const (
	sessionCookieKey = "sessionKey"
)

func Cookies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCtx := r.Context()

		session := getSessionKeyFromCookie(r, sessionCookieKey)

		ctx := context.WithValue(reqCtx, SessionKey, session)

		slog.Debug("cookies",
			slog.String("sessionKey", session),
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getSessionKeyFromCookie(r *http.Request, key string) string {
	sessionCookie, err := r.Cookie(key)
	if errors.Is(err, http.ErrNoCookie) {
		return ""
	}

	if err != nil {
		slog.Error("cookie get error", slog.String("error", err.Error()))
		return ""
	}

	//if time.Now().After(sessionCookie.Expires) {
	//	return ""
	//}

	return sessionCookie.Value
}
