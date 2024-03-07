package authapi

import (
	"context"
	"fmt"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/pkg/apiclient"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/utils"
	"net/http"
)

const (
	sessionCookieKey = "sessionKey"
)

type API interface {
	Create(email string, withEmailVerification bool) (int, error)

	FindBySessionKey(sessionKey string) (*UserDTO, string, error)
	FindByEmail(email string) (*UserDTO, error)

	GetVerificationCode(email string) (*VerificationCodeDTO, error)

	VerifyEmail(email string, verificationCode string) (string, error)
}

type authAPI struct {
	ctx    context.Context
	host   string
	client apiclient.BaseClient
}

func NewAuthAPI(ctx context.Context, host string, client apiclient.BaseClient) API {
	return &authAPI{
		ctx:    ctx,
		host:   host,
		client: client,
	}
}

func (a *authAPI) Create(email string, withEmailVerification bool) (int, error) {
	var url = fmt.Sprintf("%s/api/auth?email=%s&send_email=%t", a.host, email, withEmailVerification)

	req, err := http.NewRequestWithContext(a.ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := a.client.Send(req)
	if err != nil {
		return 0, err
	}

	defer utils.CloseReadCloser(resp.Body)

	var id int

	err = apiclient.ScanResponse(resp, &id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (a *authAPI) FindBySessionKey(sessionKey string) (*UserDTO, string, error) {
	var url = fmt.Sprintf("%s/api/user?session_key=%s", a.host, sessionKey)

	req, err := http.NewRequestWithContext(a.ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := a.client.Send(req)
	if err != nil {
		return nil, "", err
	}

	defer utils.CloseReadCloser(resp.Body)

	var user = UserDTO{}

	err = apiclient.ScanResponse(resp, &user)
	if err != nil {
		return nil, "", err
	}

	var newSessionKey string

	for _, cookie := range resp.Cookies() {
		if cookie.Name == sessionCookieKey {
			newSessionKey = cookie.Value
			break
		}
	}

	return &user, newSessionKey, nil
}

func (a *authAPI) FindByEmail(email string) (*UserDTO, error) {
	var url = fmt.Sprintf("%s/api/user?email=%s", a.host, email)

	req, err := http.NewRequestWithContext(a.ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Send(req)
	if err != nil {
		return nil, err
	}

	defer utils.CloseReadCloser(resp.Body)

	var user = UserDTO{}

	err = apiclient.ScanResponse(resp, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *authAPI) GetVerificationCode(email string) (*VerificationCodeDTO, error) {
	var url = fmt.Sprintf("%s/api/internal/user/%s/verification_code", a.host, email)

	req, err := http.NewRequestWithContext(a.ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Send(req)
	if err != nil {
		return nil, err
	}

	defer utils.CloseReadCloser(resp.Body)

	var verificationCode = VerificationCodeDTO{}

	err = apiclient.ScanResponse(resp, &verificationCode)
	if err != nil {
		return nil, err
	}

	return &verificationCode, nil
}

func (a *authAPI) VerifyEmail(email string, verificationCode string) (string, error) {
	var url = fmt.Sprintf("%s/api/%s/verify?code=%s", a.host, email, verificationCode)

	req, err := http.NewRequestWithContext(a.ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := a.client.Send(req)
	if err != nil {
		return "", err
	}

	defer utils.CloseReadCloser(resp.Body)

	err = apiclient.ScanResponse[any](resp, nil)
	if err != nil {
		return "", err
	}

	var newSessionKey string

	for _, cookie := range resp.Cookies() {
		if cookie.Name == sessionCookieKey {
			newSessionKey = cookie.Value
			break
		}
	}

	return newSessionKey, nil
}
