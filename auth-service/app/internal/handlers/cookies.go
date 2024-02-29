package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

const (
	IsCookieAcceptedKey = "isCookieAccepted"
	SessionKey          = "sessionKey"
)

const (
	noCookiesHeaderValue = "none"
	sessionCookieKey     = "sessionKey"
)

func Cookies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCtx := r.Context()

		var isCookieAccepted = r.Header.Get("Cookie") != noCookiesHeaderValue
		var session string

		ctx := context.WithValue(reqCtx, IsCookieAcceptedKey, isCookieAccepted)

		if isCookieAccepted {
			session = getSessionKeyFromCookie(r, sessionCookieKey)
		}

		ctx = context.WithValue(ctx, SessionKey, session)

		slog.Debug("cookies",
			slog.Bool("isAccepted", isCookieAccepted),
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

	if time.Now().After(sessionCookie.Expires) {
		return ""
	}

	return sessionCookie.Value
}
