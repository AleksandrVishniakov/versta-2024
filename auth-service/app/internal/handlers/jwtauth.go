package handlers

import (
	"context"
	"errors"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/pkg/e"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/pkg/jwttokens"
	"net/http"
	"strings"
)

const (
	UserIdContextKey = "userId"
	EmailContextKey  = "email"
)

const (
	RefreshTokenCookieKey = "refreshToken"
)

func NewJWTAuthMiddleware(tokensManager jwttokens.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			accessToken, err := GetAccessToken(r)
			if err != nil {
				e.WriteError(w, http.StatusUnauthorized, err.Error())
				return
			}

			accessTokenPayload, err := tokensManager.ParseAccessToken(accessToken)
			if err != nil {
				e.WriteError(w, http.StatusUnauthorized, err.Error())
				return
			}

			ctx = context.WithValue(ctx, UserIdContextKey, accessTokenPayload.UserId)
			ctx = context.WithValue(ctx, EmailContextKey, accessTokenPayload.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAccessToken(r *http.Request) (string, error) {
	headerParts := strings.Split(r.Header.Get("Authorization"), " ")
	if len(headerParts) != 2 {
		return "", errors.New("invalid authorization header length")
	}

	if headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth method")
	}

	return headerParts[1], nil
}
