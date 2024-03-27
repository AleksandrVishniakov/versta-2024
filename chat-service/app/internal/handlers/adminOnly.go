package handlers

import (
	"fmt"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/jwttokens"
	"net/http"

	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/services/api/authapi"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/e"
)

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if fmt.Sprintf("%v", r.Context().Value(StatusContextKey)) != string(authapi.StatusAdmin) {
			e.WriteError(w, http.StatusForbidden, "admin only")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewAdminOnlyFromQueryToken(
	tokensManager jwttokens.Manager,
) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken := r.URL.Query().Get("jwt")
			if accessToken == "" {
				e.WriteError(w, http.StatusUnauthorized, "no jwt token provided")
				return
			}

			accessTokenPayload, err := tokensManager.ParseAccessToken(accessToken)
			if err != nil {
				e.WriteError(w, http.StatusUnauthorized, err.Error())
				return
			}

			if accessTokenPayload.Status != authapi.StatusAdmin {
				e.WriteError(w, http.StatusForbidden, "admin only")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
