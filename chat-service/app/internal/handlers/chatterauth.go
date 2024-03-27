package handlers

import (
	"context"
	"errors"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/services/api/authapi"
	"net/http"
	"strings"
	"time"

	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/services/chatters"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/e"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/jwttokens"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/str"
)

const (
	ChatterIdContextKey = "chatterId"
	StatusContextKey    = "status"

	chatSessionCookieKey = "chatSession"
)

func NewChatterAuth(
	tokensManager jwttokens.Manager,
	chattersStorage chatters.Storage,
) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			var userStatus = authapi.StatusUser
			var chatterId int
			if len(authHeader) != 0 {
				accessToken, err := getAccessToken(r)
				if err != nil {
					e.WriteError(w, http.StatusUnauthorized, err.Error())
					return
				}

				accessTokenPayload, err := tokensManager.ParseAccessToken(accessToken)
				if err != nil {
					e.WriteError(w, http.StatusUnauthorized, err.Error())
					return
				}

				userStatus = accessTokenPayload.Status

				id, err := getChatterIdByUserId(chattersStorage, accessTokenPayload.UserId)
				if err != nil {
					e.WriteError(w, http.StatusInternalServerError, err.Error())
					return
				}

				chatterId = id
			} else {
				cookie, err := r.Cookie(chatSessionCookieKey)
				var session string

				if errors.Is(err, http.ErrNoCookie) {
					session = ""
				} else if err != nil {
					e.WriteError(w, http.StatusUnauthorized, err.Error())
					return
				} else {
					session = cookie.Value
				}

				id, session, err := getChatterIdBySession(chattersStorage, session)
				if err != nil {
					e.WriteError(w, http.StatusInternalServerError, err.Error())
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name:     chatSessionCookieKey,
					Value:    session,
					Path:     "/",
					MaxAge:   int(5 * 24 * time.Hour / time.Second),
					Secure:   false,
					HttpOnly: true,
				})

				chatterId = id
			}

			ctx := context.WithValue(r.Context(), ChatterIdContextKey, chatterId)
			ctx = context.WithValue(ctx, StatusContextKey, userStatus)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getAccessToken(r *http.Request) (string, error) {
	headerParts := strings.Split(r.Header.Get("Authorization"), " ")
	if len(headerParts) != 2 {
		return "", errors.New("invalid authorization header length")
	}

	if headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth method")
	}

	return headerParts[1], nil
}

func getChatterIdByUserId(
	chattersStorage chatters.Storage,
	userId int,
) (int, error) {
	chatter, err := chattersStorage.FindByUserId(userId)
	if errors.Is(err, chatters.ErrChatterNotFound) {
		id, err := chattersStorage.CreateWithId(userId)
		if err != nil {
			return 0, err
		}

		return id, nil
	}
	if err != nil {
		return 0, err
	}

	return chatter.Id, nil
}

func getChatterIdBySession(
	chattersStorage chatters.Storage,
	session string,
) (int, string, error) {
	var chatter *chatters.ChatterDTO
	var err error
	if session != "" {
		chatter, err = chattersStorage.FindBySession(session)
	}

	if errors.Is(err, chatters.ErrChatterNotFound) || session == "" {
		id, session, err := newChatterWithSession(chattersStorage)
		if err != nil {
			return 0, "", err
		}

		return id, session, err
	}

	if err != nil {
		return 0, "", err
	}

	return chatter.Id, session, nil
}

func newChatterWithSession(
	chattersStorage chatters.Storage,
) (int, string, error) {
	session := str.Generate(16)

	id, err := chattersStorage.CreateWithSession(session)
	if err != nil {
		return 0, "", err
	}

	return id, session, nil
}
