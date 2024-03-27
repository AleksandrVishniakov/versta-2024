package jwttokens

import (
	"errors"
	"time"

	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/services/api/authapi"
	"github.com/golang-jwt/jwt"
)

type JWTokens struct {
	RefreshToken string
	AccessToken  string
}

type AccessTokenPayload struct {
	UserId int                `json:"userId"`
	Email  string             `json:"email"`
	Status authapi.UserStatus `json:"status"`
}

type RefreshTokenPayload struct {
	UserId int `json:"userId"`
}

type Manager interface {
	CreateTokens(userId int, payload AccessTokenPayload) (*JWTokens, error)
	RefreshTokens(refreshToken string, payload *AccessTokenPayload) (*JWTokens, error)

	ParseAccessToken(token string) (*AccessTokenPayload, error)
	ParseRefreshToken(token string) (*RefreshTokenPayload, error)
}

type TokensManager struct {
	signature       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewTokensManager(
	signature []byte,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *TokensManager {
	return &TokensManager{
		signature:       signature,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

type accessTokenClaims struct {
	jwt.StandardClaims
	AccessTokenPayload
}

type refreshTokenClaims struct {
	jwt.StandardClaims
	RefreshTokenPayload
}

func (m *TokensManager) CreateTokens(userId int, payload AccessTokenPayload) (*JWTokens, error) {
	tokens := new(JWTokens)

	accessTokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &accessTokenClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(m.accessTokenTTL).Unix(),
		},
		AccessTokenPayload: payload,
	}).SignedString(m.signature)

	if err != nil {
		return nil, err
	}

	tokens.AccessToken = accessTokenStr

	refreshTokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &refreshTokenClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(m.refreshTokenTTL).Unix(),
		},
		RefreshTokenPayload: RefreshTokenPayload{UserId: userId},
	}).SignedString(m.signature)

	if err != nil {
		return nil, err
	}

	tokens.RefreshToken = refreshTokenStr

	return tokens, nil
}

func (m *TokensManager) RefreshTokens(refreshToken string, payload *AccessTokenPayload) (*JWTokens, error) {
	refreshTokenClaims, err := m.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return m.CreateTokens(refreshTokenClaims.UserId, *payload)
}

func (m *TokensManager) ParseAccessToken(token string) (*AccessTokenPayload, error) {
	accessToken, err := jwt.ParseWithClaims(token, &accessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return m.signature, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := accessToken.Claims.(*accessTokenClaims)
	if !ok {
		return nil, errors.New("undefined token claims type")
	}

	return &claims.AccessTokenPayload, nil
}

func (m *TokensManager) ParseRefreshToken(token string) (*RefreshTokenPayload, error) {
	refreshToken, err := jwt.ParseWithClaims(token, &refreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return m.signature, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := refreshToken.Claims.(*refreshTokenClaims)
	if !ok {
		return nil, errors.New("undefined token claims type")
	}

	return &claims.RefreshTokenPayload, nil
}
